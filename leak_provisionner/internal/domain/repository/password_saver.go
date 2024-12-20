package repository

import "github.com/Zando74/IHaveBeenRocked/leak_provisionner/internal/domain/entity"

type PasswordSaver interface {
	StorePasswordBatch(passwordBatch []entity.Password) error
	MarkAsFinished() error
}
