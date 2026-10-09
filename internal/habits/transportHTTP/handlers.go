package transportHTTP

import (
	"context"
	"net/http"

	"github.com/N1ktarchik/habbit-graph/internal/core/request"
	"github.com/N1ktarchik/habbit-graph/internal/core/response"
)

func (t *transportHTTP) GetGraph(w http.ResponseWriter, r *http.Request) {
	t.log.Debug("new request GET /api/graph ")

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	graph, err := t.serv.GetGraph(ctx)
	if err != nil {

		if respErr := response.ResponseWithError(w, err); respErr != nil {
			t.log.Warn("error in response function", err)
		}

		return
	}

	if err := response.ResponseWithMap(w, graph); err != nil {
		t.log.Warn("error in response function", err)
	}
}

func (t *transportHTTP) PostUserData(w http.ResponseWriter, r *http.Request) {
	t.log.Debug("new request POST /api/track")

	userDto := &userDTO{}

	if err := request.DecodeAndValidateJSON(r, userDto); err != nil {
		t.log.Warn("error in decode function", err)

		if respErr := response.ResponseWithError(w, err); respErr != nil {
			t.log.Warn("error in response function", err)
		}

		return
	}

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	if err := t.serv.SaveData(ctx, userDto.dtoToDomain()); err != nil {

		if respErr := response.ResponseWithError(w, err); respErr != nil {
			t.log.Warn("error in response function", err)
		}

		return
	}

	if err := response.ResponseWithMap(w, map[string]string{"status": "ok"}); err != nil {
		t.log.Warn("error in response function", err)
	}

}
