package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/apoldev/trackchecker/internal/app/crawler"
	"github.com/apoldev/trackchecker/internal/app/models"
	"github.com/apoldev/trackchecker/internal/app/track/usecase"
	"github.com/apoldev/trackchecker/internal/app/track/usecase/mocks"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var topic string = "topicName"

func TestTracking_Tracking(t *testing.T) {
	cases := []struct {
		name           string
		trackingNumber *models.TrackingNumber
		results        *models.Crawler
		expectError    error
	}{
		{
			name: "valid data",
			trackingNumber: &models.TrackingNumber{
				Code: "111",
			},
			results: &models.Crawler{
				Results: make([]models.CrawlerResult, 0),
			},
			expectError: nil,
		},

		{
			name:           "nil track",
			trackingNumber: nil,
			results:        nil,
			expectError:    crawler.ErrTrackIsNil,
		},

		{
			name: "no spiders",
			trackingNumber: &models.TrackingNumber{
				Code: "111",
			},
			results: &models.Crawler{
				Status: crawler.StatusNoSpiders,
			},
			expectError: nil,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			publisherMock := mocks.NewPublisher(t)
			trackResultMock := mocks.NewTrackResultRepo(t)
			crawlerMock := mocks.NewCrawler(t)
			logger := logrus.New()

			tracking := usecase.NewTracking(publisherMock, topic, logger, crawlerMock, trackResultMock)

			crawlerMock.On("Start", c.trackingNumber).
				Return(c.results, c.expectError).
				Once()

			result, err := tracking.Tracking(c.trackingNumber)

			if c.expectError != nil {
				require.Empty(t, result)
				require.EqualError(t, err, c.expectError.Error())
			} else {
				require.Equal(t, result, c.results)
				require.NoError(t, err)
			}
		})
	}
}

func TestTracking_GetTrackingResult(t *testing.T) {
	cases := []struct {
		name         string
		id           string
		expectError  error
		expectResult []*models.Crawler
	}{
		{
			name:         "valid data",
			id:           "111",
			expectError:  nil,
			expectResult: []*models.Crawler{},
		},

		{
			name:         "Error",
			id:           "111",
			expectError:  errors.New("my error"),
			expectResult: nil,
		},
	}
	for _, c := range cases {
		c := c

		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			publisherMock := mocks.NewPublisher(t)
			trackResultMock := mocks.NewTrackResultRepo(t)
			crawlerMock := mocks.NewCrawler(t)
			logger := logrus.New()

			tracking := usecase.NewTracking(publisherMock, topic, logger, crawlerMock, trackResultMock)

			ctx := context.Background()

			trackResultMock.
				On("Get", ctx, c.id).
				Return(c.expectResult, c.expectError).
				Once()

			data, err := tracking.GetTrackingResult(ctx, c.id)

			if c.expectError != nil {
				require.EqualError(t, err, c.expectError.Error())
				require.Empty(t, data)
			} else {
				require.Equal(t, c.expectResult, data)
			}
		})
	}
}

func TestTracking_SaveTrackingResult(t *testing.T) {
	cases := []struct {
		name           string
		trackingNumber string
		results        *models.Crawler
		expectError    error
	}{
		{
			name:           "valid data",
			trackingNumber: "111",
			results:        &models.Crawler{},
			expectError:    nil,
		},
		{
			name:           "error",
			trackingNumber: "111",
			results:        &models.Crawler{},
			expectError:    errors.New("error"),
		},

		{
			name:           "",
			trackingNumber: "111",
			results:        nil,
			expectError:    errors.New("error"),
		},
	}

	for _, c := range cases {
		c := c

		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			publisherMock := mocks.NewPublisher(t)
			trackResultMock := mocks.NewTrackResultRepo(t)
			crawlerMock := mocks.NewCrawler(t)
			logger := logrus.New()

			tracking := usecase.NewTracking(publisherMock, topic, logger, crawlerMock, trackResultMock)

			track := &models.TrackingNumber{
				Code: c.trackingNumber,
				UUID: uuid.NewString(),
			}

			ctx := context.Background()
			trackResultMock.
				On("Set", ctx, track, c.results).
				Return(c.expectError).
				Once()

			err := tracking.SaveTrackingResult(ctx, track, c.results)

			if c.expectError != nil {
				require.EqualError(t, err, c.expectError.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestTracking_PublishTrackingNumberToQueue(t *testing.T) {
	t.Parallel()
	publisherMock := mocks.NewPublisher(t)
	trackResultMock := mocks.NewTrackResultRepo(t)
	crawlerMock := mocks.NewCrawler(t)
	logger := logrus.New()
	tracking := usecase.NewTracking(publisherMock, topic, logger, crawlerMock, trackResultMock)

	expextError := "error publish tracking number to queue"

	reqID := uuid.NewString()
	trackingNumber := "111"
	tracks := []models.TrackingNumber{
		models.TrackingNumber{
			RequestID: reqID,
			Code:      trackingNumber,
			UUID:      uuid.NewString(),
		},
	}

	ctx := context.Background()
	publisherMock.
		On("Publish", ctx, topic, mock.Anything).
		Return(nil).
		Once()

	res, err := tracking.PublishTrackingNumbersToQueue(ctx, reqID, []string{trackingNumber})

	require.NoError(t, err)
	require.Len(t, tracks, 1)
	require.Equal(t, tracks[0].Code, res[0].Code)
	require.Equal(t, tracks[0].RequestID, res[0].RequestID)
	require.Len(t, tracks[0].UUID, 36)

	// test error
	publisherMock.
		On("Publish", ctx, topic, mock.Anything).
		Return(errors.New(expextError)).
		Once()

	res, err = tracking.PublishTrackingNumbersToQueue(ctx, reqID, []string{trackingNumber})

	require.EqualError(t, err, expextError)
	require.Empty(t, res)
}
