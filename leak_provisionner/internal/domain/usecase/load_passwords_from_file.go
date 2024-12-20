package usecase

import (
	"github.com/Zando74/IHaveBeenRocked/leak_provisionner/internal/domain/repository"
	"github.com/Zando74/IHaveBeenRocked/leak_provisionner/internal/logger"
)

var (
	log = logger.LoggerSingleton.GetInstance()
)

type LoadPasswordsFromFile struct {
	FileReader    repository.FileReader
	PasswordSaver repository.PasswordSaver
}

func (l *LoadPasswordsFromFile) Execute(filePath string) error {

	log.GetInstance().Info(logger.LoadingPasswordFromFileMessage)

	return l.FileReader.ProcessFile(filePath, l.PasswordSaver)
}
