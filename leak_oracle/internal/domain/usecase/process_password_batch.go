package usecase

import (
	"github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/domain/repository"
	"github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/logger"
)

type ProcessPasswordBatch struct {
	HashPasswordBatchUseCase *HashPasswordBatch
	HashLeakedRepository     repository.HashLeakedRepository
}

func (uc *ProcessPasswordBatch) Execute(passwords [][]byte) error {

	log.Debug(logger.ProcessBatchPasswordsMessage)

	leakedHashes, err := uc.HashPasswordBatchUseCase.Execute(passwords)
	if err != nil {
		return err
	}

	err = uc.HashLeakedRepository.SaveBatch(leakedHashes)

	if err != nil {
		return err
	}

	return nil
}
