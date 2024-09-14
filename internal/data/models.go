package data

import "database/sql"

type Models struct {
	Users UserModel 
  MetaData MetaDataModel
}

func NewModels(db *sql.DB) *Models {
	return &Models{
		Users: UserModel{DB: db}, 
    MetaData: MetaDataModel{DB: db},
	}
}
