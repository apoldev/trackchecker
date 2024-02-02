package grpctrack

import (
	"context"
	"fmt"
	appmodels "github.com/apoldev/trackchecker/internal/app/models"
	trackingService "github.com/apoldev/trackchecker/internal/app/track/proto"
	"github.com/apoldev/trackchecker/pkg/logger"
)

type QueuePublisher interface {
	PublishTrackingNumbersToQueue(id string, trackingNumbers []string) ([]appmodels.TrackingNumber, error)
}

type TrackingResultGetter interface {
	GetTrackingResult(id string) ([]*appmodels.Crawler, error)
}

type Tracking interface {
	QueuePublisher
	TrackingResultGetter
}

type TrackGRPCApi struct {
	trackingService.UnimplementedTrackingServer

	tracking Tracking
	logger   logger.Logger
}

func NewTrackGRPCApi(log logger.Logger, tracking Tracking) *TrackGRPCApi {
	return &TrackGRPCApi{
		logger:   log,
		tracking: tracking,
	}
}

func (s *TrackGRPCApi) PostTracking(
	ctx context.Context,
	in *trackingService.PostTrackingRequest,
) (*trackingService.PostTrackingResponse, error) {
	fmt.Println(in.TrackingNumbers)

	return &trackingService.PostTrackingResponse{
		TrackingId: "123",
	}, nil
}

func (s *TrackGRPCApi) GetResult(
	ctx context.Context,
	in *trackingService.GetTrackingID,
) (*trackingService.GetTrackingResponse, error) {

	return &trackingService.GetTrackingResponse{
		Status: true,
	}, nil
}