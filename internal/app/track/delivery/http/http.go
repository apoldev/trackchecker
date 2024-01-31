package http

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/apoldev/trackchecker/internal/app/models"

	"github.com/apoldev/trackchecker/pkg/logger"
	"github.com/gin-gonic/gin"
)

const (
	ErrIDEmpty             = "id is empty"
	ErrTrackingNumberEmpty = "tracking number is empty"
	ErrorPublishToQueue    = "error publish tracking number to queue"
	ErrorNotFound          = "not found"
)

// QueuePublisher is an interface for publish tracking number to queue.
type QueuePublisher interface {
	PublishTrackingNumbersToQueue(id string, trackingNumbers []string) ([]models.TrackingNumber, error)
}

type TrackingResultGetter interface {
	GetTrackingResult(id string) ([]*models.Crawler, error)
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

type ResponseTrackingResult struct {
	Status bool              `json:"status"`
	Error  string            `json:"error,omitempty"`
	Data   []*models.Crawler `json:"data,omitempty"`
}

func (h *TrackHandler) GetTrackingNumberResultHandler(c *gin.Context) {
	var err error
	var data []*models.Crawler

	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, ErrIDEmpty)
		return
	}

	data, err = h.tracking.GetTrackingResult(id)
	if err != nil || data == nil {
		c.JSON(http.StatusNotFound, ResponseTrackingResult{
			Error: ErrorNotFound,
		})
		return
	}

	c.JSON(http.StatusOK, ResponseTrackingResult{
		Status: true,
		Data:   data,
	})
}

type RequestCrawler struct {
	TrackingNumbers []string `json:"tracking_numbers"`
}

type ResponseTrackingNumber struct {
	TrackingNumber string `json:"tracking_number"`
}
type ResponseCrawler struct {
	TrackingID      string                  `json:"tracking_id"`
	TrackingNumbers []models.TrackingNumber `json:"tracking_numbers"`
}

func (h *TrackHandler) TrackingNumberCrawlerHandler(c *gin.Context) {
	var err error
	var req RequestCrawler

	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrTrackingNumberEmpty)
		return
	}

	trackingID := uuid.NewString()
	tracks, err := h.tracking.PublishTrackingNumbersToQueue(trackingID, req.TrackingNumbers)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorPublishToQueue)
		return
	}

	resp := ResponseCrawler{
		TrackingID: trackingID,
	}
	for i := range tracks {
		resp.TrackingNumbers = append(resp.TrackingNumbers, models.TrackingNumber{
			Code: tracks[i].Code,
			UUID: tracks[i].UUID,
		})
	}

	c.JSON(http.StatusCreated, resp)
}
