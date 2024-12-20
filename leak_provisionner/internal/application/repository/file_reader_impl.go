package repository

import (
	"bufio"
	"os"

	"github.com/Zando74/IHaveBeenRocked/leak_provisionner/internal/config"
	"github.com/Zando74/IHaveBeenRocked/leak_provisionner/internal/domain/entity"
	domain_repo "github.com/Zando74/IHaveBeenRocked/leak_provisionner/internal/domain/repository"
	"github.com/Zando74/IHaveBeenRocked/leak_provisionner/internal/logger"
)

var (
	log = logger.LoggerSingleton.GetInstance()
	cfg = config.ConfigSingleton.GetInstance()
)

type FileReaderImpl struct{}

func (f *FileReaderImpl) processBatch(passwordBatch *[]entity.Password, persistenceHandler domain_repo.PersistenceHandler) error {
	log.Debug(logger.SendingBatchOfPasswords)
	err := persistenceHandler.StorePasswordBatch(*passwordBatch)
	if err != nil {
		return err
	}
	return nil
}

func (f *FileReaderImpl) ProcessFile(filePath string, persistenceHandler domain_repo.PersistenceHandler) error {

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	passwordBatch := []entity.Password{}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		password := scanner.Text()

		passwordBatch = append(passwordBatch, entity.Password(password))
		if len(passwordBatch) >= cfg.FileReader.PasswordBatchSize {
			err := f.processBatch(&passwordBatch, persistenceHandler)
			if err != nil {
				return err
			}
			passwordBatch = passwordBatch[:0]
		}
	}

	if len(passwordBatch) > 0 {
		err := f.processBatch(&passwordBatch, persistenceHandler)
		if err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		log.Error(err)
		return err
	}

	return persistenceHandler.MarkAsFinished()
}
