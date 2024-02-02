package grpcserver

import (
	"github.com/apoldev/trackchecker/internal/app/config"
	grpctrack "github.com/apoldev/trackchecker/internal/app/track/delivery/grpc"
	trackingService "github.com/apoldev/trackchecker/internal/app/track/proto"
	"github.com/apoldev/trackchecker/internal/app/track/usecase"
	"github.com/apoldev/trackchecker/pkg/logger"
	"google.golang.org/grpc"
)

func NewGRPCServer(log logger.Logger, cfg config.GRPCServer, trackingUC *usecase.Tracking) *grpc.Server {

	grpcServer := grpc.NewServer()
	newTrackingApi := grpctrack.NewTrackGRPCApi(log, trackingUC)
	trackingService.RegisterTrackingServer(grpcServer, newTrackingApi)

	return grpcServer

	//grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCServer.Port))
	//if err != nil {
	//	logger.Fatal(err)
	//}
	//grpcServer.Serve(grpcListener)
	//defer grpcListener.Close()
}
