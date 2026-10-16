package transportHTTP

import (
	"log/slog"
	"net/http"

	"github.com/N1ktarchik/habbit-graph/internal/core/request"
	"github.com/N1ktarchik/habbit-graph/internal/core/response"
)

func (t *transportHTTP) GetGraph(w http.ResponseWriter, r *http.Request) {
	t.log.Debug("new request GET /api/habits/{id}/graph?year=... ")

	id, err := request.GetAndValidateID(r)
	if err != nil {
		t.log.Debug("invalid uuid", slog.Any("err", err))
		t.sendError(w, err)
	}

	graph, err := t.service.GetGraph(r.Context(), id, request.GetQueryParam(r, "year"))
	if err != nil {
		t.sendError(w, err)
		return
	}

	if err := response.ResponseOK(w, http.StatusOK, graph); err != nil {
		t.log.Warn("error in response function", slog.Any("err", err))
	}
}

func (t *transportHTTP) TrackActivity(w http.ResponseWriter, r *http.Request) {
	t.log.Debug("new request POST /api/habits/{id}/track")

	dto := &tarckDTO{}

	if err := request.DecodeJSON(r, dto); err != nil {
		t.log.Warn("error in decode function", slog.Any("err", err))
		t.sendError(w, err)
		return
	}

	id, err := request.GetAndValidateID(r)
	if err != nil {
		t.log.Debug("invalid uuid", slog.Any("err", err))
		t.sendError(w, err)
	}

	if err := t.service.SaveData(r.Context(), id, dto.UserDTOToDomain()); err != nil {
		t.sendError(w, err)
		return
	}

	if err := response.ResponseOK(w, http.StatusOK, map[string]string{"status": "ok"}); err != nil {
		t.log.Warn("error in response function", slog.Any("err", err))
	}

}
