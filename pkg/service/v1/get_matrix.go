/*
Copyright The Finder2D Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/johandry/finder2d"
	apiv1 "github.com/johandry/finder2d/api/v1"
)

// GetMatrix implement the API method from the generated protobuf
func (s *Finder2DService) GetMatrix(ctx context.Context, req *apiv1.GetMatrixRequest) (*apiv1.GetMatrixResponse, error) {
	if err := s.checkAPIVersion(req.Api); err != nil {
		return nil, err
	}

	log.Printf("[INFO] requesting %s matrix", strings.ToLower(req.Name.String()))

	var m *finder2d.Matrix
	switch req.Name {
	case apiv1.MatrixName_SOURCE:
		m = s.finder.Source
	case apiv1.MatrixName_TARGET:
		m = s.finder.Target
	}

	if m == nil {
		return nil, fmt.Errorf("matrix not found, load the matrix")
	}

	z, o := s.finder.Values()
	content := m.Sprintf(string([]byte{z}), string([]byte{o}))
	w, h := m.Size()

	log.Printf("[INFO] sending %s matrix (%d,%d)", strings.ToLower(req.Name.String()), w, h)

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
