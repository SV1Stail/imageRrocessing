package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	pb "github.com/SV1Stail/imageRrocessing/client/imageRrocessing/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	pic = "new.png"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewImageProcessingServiceClient(conn)

	imageData, err := os.ReadFile(pic)
	if err != nil {
		log.Fatalf("failed to read image: %v", err)
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		response, err := client.ConvertToMonochrome(context.Background(), &pb.MonochromeRequest{
			ImageData:   imageData,
			TargetColor: "#FFFACD",
		})
		if err != nil {
			log.Fatalf("failed to ConvertToMonochrome image: %v", err)
		}

		err = os.WriteFile("monochrome.png", response.ProcessedImageData, 0644)
		if err != nil {
			log.Fatalf("failed to save image: %v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		response, err := client.ConvertToBinary(context.Background(), &pb.BinaryRequest{
			ImageData: imageData,
		})
		if err != nil {
			log.Fatalf("failed to ConvertToBinary image: %v", err)
		}

		err = os.WriteFile("binary.png", response.ProcessedImageData, 0644)
		if err != nil {
			log.Fatalf("failed to save image: %v", err)
		}

	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		resp, err := client.ConvertToThreshold(context.Background(), &pb.ThresholdRequest{
			ImageData: imageData,
			Threshold: 150,
		})
		if err != nil {
			log.Fatalf("failed to ConvertToThreshold image: %v", err)
		}

		err = os.WriteFile("threshold.png", resp.ProcessedImageData, 0644)
		if err != nil {
			log.Fatalf("failed to save image: %v", err)
		}
	}()

	wg.Wait()

	fmt.Println("Image processed and saved as output_image.jpg")
}
