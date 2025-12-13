package server

import (
	"context"
	"errors"
	"strings"

	"github.com/shivamp1998/vpn_backend/internal/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func AuthInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	if isPublicEndpoint(info.FullMethod) {
		return handler(ctx, req)
	}

	token, err := extractTokenFromMetadata(ctx)

	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication required: %v", err)
	}

	claims, err := auth.ValidateToken(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	ctx = context.WithValue(ctx, auth.UserIdKey, claims.UserId)
	ctx = context.WithValue(ctx, auth.UserEmailKey, claims.Email)

	return handler(ctx, req)

}

func isPublicEndpoint(method string) bool {
	publicEndpoints := []string{
		"/vpn.UserService/Register",
		"/vpn.UserService/Login",
	}

	for _, endpoint := range publicEndpoints {
		if method == endpoint {
			return true
		}
	}

	return false
}

func extractTokenFromMetadata(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return "", errors.New("metadata not found")
	}

	authHeaders := md.Get("authorization")

	if len(authHeaders) == 0 {
		return "", errors.New("authorization headers missing")
	}

	authHeader := authHeaders[0]

	parts := strings.Split(authHeader, " ")

	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("invalid authorization header fromat")
	}

	return parts[1], nil
}
