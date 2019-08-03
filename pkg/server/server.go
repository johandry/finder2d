package server

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/johandry/finder2d"
	apiv1 "github.com/johandry/finder2d/api/v1"
	servicev1 "github.com/johandry/finder2d/pkg/service/v1"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server is the server that expose a gRPC and REST API
type Server struct {
	host   string
	port   string
	finder *finder2d.Finder2D
	ctx    context.Context

	// Channels
	errCh  chan error
	sigCh  chan os.Signal
	doneCh chan struct{}

	// Servers
	grpcServer *grpc.Server
	httpServer *http.Server
}

// Serve starts serving
func Serve(port, sourceFileName, zero, one string) error {
	s := &Server{
		host: "localhost",
		port: port,
	}

	if err := s.newFinder2D(sourceFileName, zero, one); err != nil {
		return err
	}

	s.ctx = context.Background()
	s.sigCh = make(chan os.Signal, 1)
	s.errCh = make(chan error, 1)

	if err := s.Start(); err != nil {
		return err
	}

	return s.Wait()
}

func (s *Server) newFinder2D(sourceFileName, zero, one string) error {
	if len(sourceFileName) == 0 {
		return nil
	}

	sourceFile, err := os.Open(sourceFileName)
	if err != nil {
		return fmt.Errorf("fail to open the frame file %q. %s", sourceFileName, err)
	}

	s.finder = finder2d.New([]byte(one)[0], []byte(zero)[0], 0, 0)

	if err := s.finder.LoadSource(sourceFile); err != nil {
		return fmt.Errorf("fail to load the source file %q. %s", sourceFileName, err)
	}

	return nil
}

func (s *Server) makeSignalCh() {
	if s.sigCh == nil {
		s.sigCh = make(chan os.Signal, 1)
	}
	signal.Notify(s.sigCh, os.Interrupt)
}

// Start starts the server to server/expose a gRPC & REST API
func (s *Server) Start() error {
	serveAddress := fmt.Sprintf("%s:%s", s.host, s.port)

	opts := []grpc.ServerOption{}
	s.grpcServer = grpc.NewServer(opts...)

	log.Printf("[DEBUG] registering service for gRPC")
	service := servicev1.New(s.finder)
	service.Register(s.grpcServer)

	reflection.Register(s.grpcServer)

	ctx, cancel := context.WithCancel(s.ctx)
	s.ctx = ctx

	gwopts := []grpc.DialOption{grpc.WithInsecure()}

	gwmux := runtime.NewServeMux(runtime.WithMarshalerOption(
		runtime.MIMEWildcard,
		&runtime.JSONPb{OrigName: true, EmitDefaults: true},
	))

	log.Printf("[DEBUG] registering service for HTTP/REST")
	if err := apiv1.RegisterFinder2DHandlerFromEndpoint(s.ctx, gwmux, serveAddress, gwopts); err != nil {
		log.Printf("[ERROR] failed to register service. %s", err)
		defer cancel()
		return s.Stop()
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/swagger/finder2d.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		io.Copy(w, bytes.NewReader(apiv1.Swagger))
	})

	mux.Handle("/", gwmux)

	httpAddress := fmt.Sprintf(":%s", s.port)

	conn, err := net.Listen("tcp", httpAddress)
	if err != nil {
		log.Printf("[ERROR] failed to create connection on ':%s'. %s", s.port, err)
		defer cancel()
		return s.Stop()
	}

	s.httpServer = &http.Server{
		Addr:    serveAddress,
		Handler: muxHandlerFuncInsecure(s.grpcServer, mux),
	}

	s.makeSignalCh()

	log.Printf("[INFO] starting HTTP/REST gateway and gRPC server on %s...", s.httpServer.Addr)

	go func() {
		defer cancel()
		err := s.httpServer.Serve(conn)
		log.Printf("[WARN] the HTTP/REST gateway and gRPC server stoped serving. Error: %s", err)
		s.errCh <- err
	}()

	return nil
}

// Wait keep the server running until it's stop or fail
func (s *Server) Wait() error {
	for {
		select {
		case <-s.ctx.Done():
			return s.ctx.Err()
		case e := <-s.errCh:
			log.Printf("[ERROR] %s", e)
			ctx, cancel := context.WithCancel(s.ctx)
			s.ctx = ctx
			cancel()
		case <-s.sigCh:
			log.Printf("[WARN] received a ^C signal, shutting down the servers")
			s.Stop()
		}
	}
}

// Stop stops the server running
func (s *Server) Stop() error {
	timeout := 5 * time.Second
	log.Printf("[WARN] shutting down Multiplex server in %s seconds...", timeout)
	ctx, cancel := context.WithTimeout(s.ctx, timeout)
	defer cancel()
	s.ctx = ctx

	return s.httpServer.Shutdown(s.ctx)
}

// the gRPC server or the HTTP muxer (HTTP/REST) at runtime. This function works without TLS
func muxHandlerFuncInsecure(grpcServer *grpc.Server, httpHandler http.Handler) http.Handler {
	return h2c.NewHandler(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
				grpcServer.ServeHTTP(w, r)
			} else {
				httpHandler.ServeHTTP(w, r)
			}
		}), &http2.Server{})
}
