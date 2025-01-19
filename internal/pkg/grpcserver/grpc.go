package grpcserver

import (
	trackingService "github.com/apoldev/trackchecker/internal/app/grpcservice"
	grpctrack "github.com/apoldev/trackchecker/internal/app/track/delivery/grpc"
	"github.com/apoldev/trackchecker/internal/app/track/usecase"
	"github.com/apoldev/trackchecker/internal/pkg/logger"
	"google.golang.org/grpc"
)

func NewGRPCServer(log logger.Logger, trackingUC *usecase.Tracking) *grpc.Server {
	grpcServer := grpc.NewServer()
	newTrackingAPI := grpctrack.NewTrackGRPCApi(log, trackingUC)
	trackingService.RegisterTrackingServer(grpcServer, newTrackingAPI)
	return grpcServer
}
