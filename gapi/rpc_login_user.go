package gapi

import (
	"context"
	"database/sql"

	db "github.com/giadat1599/small_bank/db/sqlc"
	"github.com/giadat1599/small_bank/pb"
	"github.com/giadat1599/small_bank/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	user, err := server.store.GetUser(ctx, req.GetUsername())

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "username not found: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to find user: %s", err)
	}

	err = utils.CheckPassword(req.GetPassword(), user.HashedPassword)

	if err != nil {
		return nil, status.Errorf(codes.NotFound, "incorrect password: %s", err)
	}

	acessToken, accessPayload ,err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenLifeTime)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create token: %s", err)
	}

	refreshToken, refreshPayload , err := server.tokenMaker.CreateToken(user.Username, server.config.RefreshTokenLifeTime)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create refresh token: %s", err)
	}

	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID           : refreshPayload.ID,
		Username     : user.Username,
		RefreshToken : refreshToken,
		UserAgent    : "",
		ClientIp     : "",
		IsBlocked    : false,
		ExpiresAt    : refreshPayload.ExpiredAt,
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session: %s", err)
	}

	resp := &pb.LoginUserResponse{
		User: convertUser(user),
		SessionId: session.ID.String(),
		AccessToken: acessToken,
		RefreshToken: refreshToken,
		AccessTokenExpiresAt: timestamppb.New(accessPayload.ExpiredAt),
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
	}
	return resp, nil 
}