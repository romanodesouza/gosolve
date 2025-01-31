package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/romanodesouza/gosolve/internal/search"
)

type SearchHandler struct {
	logger   slog.Logger
	searcher search.Searcher
}

func NewSearchHandler(logger slog.Logger, searcher search.Searcher) *SearchHandler {
	return &SearchHandler{
		logger:   logger,
		searcher: searcher,
	}
}

func (h SearchHandler) AssignRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /search/{n}", h.HandleSearch)
}

func (h SearchHandler) HandleSearch(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	param := req.PathValue("n")
	n, err := strconv.Atoi(param)
	// Input validation response
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(SearchResponse{
			Index:        nil,
			ErrorMessage: &[]string{fmt.Sprintf(`"%s" is not a number`, param)}[0],
		})
		if err != nil {
			h.logger.Error("could not send json response", slog.Any("err", err))
		}
		return
	}

	// Perform searching
	index, err := h.searcher.Search(uint64(n))

	// Not found response
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(SearchResponse{
			Index:        nil,
			ErrorMessage: &[]string{"couldn't find valid index number for given input"}[0],
		})
		if err != nil {
			h.logger.Error("could not send json response", slog.Any("err", err))
		}
		return
	}

	// Happy response
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(SearchResponse{
		Index:        &[]uint64{index}[0],
		ErrorMessage: nil,
	})
	if err != nil {
		h.logger.Error("could not send json response", slog.Any("err", err))
	}
}

type SearchResponse struct {
	Index        *uint64 `json:"index"`
	ErrorMessage *string `json:"errorMessage"`
}
