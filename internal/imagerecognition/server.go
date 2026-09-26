package imagerecognition

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "landan-desktop-fyne/proto/ir"
)

// Server implements pb.ImageRecognitionServiceServer on top of a Recognizer.
type Server struct {
	pb.UnimplementedImageRecognitionServiceServer
	recognizer Recognizer
}

func NewServer(oRecognizer Recognizer) *Server {
	return &Server{recognizer: oRecognizer}
}

func (s *Server) Recognize(oCtx context.Context, oRequest *pb.RecognizeRequest) (*pb.RecognizeResponse, error) {
	if len(oRequest.GetImage()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "image is empty")
	}
	oResult, err := s.recognizer.Recognize(oCtx, oRequest.GetImage())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "recognize: "+err.Error())
	}

	aLabels := make([]*pb.Label, 0, len(oResult.Labels))
	for _, oLabel := range oResult.Labels {
		aLabels = append(aLabels, &pb.Label{Name: oLabel.Name, Confidence: oLabel.Confidence})
	}
	return &pb.RecognizeResponse{
		Labels: aLabels,
		Width:  int32(oResult.Width),
		Height: int32(oResult.Height),
		Format: oResult.Format,
	}, nil
}
