package rpc_api

import (
	"context"
	"github.com/RafalSalwa/interview-app-srv/pkg/models"
	pb "github.com/RafalSalwa/interview-app-srv/proto/grpc"
	"go.opentelemetry.io/otel"
	otelcodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (authServer *AuthServer) SignInUser(ctx context.Context, req *pb.SignInUserInput) (*pb.SignInUserResponse, error) {
	loginUser := &models.LoginUserRequest{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
	}

	ur, err := authServer.authService.SignInUser(loginUser)
	if err != nil {
		authServer.logger.Error().Err(err).Msg("rpc:service:signin")
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	res := &pb.SignInUserResponse{
		Status:       "success",
		AccessToken:  ur.Token,
		RefreshToken: ur.RefreshToken,
	}
	return res, nil
}

func (authServer *AuthServer) SignUpUser(ctx context.Context, req *pb.SignUpUserInput) (*pb.SignUpUserResponse, error) {
	ctx, span := otel.GetTracerProvider().Tracer("auth_service-rpc").Start(ctx, "GRPC SignUpUser")
	defer span.End()

	signUpUser := &models.CreateUserRequest{
		Email:           req.GetEmail(),
		Password:        req.GetPassword(),
		PasswordConfirm: req.GetPasswordConfirm(),
	}

	ur, err := authServer.authService.SignUpUser(ctx, signUpUser)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelcodes.Error, err.Error())
		authServer.logger.Error().Err(err).Msg("rpc:service:signup")
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	res := &pb.SignUpUserResponse{
		Id:                ur.Id,
		Username:          ur.Username,
		VerificationToken: ur.VerificationCode,
		CreatedAt:         nil,
	}
	return res, nil
}

func (authServer *AuthServer) GetVerificationKey(ctx context.Context, in *pb.VerificationCodeRequest) (*pb.VerificationCodeResponse, error) {
	ur, err := authServer.authService.GetVerificationKey(ctx, in.Email)
	if err != nil {
		authServer.logger.Error().Err(err).Msg("rpc:service:getkey")
		return nil, err
	}
	return &pb.VerificationCodeResponse{Code: ur.VerificationCode}, nil
}
