package http

import (
	"context"
	"net/http"

	"github.com/apoldev/trackchecker/internal/pkg/logger"

	appmodels "github.com/apoldev/trackchecker/internal/app/models"
	"github.com/apoldev/trackchecker/internal/app/restapi/models"
	"github.com/apoldev/trackchecker/internal/app/restapi/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/google/uuid"
)

const (
	ErrorPublishToQueue = "error publish tracking number to queue"
	ErrorNotFound       = "not found"
)

// QueuePublisher is an interface for publish tracking number to queue.
type QueuePublisher interface {
	PublishTrackingNumbersToQueue(
		ctx context.Context,
		id string,
		trackingNumbers []string,
	) ([]appmodels.TrackingNumber, error)
}

type TrackingResultGetter interface {
	GetTrackingResult(ctx context.Context, id string) ([]*appmodels.Crawler, error)
}

type Tracking interface {
	QueuePublisher
	TrackingResultGetter
}

type TrackHandler struct {
	logger   logger.Logger
	tracking Tracking
}

func NewTrackHandler(log logger.Logger, tracking Tracking) *TrackHandler {
	return &TrackHandler{
		logger:   log,
		tracking: tracking,
	}
}

func (h *TrackHandler) PostTrackingResultHandler(params operations.PostTrackParams) middleware.Responder {
	ctx := context.Background()
	trackingID := uuid.NewString()
	tracks, err := h.tracking.PublishTrackingNumbersToQueue(ctx, trackingID, params.Body.TrackingNumbers)
	if err != nil {
		return operations.NewPostTrackDefault(http.StatusBadRequest).WithPayload(&models.Error{
			Message: swag.String(ErrorPublishToQueue),
		})
	}

	result := models.RequestResult{
		TrackingID:      trackingID,
		TrackingNumbers: nil,
	}
	for i := range tracks {
		result.TrackingNumbers = append(result.TrackingNumbers, &models.TrackingNumber{
			Code: tracks[i].Code,
			UUID: tracks[i].UUID,
		})
	}
	return operations.NewPostTrackCreated().WithPayload(&result)
}

func (h *TrackHandler) GetTrackingResultHandler(params operations.GetResultsParams) middleware.Responder {
	ctx := context.Background()
	crawlers, err := h.tracking.GetTrackingResult(ctx, params.ID)
	if err != nil || len(crawlers) == 0 {
		return operations.NewGetResultsDefault(http.StatusNotFound).WithPayload(&models.Error{
			Message: swag.String(ErrorNotFound),
		})
	}
	data := make([]*models.Result, 0, len(crawlers))
	for i := range crawlers {
		c := crawlers[i]
		results := make([]*models.SpiderResults, 0, len(c.Results))
		for j := range c.Results {
			r := c.Results[j]
			results = append(results, &models.SpiderResults{
				Error:          r.Err,
				ExecuteTime:    r.ExecuteTime,
				Result:         r.Result,
				Spider:         r.Spider,
				TrackingNumber: r.TrackingNumber,
			})
		}
		data = append(data, &models.Result{
			UUID:    c.UUID,
			Status:  c.Status,
			Code:    c.Code,
			ID:      c.RequestID,
			Results: results,
		})
	}
	return operations.NewGetResultsOK().WithPayload(&models.TrackingResult{
		Status: true,
		Data:   data,
	})
}
