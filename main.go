package main

import (
	"net"

	"github.com/SV1Stail/imageRrocessing/server"
	pb "github.com/SV1Stail/imageRrocessing/server/github.com/SV1Stail/imageRrocessing/gen"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"google.golang.org/grpc"
)

func main() {

	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Err(err).Msg("failed start listen")
	}

	grpcServer := grpc.NewServer()
	pb.RegisterImageProcessingServiceServer(grpcServer, &server.Server{})

	log.Info().Msg("Server is running on port 50051...")

	if err := grpcServer.Serve(lis); err != nil {
		log.Err(err).Msg("failed to serve")
	}

}
