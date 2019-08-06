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

// GetMatch implement the API method from the generated protobuf
func (s *Finder2DService) GetMatch(ctx context.Context, req *apiv1.GetMatchRequest) (*apiv1.GetMatchResponse, error) {
	if req.Id < 0 || int(req.Id) >= len(s.finder.Matches) {
		errMsg := fmt.Sprintf("not found match with id=%d", req.Id)
		log.Printf("[ERROR] %s", errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	match := s.finder.Matches[int(req.Id)]
	m := &apiv1.Match{
		X:          int32(match.X),
		Y:          int32(match.Y),
		Percentage: float32(match.Percentage),
	}

	targetW, targetH := s.finder.Target.Size()
	matrix := s.finder.Source.Sample(match.X, match.Y, targetW, targetH)
	z, o := s.finder.Values()
	content := matrix.Sprintf(string([]byte{z}), string([]byte{o}))
	matx := &apiv1.Matrix{
		Width:   int32(targetW),
		Height:  int32(targetH),
		Content: content,
	}

	log.Printf("[INFO] match id=%d requested and returned", req.Id)

	return &apiv1.GetMatchResponse{
		Api:    apiVersion,
		Match:  m,
		Matrix: matx,
	}, nil
}
