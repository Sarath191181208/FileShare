package data

import (
	"database/sql"
	"time"
)

type MetaData struct {
	ID          int64     `json:"id"`
	UserId      int64     `json:"user_id"`
	Name        string    `json:"name"`
	UploadDate  time.Time `json:"upload_date"`
	Size        int64     `json:"size"`
	ContentType string    `json:"content_type"`
	FileUrl     string    `json:"-"`
}

type MetaDataModel struct {
	DB *sql.DB
}

func (m *MetaDataModel) Insert(metaData *MetaData) error {
	stmt := `INSERT INTO metadata (user_id, name, size, content_type, file_url) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return m.DB.QueryRow(stmt, metaData.UserId, metaData.Name, metaData.Size, metaData.ContentType, metaData.FileUrl).Scan(&metaData.ID)
}

func (m *MetaDataModel) Get(id int64) (*MetaData, error) {
	stmt := `SELECT id, name, upload_date, size, content_type, file_url FROM metadata WHERE id = $1`
	meta := &MetaData{}
	err := m.DB.QueryRow(stmt, id).Scan(&meta.ID, &meta.Name, &meta.UploadDate, &meta.Size, &meta.ContentType, &meta.FileUrl)
	if err != nil {
		return nil, err
	}
	return meta, nil
}

func (m *MetaDataModel) GetByUserID(id int64) ([]*MetaData, error) {
	stmt := `SELECT id, name, upload_date, size, content_type, file_url FROM metadata WHERE user_id = $1 LIMIT 10`
	rows, err := m.DB.Query(stmt, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metadata []*MetaData
	for rows.Next() {
		meta := &MetaData{}
		err := rows.Scan(&meta.ID, &meta.Name, &meta.UploadDate, &meta.Size, &meta.ContentType, &meta.FileUrl)
		if err != nil {
			return nil, err
		}
		metadata = append(metadata, meta)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return metadata, nil
}

