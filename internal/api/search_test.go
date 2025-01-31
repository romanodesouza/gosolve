package api_test

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/romanodesouza/gosolve/internal/api"
	"github.com/romanodesouza/gosolve/internal/search"
	mock_search "github.com/romanodesouza/gosolve/internal/search/mocks"
	"go.uber.org/mock/gomock"
)

func TestSearchHandler_HandleSearch(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedStatus int
		expectedBody   string
		searcherMock   func(t *testing.T) search.Searcher
	}{
		{
			name:           "it should return http bad request for non-number input values",
			input:          "oops",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"index":null,"errorMessage":"\"oops\" is not a number"}`,
			searcherMock:   func(t *testing.T) search.Searcher { return nil },
		},
		{

			name:           "it should return a null index and error message when index has not been found",
			input:          "1",
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"index":null,"errorMessage":"couldn't find valid index number for given input"}`,
			searcherMock: func(t *testing.T) search.Searcher {
				ctrl := gomock.NewController(t)
				m := mock_search.NewMockSearcher(ctrl)
				m.EXPECT().
					Search(uint64(1)).
					Return(uint64(0), search.ErrNumberNotFound)

				return m
			},
		},
		{

			name:           "it should return status ok when index has been found",
			input:          "100",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"index":0,"errorMessage":null}`,
			searcherMock: func(t *testing.T) search.Searcher {
				ctrl := gomock.NewController(t)
				m := mock_search.NewMockSearcher(ctrl)
				m.EXPECT().
					Search(uint64(100)).
					Return(uint64(0), nil)

				return m
			},
		},
	}

	logger := slog.Default()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("/search/%s", tt.input)
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			mux := http.NewServeMux()
			api.NewSearchHandler(*logger, tt.searcherMock(t)).AssignRoutes(mux)
			mux.ServeHTTP(rr, req)

			if tt.expectedStatus != rr.Code {
				t.Errorf("expected %d http status code; got %d", tt.expectedStatus, rr.Code)
			}

			if strings.TrimSpace(tt.expectedBody) != strings.TrimSpace(rr.Body.String()) {
				t.Errorf(`expected %s body; got %s`, tt.expectedBody, rr.Body.String())
			}
		})
	}
}
