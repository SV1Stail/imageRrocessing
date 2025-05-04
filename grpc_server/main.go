package main

import (
	"net"
	"os"

	pb "github.com/SV1Stail/imageRrocessing/grpc_server/imageRrocessing/gen"
	"github.com/SV1Stail/imageRrocessing/grpc_server/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"google.golang.org/grpc"
)

var (
	GrpcPort string
)

func init() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	GrpcPort = os.Getenv("GRPC_PORT")
}

func main() {
	lis, err := net.Listen("tcp", ":"+GrpcPort)
	if err != nil {
		log.Err(err).Msg("failed start listen")
	}

	grpcServer := grpc.NewServer()
	pb.RegisterImageProcessingServiceServer(grpcServer, &server.Server{})

	log.Info().Msgf("Server is running on port %s...", GrpcPort)

	if err := grpcServer.Serve(lis); err != nil {
		log.Err(err).Msg("failed to serve")
	}

}
