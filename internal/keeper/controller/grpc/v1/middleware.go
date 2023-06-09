package v1

import (
	"context"
	"regexp"

	"github.com/alkurbatov/goph-keeper/internal/keeper/entity"
	"github.com/alkurbatov/goph-keeper/internal/keeper/infra/logger"
	"github.com/alkurbatov/goph-keeper/internal/libraries/creds"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var methodsWithoutAuth = regexp.MustCompile(`/(Login|Register)`)

// LoggingUnaryInterceptor is gRPC unary server interceptor
// which logs incoming requests and responses.
func LoggingUnaryInterceptor(log *logger.Logger) grpc.UnaryServerInterceptor {
	interceptor := func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		l := log.With().
			Str("req-id", uuid.NewV4().String()).
			Logger()

		l.Info().
			Str("method", info.FullMethod).
			Msg("")

		resp, err := handler(l.WithContext(ctx), req)

		status, ok := status.FromError(err)
		if ok {
			l.Info().
				Str("status", status.Code().String()).
				Msg("")

			return resp, err
		}

		l.Info().
			Err(err).
			Msg("")

		return resp, err
	}

	return interceptor
}

// AuthUnaryInterceptor is gRPC unary server interceptor extracts access token
// from metadata and verifies it.
// If the token is valid, request is passed further.
// Token's subject ID is injected as user ID into the context to use later.
func AuthUnaryInterceptor(secret creds.Password) grpc.UnaryServerInterceptor {
	interceptor := func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		if methodsWithoutAuth.MatchString(info.FullMethod) {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, entity.ErrInvalidCredentials.Error())
		}

		values := md.Get("authorization")
		if len(values) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, entity.ErrInvalidCredentials.Error())
		}

		// NB (alkurbatov): It is ok to pass empty token further, Decode will mark it as invalid anyway.
		claims, err := entity.TokenFromString(values[0]).Decode(secret)
		if err != nil {
			logger.FromContext(ctx).Error().Err(err).Msg("Unauthorized access")

			return nil, status.Errorf(codes.Unauthenticated, entity.ErrInvalidCredentials.Error())
		}

		userID, err := uuid.FromString(claims.Subject)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}

		user := entity.User{
			ID:       userID,
			Username: claims.Username,
		}

		return handler(user.WithContext(ctx), req)
	}

	return interceptor
}
