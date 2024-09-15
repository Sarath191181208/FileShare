package data

import (
	"database/sql"
	"fmt"
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

func (m *MetaDataModel) Update(metaData *MetaData) error {
  stmt := `UPDATE metadata SET name = $1 WHERE id = $2`
  _, err := m.DB.Exec(stmt, metaData.Name, metaData.ID)
  return err
}

func (m *MetaDataModel) Delete(id int64) error {
  stmt := `DELETE FROM metadata WHERE id = $1`
  _, err := m.DB.Exec(stmt, id)
  return err
}


func (m *MetaDataModel) FetchTop() (*MetaData, error) {
  stmt := `SELECT id, file_url FROM metadata ORDER BY upload_date ASC LIMIT 1`
  meta := &MetaData{}
  err := m.DB.QueryRow(stmt).Scan(&meta.ID, &meta.FileUrl)
  if err != nil {
    return nil, err
  }
  return meta, nil
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
	stmt := `SELECT id, user_id, name, upload_date, size, content_type, file_url FROM metadata WHERE user_id = $1 LIMIT 10`
	rows, err := m.DB.Query(stmt, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metadata []*MetaData
	for rows.Next() {
		meta := &MetaData{}
		err := rows.Scan(&meta.ID, &meta.UserId, &meta.Name, &meta.UploadDate, &meta.Size, &meta.ContentType, &meta.FileUrl)
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

func (m *MetaDataModel) Search(user_id int64, filename string, content_type string, time *time.Time) ([]*MetaData, error) {
	// Start building the SQL query
	stmt := `SELECT id, name, upload_date, size, content_type, file_url FROM metadata WHERE user_id = $1 `
  params := []string { "$2", "$3", "$4" }
  paramIndex := 0
	args := []interface{}{user_id}

	// Add conditions based on non-empty parameters
  if filename != "" {
    stmt += fmt.Sprintf(" AND name = %s ", params[paramIndex])
    args = append(args, filename)
    paramIndex++
  }

  if content_type != "" {
    stmt += fmt.Sprintf(" AND content_type = %s ", params[paramIndex])
    args = append(args, content_type)
    paramIndex++
  }

  if time != nil {
    stmt += fmt.Sprintf(" AND upload_date > %s ", params[paramIndex])
    args = append(args, time)
    paramIndex++
  }


	stmt += ` LIMIT 10`

	// Prepare the query
	rows, err := m.DB.Query(stmt, args...)
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
