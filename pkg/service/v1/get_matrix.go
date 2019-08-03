package v1

import (
	"context"

	apiv1 "github.com/johandry/finder2d/api/v1"
)

// GetMatrix implement the API method from the generated protobuf
func (s *Finder2DService) GetMatrix(ctx context.Context, req *apiv1.GetMatrixRequest) (*apiv1.GetMatrixResponse, error) {
	m := s.finder.Source

	content := m.String()
	w, h := m.Size()

	matrix := &apiv1.Matrix{
		Content: content,
		Width:   int32(w),
		Height:  int32(h),
	}

	return &apiv1.GetMatrixResponse{
		Api:    apiVersion,
		Name:   apiv1.MatrixName_SOURCE,
		Matrix: matrix,
	}, nil
}
