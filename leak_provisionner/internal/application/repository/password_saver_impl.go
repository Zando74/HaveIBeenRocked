package repository

import (
	"context"
	"fmt"

	"github.com/Zando74/IHaveBeenRocked/leak_provisionner/internal/application/grpc_proto"
	"github.com/Zando74/IHaveBeenRocked/leak_provisionner/internal/config"
	"github.com/Zando74/IHaveBeenRocked/leak_provisionner/internal/domain/entity"
	"github.com/Zando74/IHaveBeenRocked/leak_provisionner/internal/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PasswordSaverGRPCImpl struct {
	Client grpc_proto.RawPasswordListUploadClient
	Stream grpc.ClientStreamingClient[grpc_proto.RawPasswordList, grpc_proto.Status]
}

func NewPasswordSaverGRPCImpl() (*PasswordSaverGRPCImpl, error) {
	cfg := config.ConfigSingleton.GetInstance()
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", cfg.GrpcServer.Host, cfg.GrpcServer.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := grpc_proto.NewRawPasswordListUploadClient(conn)

	return &PasswordSaverGRPCImpl{
		Client: client,
		Stream: nil,
	}, nil
}

func (ps *PasswordSaverGRPCImpl) openStream() error {

	stream, err := ps.Client.UploadRawPasswordList(context.Background())
	if err != nil {
		return err
	}

	ps.Stream = stream

	return nil
}

func (ps *PasswordSaverGRPCImpl) closeStream() error {

	if ps.Stream == nil {
		return nil
	}
	status, err := ps.Stream.CloseAndRecv()

	if err != nil {
		return err
	}

	if !status.Success {
		log.Error(logger.StreamCloseErrorMessage, status.Message)
	} else {
		log.Info(logger.StreamClosedMessage)
	}

	return err
}

func (ps *PasswordSaverGRPCImpl) StorePasswordBatch(passwords []entity.Password) error {

	if ps.Stream == nil {
		err := ps.openStream()
		if err != nil {
			return err
		}
	}

	passwordList := make([][]byte, len(passwords))
	for i, password := range passwords {
		passwordList[i] = password
	}

	if err := ps.Stream.Send(&grpc_proto.RawPasswordList{Passwords: passwordList}); err != nil {
		return err
	}
	return nil
}

func (ps *PasswordSaverGRPCImpl) MarkAsFinished() error {
	return ps.closeStream()
}
