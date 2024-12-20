package repository

import (
	"github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/domain/entity"
	"github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/domain/value_object"
)

type HashLeakedRepository interface {
	SaveBatch(hashes []*entity.LeakedHash) error
	Retrieve(hash value_object.PasswordHash) (*entity.LeakedHash, error)
}
