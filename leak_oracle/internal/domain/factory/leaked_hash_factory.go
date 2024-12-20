package factory

import (
	"github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/domain/entity"
	"github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/domain/value_object"
)

type LeakedHashFactory struct{}

func (fc *LeakedHashFactory) Build(password []byte) (*entity.LeakedHash, error) {
	passwordHash, err := value_object.NewPasswordHashFromPassword(password)

	if err != nil {
		return nil, err
	}

	return &entity.LeakedHash{Hash: passwordHash}, nil
}
