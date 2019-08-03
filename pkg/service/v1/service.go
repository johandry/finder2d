package v1

import (
	"github.com/johandry/finder2d"
	apiv1 "github.com/johandry/finder2d/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const apiVersion = "v1"

// Finder2DService defines the Finder2D service
type Finder2DService struct {
	finder *finder2d.Finder2D
}

// New create a new service
func New(f *finder2d.Finder2D) *Finder2DService {
	return &Finder2DService{
		finder: f,
	}
}

// Register registers this service to the given gRPC server
func (s *Finder2DService) Register(server *grpc.Server) {
	apiv1.RegisterFinder2DServer(server, s)
}

func (s *Finder2DService) checkAPIVersion(version string) error {
	if len(version) != 0 && version != apiVersion {
		return status.Errorf(codes.Unimplemented, "API version %q is not supported. This services implements API version %q", version, apiVersion)
	}
	return nil
}
