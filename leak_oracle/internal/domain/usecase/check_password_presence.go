package usecase

import (
	"github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/domain/repository"
	"github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/domain/value_object"
	"github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/logger"
)

var (
	log = logger.LoggerSingleton.GetInstance()
)

type CheckPasswordPresence struct {
	LeakedHashRepository repository.HashLeakedRepository
}

func (uc *CheckPasswordPresence) Execute(password string) (bool, error) {

	log.Debug(logger.CheckingPasswordPresenceMessage)

	passwordHash, err := value_object.NewPasswordHashFromPassword([]byte(password))

	if err != nil {
		return false, err
	}

	hash, err := uc.LeakedHashRepository.Retrieve(passwordHash)

	if err != nil {
		return false, err
	}

	return hash != nil, nil
}
