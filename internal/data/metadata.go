package data

import (
	"database/sql"
	"time"
)

type MetaData struct{
  ID int64 `json:"id"`
  Name string `json:"name"`
  UploadDate time.Time `json:"upload_date"`
  Size int64 `json:"size"`
  ContentType string `json:"content_type"`
}

type MetaDataModel struct{
  DB *sql.DB
}

func (m *MetaDataModel) Insert(metaData *MetaData) error {
  stmt := `INSERT INTO metadata (name, size, content_type) VALUES ($1, $2, $3, $4) RETURNING id`
  return m.DB.QueryRow(stmt, metaData.Name, metaData.Size, metaData.ContentType).Scan(&metaData.ID)
}
