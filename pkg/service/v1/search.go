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

	apiv1 "github.com/johandry/finder2d/api/v1"
)

// Search implement the API method from the generated protobuf
func (s *Finder2DService) Search(ctx context.Context, req *apiv1.SearchRequest) (*apiv1.SearchResponse, error) {
	if s.finder.Source == nil {
		errMsg := "the Finder2D does not have a frame or source matrix, load the source matrix first"
		log.Printf("[ERROR] %s", errMsg)
		return nil, fmt.Errorf(errMsg)
	}
	if s.finder.Target == nil {
		errMsg := "the Finder2D does not have an image or target matrix, load the target matrix first"
		log.Printf("[ERROR] %s", errMsg)
		return nil, fmt.Errorf(errMsg)
	}
	if req.Percentage != 0 {
		s.finder.Percentage = float64(req.Percentage)
	}
	if req.Delta != 0 {
		s.finder.Delta = int(req.Delta)
	}

	if err := s.finder.SearchSimple(); err != nil {
		errMsg := fmt.Sprintf("failed to search the target matrix. %s", err)
		log.Printf("[ERROR] %s", errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	n := len(s.finder.Matches)
	log.Printf("[INFO] searched target matrix in source matrix with matching percentage %f and blurry delta %d, found %d matches", s.finder.Percentage, s.finder.Delta, n)

	return &apiv1.SearchResponse{
		Api:          apiVersion,
		TotalMatches: int32(n),
	}, nil
}
