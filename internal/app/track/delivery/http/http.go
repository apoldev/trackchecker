package http

import (
	"encoding/json"
	"net/http"

	"github.com/apoldev/trackchecker/internal/app/models"

	"github.com/apoldev/trackchecker/internal/app/track/usecase"
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
	PublishTrackingNumberToQueue(id string) (models.TrackingNumber, error)
}

type TrackingResultGetter interface {
	GetTrackingResult(id string) ([]byte, error)
}

type Tracking interface {
	QueuePublisher
	TrackingResultGetter
}

type TrackHandler struct {
	logger   logger.Logger
	tracking Tracking
}

func NewTrackHandler(log logger.Logger, tracking *usecase.Tracking) *TrackHandler {
	return &TrackHandler{
		logger:   log,
		tracking: tracking,
	}
}

type ResponseTrackingResult struct {
	Status bool            `json:"status"`
	Error  string          `json:"error,omitempty"`
	Data   json.RawMessage `json:"data,omitempty"`
}

func (h *TrackHandler) GetTrackingNumberResultHandler(c *gin.Context) {
	var err error
	var data []byte

	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, ErrIDEmpty)
		return
	}

	data, err = h.tracking.GetTrackingResult(id)
	if err != nil {
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
	TrackingNumber string `json:"tracking_number"`
}

type ResponseCrawler struct {
	models.TrackingNumber
}

func (h *TrackHandler) TrackingNumberCrawlerHandler(c *gin.Context) {
	var err error
	var req RequestCrawler

	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrTrackingNumberEmpty)
		return
	}

	track, err := h.tracking.PublishTrackingNumberToQueue(req.TrackingNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorPublishToQueue)
		return
	}

	c.JSON(http.StatusCreated, ResponseCrawler{
		track,
	})
}
