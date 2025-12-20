package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/rs/cors"
	gen "github.com/shivamp1998/vpn_backend/proto/gen"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"connectrpc.com/connect"
	genconnect "github.com/shivamp1998/vpn_backend/proto/gen/genconnect"
)

func StartConnectServer(mainServer *Server) {
	connectPort := os.Getenv("CONNECT_PORT")

	if connectPort == "" {
		connectPort = ":50052"
	}

	userServiceHandler := &connectUserServiceHandler{server: mainServer}
	configServiceHandler := &connectConfigServiceHandler{server: mainServer}
	serverServiceHandler := &connectServerServiceHandler{server: mainServer}

	mux := http.NewServeMux()

	userServicePath, userServiceHTTPHandler := genconnect.NewUserServiceHandler(
		userServiceHandler,
		connect.WithInterceptors(newConnectAuthInterceptor()),
	)
	mux.Handle(userServicePath, userServiceHTTPHandler)

	serverServicePath, serverServiceHTTPHandler := genconnect.NewServerServiceHandler(
		serverServiceHandler,
		connect.WithInterceptors(newConnectAuthInterceptor()),
	)
	mux.Handle(serverServicePath, serverServiceHTTPHandler)

	configServicePath, configServiceHTTPHandler := genconnect.NewConfigServiceHandler(
		configServiceHandler,
		connect.WithInterceptors(newConnectAuthInterceptor()),
	)
	mux.Handle(configServicePath, configServiceHTTPHandler)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	handler := c.Handler(mux)
	http2Server := &http2.Server{}
	httpServer := &http.Server{
		Addr:    connectPort,
		Handler: h2c.NewHandler(handler, http2Server),
	}

	log.Printf("Connect RPC server listening on %s", connectPort)
	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatalf("Connect server failed: %v", err)
	}
}

func newConnectAuthInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			procedure := req.Spec().Procedure

			if IsPublicEndpoint(procedure) {
				return next(ctx, req)
			}

			authHeader := req.Header().Get("authorization")

			token, err := ExtractTokenFromHeader(authHeader)

			if err != nil {
				return nil, connect.NewError(
					connect.CodeUnauthenticated,
					err,
				)
			}

			ctx, err = ValidateAndSetUserContext(ctx, token)
			if err != nil {
				return nil, connect.NewError(
					connect.CodeUnauthenticated,
					err,
				)
			}

			return next(ctx, req)

		}
	}
}

type connectUserServiceHandler struct {
	server *Server
}

func (h *connectUserServiceHandler) Login(
	ctx context.Context,
	req *connect.Request[gen.LoginRequest],
) (*connect.Response[gen.AuthenticationResponse], error) {
	resp, err := h.server.Login(ctx, req.Msg)

	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (h *connectUserServiceHandler) Register(
	ctx context.Context,
	req *connect.Request[gen.RegisterRequest],
) (*connect.Response[gen.AuthenticationResponse], error) {
	resp, err := h.server.Register(ctx, req.Msg)

	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

type connectServerServiceHandler struct {
	server *Server
}

func (h *connectServerServiceHandler) CreateServer(
	ctx context.Context,
	req *connect.Request[gen.CreateServerRequest],
) (*connect.Response[gen.CreateServerResponse], error) {
	resp, err := h.server.CreateServer(ctx, req.Msg)

	if err != nil {
		return nil, err
	}

	return connect.NewResponse(resp), nil
}

func (h *connectServerServiceHandler) ListServers(
	ctx context.Context,
	req *connect.Request[gen.ListServerRequest],
) (*connect.Response[gen.ListServerResponse], error) {
	resp, err := h.server.ListServers(ctx, req.Msg)

	if err != nil {
		return nil, err
	}

	return connect.NewResponse(resp), nil
}

func (h *connectServerServiceHandler) GetServer(
	ctx context.Context,
	req *connect.Request[gen.GetServerRequest],
) (*connect.Response[gen.GetServerResponse], error) {
	resp, err := h.server.GetServer(ctx, req.Msg)

	if err != nil {
		return nil, err
	}

	return connect.NewResponse(resp), nil
}

type connectConfigServiceHandler struct {
	server *Server
}

func (h *connectConfigServiceHandler) GenerateConfig(
	ctx context.Context,
	req *connect.Request[gen.GenerateConfigRequest],
) (*connect.Response[gen.GenerateConfigResponse], error) {
	resp, err := h.server.GenerateConfig(ctx, req.Msg)

	if err != nil {
		return nil, err
	}

	return connect.NewResponse(resp), nil
}

func (h *connectConfigServiceHandler) GetConfig(
	ctx context.Context,
	req *connect.Request[gen.GetConfigRequest],
) (*connect.Response[gen.GetConfigResponse], error) {
	resp, err := h.server.GetConfig(ctx, req.Msg)

	if err != nil {
		return nil, err
	}

	return connect.NewResponse(resp), nil
}

func (h *connectConfigServiceHandler) RotateKeys(
	ctx context.Context,
	req *connect.Request[gen.GenerateConfigRequest],
) (*connect.Response[gen.GetConfigResponse], error) {
	// TODO: Implement RotateKeys in grpc_server.go first
	return nil, connect.NewError(
		connect.CodeUnimplemented,
		errors.New("RotateKeys is not yet implemented"),
	)
}
