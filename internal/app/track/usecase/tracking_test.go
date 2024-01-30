package usecase_test

import (
	"errors"
	"github.com/apoldev/trackchecker/internal/app/crawler"
	"github.com/apoldev/trackchecker/internal/app/models"
	"github.com/apoldev/trackchecker/internal/app/track/usecase"
	"github.com/apoldev/trackchecker/internal/app/track/usecase/mocks"
	"github.com/apoldev/trackchecker/pkg/scraper"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTracking_Tracking(t *testing.T) {
	cases := []struct {
		name           string
		trackingNumber string
		spiders        []*models.Spider
		results        map[string]models.CrawlerResult
		expectError    error
	}{
		{
			name:           "valid data",
			trackingNumber: "111",
			spiders: []*models.Spider{
				{
					Scraper: scraper.Scraper{
						Code: "code",
					},
				},
			},

			results: map[string]models.CrawlerResult{
				"code": {
					Result: []byte(`{}`),
				},
			},
			expectError: nil,
		},

		{
			name:           "empty spiders",
			trackingNumber: "111",
			spiders:        nil,
			results:        nil,
			expectError:    crawler.ErrNoSpiders,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			publisherMock := mocks.NewPublisher(t)
			trackResultMock := mocks.NewTrackResultRepo(t)
			spiderMock := mocks.NewSpiderRepo(t)
			logger := logrus.New()

			tracking := usecase.NewTracking(publisherMock, logger, spiderMock, trackResultMock)

			track := &models.TrackingNumber{
				Code: c.trackingNumber,
			}

			spiderMock.On("FindSpidersByTrackingNumber", c.trackingNumber).
				Return(c.spiders).
				Once()

			result, err := tracking.Tracking(track)

			if c.expectError != nil {
				require.Empty(t, result)
				require.EqualError(t, err, c.expectError.Error())
			} else {
				require.Len(t, result, len(c.results))
				require.NoError(t, err)
			}

		})
	}

}

func TestTracking_GetTrackingResult(t *testing.T) {
	cases := []struct {
		name        string
		id          string
		expectError error
		expectBytes []byte
	}{
		{
			name:        "valid data",
			id:          "111",
			expectError: nil,
			expectBytes: []byte(`{"data": "ok"}`),
		},

		{
			name:        "Error",
			id:          "111",
			expectError: errors.New("my error"),
			expectBytes: nil,
		},
	}

	for _, c := range cases {
		c := c

		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			publisherMock := mocks.NewPublisher(t)
			trackResultMock := mocks.NewTrackResultRepo(t)
			spiderMock := mocks.NewSpiderRepo(t)
			logger := logrus.New()

			tracking := usecase.NewTracking(publisherMock, logger, spiderMock, trackResultMock)

			trackResultMock.
				On("Get", c.id).
				Return(c.expectBytes, c.expectError).
				Once()

			data, err := tracking.GetTrackingResult(c.id)

			if c.expectError != nil {
				require.EqualError(t, err, c.expectError.Error())
				require.Empty(t, data)
			} else {
				require.Equal(t, c.expectBytes, data)
			}
		})
	}

}

func TestTracking_SaveTrackingResult(t *testing.T) {
	cases := []struct {
		name           string
		trackingNumber string
		results        map[string]models.CrawlerResult
		expectError    error
	}{
		{
			name:           "valid data",
			trackingNumber: "111",
			results: map[string]models.CrawlerResult{
				"spider1": models.CrawlerResult{},
			},
			expectError: nil,
		},
		{
			name:           "error",
			trackingNumber: "111",
			results: map[string]models.CrawlerResult{
				"spider1": models.CrawlerResult{},
			},
			expectError: errors.New("error"),
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
			spiderMock := mocks.NewSpiderRepo(t)
			logger := logrus.New()

			tracking := usecase.NewTracking(publisherMock, logger, spiderMock, trackResultMock)

			track := &models.TrackingNumber{
				Code: c.trackingNumber,
				UUID: uuid.NewString(),
			}

			trackResultMock.
				On("Set", track.UUID, mock.Anything).
				Return(c.expectError).
				Once()

			err := tracking.SaveTrackingResult(track, c.results)

			if c.expectError != nil {
				require.EqualError(t, err, c.expectError.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}

}

func TestTracking_PublishTrackingNumberToQueue(t *testing.T) {
	publisherMock := mocks.NewPublisher(t)
	trackResultMock := mocks.NewTrackResultRepo(t)
	spiderMock := mocks.NewSpiderRepo(t)
	logger := logrus.New()
	tracking := usecase.NewTracking(publisherMock, logger, spiderMock, trackResultMock)

	expextError := "error publish tracking number to queue"
	trackingNumber := "111"
	track := models.TrackingNumber{
		Code: trackingNumber,
		UUID: uuid.NewString(),
	}

	publisherMock.
		On("Publish", mock.Anything).
		Return(nil).
		Once()

	res, err := tracking.PublishTrackingNumberToQueue(trackingNumber)

	require.NoError(t, err)
	require.Equal(t, track.Code, res.Code)
	require.Len(t, track.UUID, 36)

	// test error
	publisherMock.
		On("Publish", mock.Anything).
		Return(errors.New(expextError)).
		Once()

	res, err = tracking.PublishTrackingNumberToQueue(trackingNumber)

	require.EqualError(t, err, expextError)
	require.Equal(t, models.TrackingNumber{}, res)

}
