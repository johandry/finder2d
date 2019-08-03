package v1

import (
	"context"

	apiv1 "github.com/johandry/finder2d/api/v1"
)

// LoadMatrix implement the API method from the generated protobuf
func (s *Finder2DService) LoadMatrix(ctx context.Context, req *apiv1.LoadMatrixRequest) (*apiv1.LoadMatrixResponse, error) {
	return &apiv1.LoadMatrixResponse{
		Api: apiVersion,
	}, nil
}
