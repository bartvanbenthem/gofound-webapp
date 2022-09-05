package models

import (
	"database/sql"
	"errors"
	"time"
)

type BlogPostModelInterface interface {
	Insert(title string, content string, author string, imgURL string) (int, error)
	Get(id int) (*BlogPost, error)
	Latest() ([]*BlogPost, error)
	All() ([]*BlogPost, error)
}

type BlogPost struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Author  string
	ImgURL  string
	Version int
}

type BlogPostModel struct {
	DB *sql.DB
}

func (m *BlogPostModel) Insert(title string, content string, author string, imgURL string) (int, error) {
	stmt := `INSERT INTO blogposts (title, content, created, author, img_url)
    VALUES(?, ?, UTC_TIMESTAMP(), ?, ?)`

	result, err := m.DB.Exec(stmt, title, content, author, imgURL)
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
	stmt := `SELECT id, title, content, created, author, img_url FROM blogposts
    WHERE id = ?`

	row := m.DB.QueryRow(stmt, id)

	s := &BlogPost{}

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Author, &s.ImgURL)
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
	stmt := `SELECT id, title, content, created, author, img_url FROM blogposts
    ORDER BY id DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	blogposts := []*BlogPost{}

	for rows.Next() {
		s := &BlogPost{}

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Author, &s.ImgURL)
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

func (m *BlogPostModel) All() ([]*BlogPost, error) {
	stmt := `SELECT id, title, content, created, author, img_url FROM blogposts
    ORDER BY id DESC`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	blogposts := []*BlogPost{}

	for rows.Next() {
		s := &BlogPost{}

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Author, &s.ImgURL)
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
