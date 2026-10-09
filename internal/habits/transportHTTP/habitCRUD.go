package transportHTTP

import (
	"log/slog"
	"net/http"

	"github.com/N1ktarchik/habbit-graph/internal/core/domain"
	"github.com/N1ktarchik/habbit-graph/internal/core/request"
	"github.com/N1ktarchik/habbit-graph/internal/core/response"
)

func (t *transportHTTP) CreateHabit(w http.ResponseWriter, r *http.Request) {
	t.log.Debug("new request POST /api/habits ")

	habit := &habitDTO{}
	if err := request.DecodeJSON(r, habit); err != nil {
		t.sendError(w, err)
		return
	}

	savedHabit, err := t.service.CreateHabit(r.Context(), habit.HabitDTOToDomain())
	if err != nil {
		t.sendError(w, err)
		return
	}

	if err := response.ResponseOK(w, http.StatusCreated, HabitDTOFromDomain(*savedHabit)); err != nil {
		t.log.Warn("error in response function", slog.Any("err", err))
	}

}

func (t *transportHTTP) GetHabits(w http.ResponseWriter, r *http.Request) {
	t.log.Debug("new request GET /api/habits ")

	habits, err := t.service.GetHabits(r.Context())
	if err != nil {
		t.sendError(w, err)
		return
	}

	responseHabits := make([]habitResponseDTO, 0, len(habits))
	for _, h := range habits {
		responseHabits = append(responseHabits, HabitDTOFromDomain(h))
	}

	if err := response.ResponseOK(w, http.StatusOK, responseHabits); err != nil {
		t.log.Warn("error in response function", slog.Any("err", err))
	}
}

func (t *transportHTTP) DeleteHabit(w http.ResponseWriter, r *http.Request) {
	t.log.Debug("new request DELETE /api/habits ")

	if err := t.service.DeleteHabit(r.Context(), request.GetID(r)); err != nil {
		t.sendError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (t *transportHTTP) UpdateHabit(w http.ResponseWriter, r *http.Request) {
	t.log.Debug("new request PATCH /api/habits ")

	habit := &habitDTO{}
	if err := request.DecodeJSON(r, habit); err != nil {
		t.sendError(w, err)
		return
	}

	updateHabitInput := &domain.UpdateHabitInput{
		Title:       habit.Title,
		Description: habit.Description,
	}

	savedHabit, err := t.service.UpdateHabit(r.Context(), request.GetID(r), updateHabitInput)
	if err != nil {
		t.sendError(w, err)
		return
	}

	if err := response.ResponseOK(w, http.StatusOK, HabitDTOFromDomain(*savedHabit)); err != nil {
		t.log.Warn("error in response function", slog.Any("err", err))
	}

}
