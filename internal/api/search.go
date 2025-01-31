package api

import (
	"fmt"
	"log/slog"
	"net/http"

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

func (h *SearchHandler) HandleSearch(w http.ResponseWriter, req *http.Request) {
	n := req.PathValue("n")
	fmt.Println(n)
}
