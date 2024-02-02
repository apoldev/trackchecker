package grpctrack

import (
	"context"
	"github.com/apoldev/trackchecker/internal/pkg/logger"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	appmodels "github.com/apoldev/trackchecker/internal/app/models"
	trackingService "github.com/apoldev/trackchecker/internal/app/track/proto"
)

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
	if len(in.TrackingNumbers) == 0 {
		return nil, status.Error(codes.InvalidArgument, "tracking numbers is empty")
	}

	trackingID := uuid.NewString()
	tracks, err := s.tracking.PublishTrackingNumbersToQueue(ctx, trackingID, in.TrackingNumbers)
	if err != nil {
		return nil, status.Error(codes.Internal, "error publish tracking number to queue")
	}

	result := trackingService.PostTrackingResponse{
		TrackingId:      trackingID,
		TrackingNumbers: nil,
	}
	for i := range tracks {
		result.TrackingNumbers = append(result.TrackingNumbers, &trackingService.PostTrack{
			Code: tracks[i].Code,
			Uuid: tracks[i].UUID,
		})
	}

	return &result, nil
}

func (s *TrackGRPCApi) GetResult(
	ctx context.Context,
	in *trackingService.GetTrackingID,
) (*trackingService.GetTrackingResponse, error) {

	crawlers, err := s.tracking.GetTrackingResult(ctx, in.GetId())
	if err != nil || len(crawlers) == 0 {
		return nil, status.Error(codes.NotFound, "tracking results not found")
	}

	data := make([]*trackingService.TrackResponse, 0, len(crawlers))
	for i := range crawlers {
		c := crawlers[i]
		results := make([]*trackingService.TrackResult, 0, len(c.Results))
		for j := range c.Results {
			r := c.Results[j]

			bytes, marshalErr := r.Result.MarshalJSON()
			if marshalErr != nil {
				s.logger.Warnf("error marshal tracking result: %v", marshalErr)
				continue
			}

			results = append(results, &trackingService.TrackResult{
				Error:          r.Err,
				ExecuteTime:    float32(r.ExecuteTime),
				Result:         string(bytes),
				Spider:         r.Spider,
				TrackingNumber: r.TrackingNumber,
			})
		}

		data = append(data, &trackingService.TrackResponse{
			Uuid:   c.UUID,
			Status: c.Status,
			Code:   c.Code,
			Id:     c.RequestID,
			Result: results,
		})
	}
	return &trackingService.GetTrackingResponse{
		Status:   true,
		Tracking: data,
	}, nil
}
