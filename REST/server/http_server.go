package server

import (
	"context"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/SV1Stail/imageRrocessing/REST/server/gen"
)

var grpcClient pb.ImageProcessingServiceClient

func InitGRPCClient() error {
	grpcServerAddress := os.Getenv("GRPC_SERVER_ADDRESS")

	if grpcServerAddress == "" {
		log.Err(ErrBadRequest).Msg("empty grpc server address")
		return ErrBadRequest
	}

	conn, err := grpc.NewClient(
		grpcServerAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Err(err).Msg("Failed to connect to gRPC server")
		return ErrInternal
	}
	grpcClient = pb.NewImageProcessingServiceClient(conn)

	return nil
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Err(ErrNotAllowed).Msg("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		log.Err(ErrBadRequest).Msg("Failed to read image")
		http.Error(w, "Failed to read image", http.StatusBadRequest)
		return
	}
	defer file.Close()

	imgBytes, err := io.ReadAll(file)
	if err != nil {
		log.Err(ErrInternal).Msg("Failed to read image data")
		http.Error(w, "Failed to read image data", http.StatusInternalServerError)
		return
	}

	var resp *pb.ImageResponse
	switch r.FormValue("type") {
	case "binary":
		resp, err = grpcClient.ConvertToBinary(context.Background(), &pb.BinaryRequest{
			ImageData: imgBytes,
		})
	case "monochrome":
		resp, err = grpcClient.ConvertToMonochrome(context.Background(), &pb.MonochromeRequest{
			ImageData:   imgBytes,
			TargetColor: r.FormValue("color"),
		})
	case "threshold":
		var thresholdInt int

		thresholdInt, err = strconv.Atoi(r.FormValue("threshold"))
		if err != nil {
			log.Err(ErrBadRequest).Msg("Invalid threshold value")
			http.Error(w, "Invalid threshold value", http.StatusBadRequest)
			return
		}

		resp, err = grpcClient.ConvertToThreshold(context.Background(), &pb.ThresholdRequest{
			ImageData: imgBytes,
			Threshold: int32(thresholdInt),
		})
	}
	if err != nil {
		log.Err(err).Msg("gRPC call failed")
		http.Error(w, "gRPC call failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Write(resp.ProcessedImageData)
}
