package models

import (
	"database/sql"
	"errors"
	"time"
)

type BlogPostModelInterface interface {
	Insert(title string, content string, expires int) (int, error)
	Get(id int) (*BlogPost, error)
	Latest() ([]*BlogPost, error)
}

type BlogPost struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type BlogPostModel struct {
	DB *sql.DB
}

func (m *BlogPostModel) Insert(title string, content string, expires int) (int, error) {
	stmt := `INSERT INTO blogposts (title, content, created, expires)
    VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *BlogPostModel) Get(id int) (*BlogPost, error) {
	stmt := `SELECT id, title, content, created, expires FROM blogposts
    WHERE expires > UTC_TIMESTAMP() AND id = ?`

	row := m.DB.QueryRow(stmt, id)

	s := &BlogPost{}

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

func (m *BlogPostModel) Latest() ([]*BlogPost, error) {
	stmt := `SELECT id, title, content, created, expires FROM blogposts
    ORDER BY id DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	blogposts := []*BlogPost{}

	for rows.Next() {
		s := &BlogPost{}

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		blogposts = append(blogposts, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return blogposts, nil
}
