package search

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"slices"
	"strconv"
)

var ErrNumberNotFound = errors.New("could not find given number")

type Searcher interface {
	Search(n uint64) (int, error)
}

type Index []uint64

var _ Searcher = Index{}

func NewIndex(r io.Reader, logger slog.Logger) (*Index, error) {
	var idx []uint64
	var lineCounter int
	buf := bufio.NewReader(r)

	for {
		lineCounter++
		line, _, err := buf.ReadLine()
		if err == io.EOF {
			logger.Debug("EOF for input file")
			break
		} else if err != nil {
			return nil, fmt.Errorf("could not read line: %w", err)
		}
		n, err := strconv.Atoi(string(line))
		if err != nil {
			logger.Warn("invalid input value", slog.Int("line_number", lineCounter))
			continue
		}

		idx = append(idx, uint64(n))
	}

	index := Index(idx)
	return &index, nil
}

func (i Index) Search(n uint64) (int, error) {
	idx, found := slices.BinarySearch(i, n)
	if found {
		return idx, nil
	}
	return -1, ErrNumberNotFound
}
