package server

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"strconv"
	"strings"

	pb "github.com/SV1Stail/imageRrocessing/server/github.com/SV1Stail/imageRrocessing/gen"
	"github.com/disintegration/imaging"
	"github.com/rs/zerolog/log"
)

type Server struct {
	pb.UnimplementedImageProcessingServiceServer
}

func hexToRGBA(hex string) (color.RGBA, error) {
	log.Debug().Msg("start hexToRGBA")

	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return color.RGBA{}, ErrWrongFormat
	}

	r, err := strconv.ParseUint(hex[0:2], 16, 8)
	if err != nil {
		return color.RGBA{}, ErrWrongFormat
	}

	g, err := strconv.ParseUint(hex[2:4], 16, 8)
	if err != nil {
		return color.RGBA{}, ErrWrongFormat
	}

	b, err := strconv.ParseUint(hex[4:6], 16, 8)
	if err != nil {
		return color.RGBA{}, ErrWrongFormat
	}

	log.Debug().Msg("success hexToRGBA")

	return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}, nil
}

// Функция для применения цвета к изображению
func applyColor(img image.Image, targetColor color.Color) image.Image {
	log.Debug().Msg("start applyColor")

	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			originalColor := img.At(x, y)
			grayColor := color.GrayModel.Convert(originalColor).(color.Gray)
			r, g, b, a := targetColor.RGBA()
			grayValue := float64(grayColor.Y) / 255.0
			newColor := color.RGBA{
				R: uint8(float64(r>>8) * grayValue),
				G: uint8(float64(g>>8) * grayValue),
				B: uint8(float64(b>>8) * grayValue),
				A: uint8(a >> 8),
			}
			dst.Set(x, y, newColor)
		}
	}

	log.Debug().Msg("success hexToRGBA")

	return dst
}
func (s *Server) ConvertToMonochrome(ctx context.Context, req *pb.MonochromeRequest) (*pb.ImageResponse, error) {
	log.Info().Msg("start Convert To Monochrome")
	if req.TargetColor == "" {
		log.Err(ErrNoColor).Msg("Convert To Monochrome failed")

		return nil, ErrNoColor
	}
	if req.ImageData == nil {
		log.Err(ErrNoData).Msg("Convert To Monochrome failed")

		return nil, ErrNoData
	}

	targetColor, err := hexToRGBA(req.TargetColor)
	if err != nil {
		log.Err(err).Msg("Convert To Monochrome failed")

		return nil, err
	}

	img, format, err := image.Decode(bytes.NewReader(req.ImageData))
	if err != nil {
		log.Err(err).Msg("Convert To Monochrome failed")

		return nil, err
	}

	coloredImg := applyColor(img, targetColor)

	var buf bytes.Buffer
	switch format {
	case "jpeg":
		err = jpeg.Encode(&buf, coloredImg, nil)
	case "png":
		err = png.Encode(&buf, coloredImg)
	default:
		return nil, ErrWrongFormat
	}
	if err != nil {
		log.Err(err).Msg("Convert To Monochrome failed")

		return nil, err
	}

	log.Info().Msg("Convert To Monochrome success")

	return &pb.ImageResponse{
		ProcessedImageData: buf.Bytes(),
	}, nil
}

func (s *Server) ConvertToBinary(ctx context.Context, req *pb.BinaryRequest) (*pb.ImageResponse, error) {
	log.Info().Msg("start ConvertToBinary")

	if req.ImageData == nil {
		log.Err(ErrNoData).Msg("ConvertToBinary failed")

		return nil, ErrNoData
	}

	img, format, err := image.Decode(bytes.NewReader(req.ImageData))
	if err != nil {
		return nil, err
	}

	imgGray := imaging.Grayscale(img)

	var buf bytes.Buffer
	switch format {
	case "jpeg":
		err = jpeg.Encode(&buf, imgGray, nil)
	case "png":
		err = png.Encode(&buf, imgGray)
	default:
		return nil, ErrWrongFormat
	}
	if err != nil {
		log.Err(err).Msg("ConvertToBinary failed")

		return nil, err
	}

	log.Info().Msg("ConvertToBinary success")

	return &pb.ImageResponse{
		ProcessedImageData: buf.Bytes(),
	}, nil
}

func (s *Server) ConvertToThreshold(ctx context.Context, req *pb.ThresholdRequest) (*pb.ImageResponse, error) {
	log.Info().Msg("start ConvertToThreshold")

	if req.ImageData == nil {
		log.Err(ErrNoData).Msg("ConvertToThreshold failed")

		return nil, ErrNoData
	}
	if req.Threshold < 0 || req.Threshold > 255 {
		log.Err(ErrWrongFormat).Msg("ConvertToThreshold failed")

		return nil, ErrWrongFormat
	}

	img, format, err := image.Decode(bytes.NewReader(req.ImageData))
	if err != nil {
		log.Err(ErrWrongFormat).Msg("ConvertToThreshold failed")
		return nil, err
	}

	bounds := img.Bounds()
	result := image.NewRGBA(bounds)

	threshold := uint8(req.Threshold)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			originalColor := img.At(x, y)
			r, g, b, _ := originalColor.RGBA()
			gray := uint8((0.299*float64(r>>8) + 0.587*float64(g>>8) + 0.114*float64(b>>8)))
			if gray >= threshold {
				result.Set(x, y, color.RGBA{R: 255, G: 255, B: 255, A: 255})
			} else {
				result.Set(x, y, color.RGBA{R: 0, G: 0, B: 0, A: 255})
			}
		}
	}

	var buf bytes.Buffer
	switch format {
	case "jpeg":
		err = jpeg.Encode(&buf, result, nil)
	case "png":
		err = png.Encode(&buf, result)
	default:
		return nil, ErrWrongFormat
	}
	if err != nil {
		log.Err(err).Msg("ConvertToBinary failed")

		return nil, err
	}

	log.Info().Msg("ConvertToThreshold success")

	return &pb.ImageResponse{
		ProcessedImageData: buf.Bytes(),
	}, nil

}
