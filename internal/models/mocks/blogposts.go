package mocks

import (
	"time"

	"github.com/bartvanbenthem/gofound-webapp/internal/models"
)

var mockBlogPost = &models.BlogPost{
	ID:      1,
	Title:   "An old silent pond",
	Content: "An old silent pond...",
	Created: time.Now(),
	Expires: time.Now(),
}

type BlogPostModel struct{}

func (m *BlogPostModel) Insert(title string, content string, expires int) (int, error) {
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
