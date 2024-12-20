package repository

import (
	"strings"
	"sync"
	"time"

	"github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/application/repository/postgresql"
	"github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/application/repository/postgresql/model"
	"github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/config"
	"github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/domain/entity"
	"github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/domain/value_object"
	"github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/logger"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	log = logger.LoggerSingleton.GetInstance()
	cfg = config.ConfigSingleton.GetInstance()
)

type HashLeakedRepositoryImpl struct {
	semaphore chan struct{}
	once      sync.Once
	db        *gorm.DB
}

func NewHashLeakedRepositoryImpl() *HashLeakedRepositoryImpl {
	return &HashLeakedRepositoryImpl{
		semaphore: make(chan struct{}, cfg.Database.MaxIdleConns),
		db:        nil,
	}
}

func (r *HashLeakedRepositoryImpl) Init() {
	r.once.Do(func() {
		r.db = postgresql.DBConnectionSingleton.GetInstance()
	})
}

func (r *HashLeakedRepositoryImpl) retrySavingTransaction(hashes []*entity.LeakedHash) error {
	time.Sleep(time.Duration(cfg.Database.RetryTimeout) * time.Second)
	return r.SaveBatch(hashes)
}

func (r *HashLeakedRepositoryImpl) SaveBatch(hashes []*entity.LeakedHash) error {

	if r.db == nil {
		r.Init()
	}

	r.semaphore <- struct{}{}
	defer func() { <-r.semaphore }()

	result := r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "hash"}},
		DoNothing: true,
	}).Create(hashes)

	if result.Error == nil {
		log.Debug(logger.StoreBatchOfPasswordsMessage)
		return nil
	}

	if strings.Contains(result.Error.Error(), "40P01") {
		log.Error(logger.DeadlockDetectedMessage)
		return r.retrySavingTransaction(hashes)
	} else if strings.Contains(result.Error.Error(), "53300") {
		log.Error(logger.TooManyConnectionsMessage)
		return r.retrySavingTransaction(hashes)
	} else {
		log.Error(result.Error)
		return result.Error
	}
}

func (r *HashLeakedRepositoryImpl) Retrieve(hash value_object.PasswordHash) (*entity.LeakedHash, error) {

	if r.db == nil {
		r.Init()
	}

	var leakedHashModel model.LeakedHash

	result := r.db.Where("hash = ?", hash).First(&leakedHashModel)

	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return model.NewLeakedHashEntityFromModel(leakedHashModel), nil
}

func (r *HashLeakedRepositoryImpl) Len() int64 {

	if r.db == nil {
		r.Init()
	}

	var count int64

	result := r.db.Model(&model.LeakedHash{}).Count(&count)

	if result.Error != nil {
		return 0
	}

	return count
}
