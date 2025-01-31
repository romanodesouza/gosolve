package search

import (
	"errors"
	"io"
	"log/slog"
)

var ErrNumberNotFound = errors.New("could not find given number")

type Searcher interface {
	Search(n uint64) (uint64, error)
}

type Index []uint64

var _ Searcher = Index{}

func NewIndex(r io.Reader, logger slog.Logger) (*Index, error) {
	return &Index{}, nil
}

func (i Index) Search(n uint64) (uint64, error) {
	return 0, nil
}
