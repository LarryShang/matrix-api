package service

import (
	"encoding/csv"
	"io"
	"math/big"
	"strconv"
	"strings"

	"matrix-api/internal/domain"
)

type ConcurrentMatrixService struct {
	// Composition: Embed SerialService to reuse Echo/Invert logic
	// since those don't need concurrency changes.
	*SerialMatrixService
}

func NewConcurrentMatrixService() *ConcurrentMatrixService {
	return &ConcurrentMatrixService{
		SerialMatrixService: &SerialMatrixService{},
	}
}

// streamMatrix (Producer)
func (s *ConcurrentMatrixService) streamMatrix(file io.Reader) (<-chan int, <-chan error) {
	out := make(chan int)
	errChan := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errChan)

		reader := csv.NewReader(file)
		reader.FieldsPerRecord = -1

		firstRow, err := reader.Read()
		if err == io.EOF {
			errChan <- domain.ErrEmpty
			return
		}
		if err != nil {
			errChan <- domain.ErrParsingCSV
			return
		}

		numCols := len(firstRow)
		numRows := 1

		for _, valStr := range firstRow {
			val, err := strconv.Atoi(valStr)
			if err != nil {
				errChan <- domain.ErrNotInt
				return
			}
			out <- val
		}

		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				errChan <- domain.ErrParsingCSV
				return
			}
			numRows++
			if len(record) != numCols {
				errChan <- domain.ErrNotAligned
				return
			}
			for _, valStr := range record {
				val, err := strconv.Atoi(valStr)
				if err != nil {
					errChan <- domain.ErrNotInt
					return
				}
				out <- val
			}
		}

		if numRows != numCols {
			errChan <- domain.ErrNotSquare
			return
		}
	}()
	return out, errChan
}

// Consumers
func (s *ConcurrentMatrixService) Flatten(file io.Reader) (string, error) {
	dataCh, errCh := s.streamMatrix(file)
	var builder strings.Builder
	isFirst := true

	for val := range dataCh {
		if !isFirst {
			builder.WriteString(",")
		}
		builder.WriteString(strconv.Itoa(val))
		isFirst = false
	}
	if err := <-errCh; err != nil {
		return "", err
	}
	return builder.String(), nil
}

func (s *ConcurrentMatrixService) Sum(file io.Reader) (*big.Int, error) {
	dataCh, errCh := s.streamMatrix(file)
	sum := big.NewInt(0)
	for val := range dataCh {
		sum.Add(sum, big.NewInt(int64(val)))
	}
	if err := <-errCh; err != nil {
		return nil, err
	}
	return sum, nil
}

func (s *ConcurrentMatrixService) Multiply(file io.Reader) (*big.Int, error) {
	dataCh, errCh := s.streamMatrix(file)
	total := big.NewInt(1)
	hasZero := false

	for val := range dataCh {
		if val == 0 {
			hasZero = true
			continue
		}
		if !hasZero {
			total.Mul(total, big.NewInt(int64(val)))
		}
	}
	if err := <-errCh; err != nil {
		return nil, err
	}
	if hasZero {
		return big.NewInt(0), nil
	}
	return total, nil
}
