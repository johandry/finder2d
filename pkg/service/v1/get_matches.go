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

	apiv1 "github.com/johandry/finder2d/api/v1"
)

// GetMatches implement the API method from the generated protobuf
func (s *Finder2DService) GetMatches(ctx context.Context, req *apiv1.GetMatchesRequest) (*apiv1.GetMatchesResponse, error) {
	ms := []*apiv1.Match{}

	for _, match := range s.finder.Matches {
		m := &apiv1.Match{
			X:          int32(match.X),
			Y:          int32(match.Y),
			Percentage: float32(match.Percentage),
		}
		ms = append(ms, m)
	}

	log.Printf("[INFO] list of matches requested, returned %d matches", len(ms))

	return &apiv1.GetMatchesResponse{
		Api:     apiVersion,
		Matches: ms,
	}, nil
}
