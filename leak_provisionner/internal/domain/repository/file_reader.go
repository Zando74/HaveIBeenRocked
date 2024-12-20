package repository

import "github.com/Zando74/IHaveBeenRocked/leak_provisionner/internal/domain/entity"

type PersistenceHandler interface {
	StorePasswordBatch(password []entity.Password) error
	MarkAsFinished() error
}

type FileReader interface {
	ProcessFile(filePath string, persistenceHandler PersistenceHandler) error
}
