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
