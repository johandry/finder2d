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
	"log"
	"strings"

	apiv1 "github.com/johandry/finder2d/api/v1"
)

// LoadMatrix implement the API method from the generated protobuf
func (s *Finder2DService) LoadMatrix(ctx context.Context, req *apiv1.LoadMatrixRequest) (*apiv1.LoadMatrixResponse, error) {
	r := strings.NewReader(req.Matrix.Content)
	var err error
	var w, h int
	switch req.Name {
	case apiv1.MatrixName_SOURCE:
		err = s.finder.LoadSource(r)
		w, h = s.finder.Source.Size()
	case apiv1.MatrixName_TARGET:
		err = s.finder.LoadTarget(r)
		w, h = s.finder.Target.Size()
	}
	if err != nil {
		return nil, err
	}

	log.Printf("[INFO] %s matrix (%d,%d) loaded", strings.ToLower(req.Name.String()), w, h)

	return &apiv1.LoadMatrixResponse{
		Api: apiVersion,
	}, nil
}
