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
