package usecase

import (
	"github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/domain/entity"
	"github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/domain/factory"
	"github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/domain/repository"
	"github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/logger"
)

type HashPasswordBatch struct {
	LeakedHashRepository repository.HashLeakedRepository
	LeakedHashFactory    *factory.LeakedHashFactory
}

func (uc *HashPasswordBatch) Execute(passwords [][]byte) ([]*entity.LeakedHash, error) {

	log.Debug(logger.HashBatchOfPasswordsMessage)

	leakedHashes := []*entity.LeakedHash{}

	for _, password := range passwords {
		leakedHash, err := uc.LeakedHashFactory.Build(password)

		if err != nil {
			return nil, err
		}

		leakedHashes = append(leakedHashes, leakedHash)
	}

	return leakedHashes, nil
}
