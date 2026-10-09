package transportHTTP

import "github.com/N1ktarchik/habbit-graph/internal/core/domain"

type tarckDTO struct {
	Date  string `json:"date"`
	Value int    `json:"value"`
}

func (t *tarckDTO) UserDTOToDomain() domain.Track {
	return domain.Track{
		Date:  t.Date,
		Value: t.Value,
	}
}

func (t *tarckDTO) UserDTOFromDomain(dt domain.Track) tarckDTO {
	return tarckDTO{
		Date:  dt.Date,
		Value: dt.Value,
	}
}

type habitDTO struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

type habitResponseDTO struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Created_at  string `json:"created_at"`
}

func (h *habitDTO) HabitDTOToDomain() domain.Habit {

	return domain.Habit{
		Title:       pointerToStr(h.Title),
		Description: pointerToStr(h.Description),
	}
}

func HabitDTOFromDomain(dh domain.Habit) habitResponseDTO {
	return habitResponseDTO{
		Id:          dh.Id,
		Title:       dh.Title,
		Description: dh.Description,
		Created_at:  dh.Created_at,
	}
}

func pointerToStr(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}
