package cli_test

import (
	"github.com/Zando74/IHaveBeenRocked/leak_provisionner/internal/application/controller"
	"github.com/Zando74/IHaveBeenRocked/leak_provisionner/internal/application/controller/cli"
	"github.com/Zando74/IHaveBeenRocked/leak_provisionner/internal/application/grpc_proto"
	"github.com/Zando74/IHaveBeenRocked/leak_provisionner/internal/application/repository"
	"github.com/Zando74/IHaveBeenRocked/leak_provisionner/internal/domain/usecase"
	"github.com/Zando74/IHaveBeenRocked/leak_provisionner/internal/integration/mock"
	"github.com/spf13/cobra"
)

var (
	PasswordGrpcServerMock = &mock.GrpcServerMock{
		DB: make(map[string]bool),
	}
)

func InitTestCli() *cobra.Command {

	conn, _ := mock.InitTestGrpcServerConnection(PasswordGrpcServerMock)

	client := grpc_proto.NewRawPasswordListUploadClient(conn)

	rootCmd := &cobra.Command{Use: "leak_provisionner"}

	loadPwdCmd := &cli.LoadPasswordsFromFileCommand{
		LoadPasswordFromFileUseCase: &usecase.LoadPasswordsFromFile{
			FileReader: &repository.FileReaderImpl{},
			PasswordSaver: &repository.PasswordSaverGRPCImpl{
				Client: client,
				Stream: nil,
			},
		},
	}

	controller.RegisterCommands(rootCmd, loadPwdCmd)

	return rootCmd
}
