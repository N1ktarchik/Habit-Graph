package transportHTTP

import "github.com/N1ktarchik/habbit-graph/internal/core/domain"

type userDTO struct {
	Date  string `json:"date"`
	Value int    `json:"value"`
}

func (u *userDTO) dtoToDomain() *domain.UserData {
	return &domain.UserData{
		Date:  u.Date,
		Value: u.Value,
	}
}
