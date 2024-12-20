package model

import (
	"github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/domain/entity"
	"github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/domain/value_object"
)

type LeakedHash struct {
	Hash string `gorm:"primaryKey"`
}

func (LeakedHash) TableName() string {
	return "leaked_hashes"
}

func NewLeakedHashEntityFromModel(leakedHashModel LeakedHash) *entity.LeakedHash {
	return &entity.LeakedHash{Hash: value_object.PasswordHash(leakedHashModel.Hash)}
}
