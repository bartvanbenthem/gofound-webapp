package mocks

import (
	"time"

	"github.com/bartvanbenthem/gofound-webapp/internal/models"
)

var mockBlogPost = &models.BlogPost{
	ID:      1,
	Title:   "mock post",
	Content: "mock post content...",
	Created: time.Now(),
	Author:  "Admin",
	ImgURL:  "http://image_url/img",
}

type BlogPostModel struct{}

func (m *BlogPostModel) Insert(title string, content string, author string, imgURL string) (int, error) {
	return 2, nil
}

func (m *BlogPostModel) Get(id int) (*models.BlogPost, error) {
	switch id {
	case 1:
		return mockBlogPost, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (m *BlogPostModel) Latest() ([]*models.BlogPost, error) {
	return []*models.BlogPost{mockBlogPost}, nil
}

func (m *BlogPostModel) All() ([]*models.BlogPost, error) {
	return []*models.BlogPost{mockBlogPost}, nil
}
