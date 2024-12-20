package postgresql

import "github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/application/repository/postgresql/model"

func Init() {
	DBConnectionSingleton.GetInstance().AutoMigrate(model.LeakedHash{})
}
