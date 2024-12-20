package entity

import "github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/domain/value_object"

type LeakedHash struct {
	Hash value_object.PasswordHash
}
