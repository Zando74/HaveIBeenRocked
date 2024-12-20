package controller

import (
	"github.com/Zando74/IHaveBeenRocked/leak_provisionner/internal/application/controller/cli"
	"github.com/Zando74/IHaveBeenRocked/leak_provisionner/internal/application/repository"
	"github.com/Zando74/IHaveBeenRocked/leak_provisionner/internal/domain/usecase"
	"github.com/Zando74/IHaveBeenRocked/leak_provisionner/internal/logger"
	"github.com/spf13/cobra"
)

var (
	log = logger.LoggerSingleton.GetInstance()
)

type Command interface {
	Init()
	GetCommand() *cobra.Command
}

func RegisterCommands(rootCmd *cobra.Command, commands ...Command) {
	for _, cmd := range commands {
		cmd.Init()
		rootCmd.AddCommand(cmd.GetCommand())
	}
}

type Controller struct {
	Commands []Command
}

func (c *Controller) RunCli() {
	rootCmd := &cobra.Command{Use: "leak_provisionner"}

	passwordSaverImpl, err := repository.NewPasswordSaverGRPCImpl()
	if err != nil {
		log.Fatal(err)
	}

	loadPwdCmd := &cli.LoadPasswordsFromFileCommand{
		LoadPasswordFromFileUseCase: &usecase.LoadPasswordsFromFile{
			FileReader:    &repository.FileReaderImpl{},
			PasswordSaver: passwordSaverImpl,
		},
	}

	RegisterCommands(rootCmd, loadPwdCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
