package search_test

import (
	"errors"
	"testing"

	"github.com/romanodesouza/gosolve/internal/search"
)

func TestIndex_Search(t *testing.T) {
	tests := []struct {
		name          string
		index         search.Index
		search        uint64
		expectedIndex int
		expectedError error
	}{
		{
			name:          "it should return correct index of valid value",
			index:         search.Index([]uint64{0, 100, 200}),
			search:        uint64(100),
			expectedIndex: 1,
			expectedError: nil,
		},
		{
			name:          "it should return an error when index has not been found",
			index:         search.Index([]uint64{0, 100, 200}),
			search:        uint64(300),
			expectedIndex: -1,
			expectedError: search.ErrNumberNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n, err := tt.index.Search(tt.search)

			if tt.expectedIndex != n {
				t.Errorf("expected index: %d; got %d", tt.expectedIndex, n)
			}

			if tt.expectedError == nil && err != nil {
				t.Errorf("unexpected error: %#v", err)
			}

			if !errors.Is(err, tt.expectedError) {
				t.Errorf(`expected error: "%s"; got "%s"`, tt.expectedError.Error(), err.Error())
			}
		})
	}
}
