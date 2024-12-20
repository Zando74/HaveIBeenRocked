package cli

import (
	"github.com/Zando74/IHaveBeenRocked/leak_provisionner/internal/domain/usecase"
	"github.com/Zando74/IHaveBeenRocked/leak_provisionner/internal/logger"
	"github.com/spf13/cobra"
)

var (
	log = logger.LoggerSingleton.GetInstance()
)

type LoadPasswordsFromFileCommand struct {
	LoadPasswordFromFileUseCase *usecase.LoadPasswordsFromFile
	Cmd                         *cobra.Command
}

func (c *LoadPasswordsFromFileCommand) Init() {

	filePath := ""

	c.Cmd = &cobra.Command{
		Use:   "load_passwords_from_file",
		Short: "Load passwords from file",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.LoadPasswordFromFileUseCase.Execute(filePath); err != nil {
				log.Error(err)
				return err
			}
			return nil
		},
	}

	c.Cmd.Flags().StringVarP(&filePath, "file", "f", "", "file path")
	if err := c.Cmd.MarkFlagRequired("file"); err != nil {
		log.Error(err)
	}
}

func (c *LoadPasswordsFromFileCommand) GetCommand() *cobra.Command {
	return c.Cmd
}
