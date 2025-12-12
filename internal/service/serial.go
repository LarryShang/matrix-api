package service

import (
	"encoding/csv"
	"io"
	"math/big"
	"strconv"
	"strings"

	"matrix-api/internal/domain"
)

type SerialMatrixService struct{}

func NewSerialMatrixService() *SerialMatrixService {
	return &SerialMatrixService{}
}

// Echo and Invert are shared logic, but implemented here for the struct
func (s *SerialMatrixService) Echo(matrix [][]int) string {
	var builder strings.Builder
	for i, row := range matrix {
		for j, val := range row {
			builder.WriteString(strconv.Itoa(val))
			if j < len(row)-1 {
				builder.WriteString(",")
			}
		}
		if i < len(matrix)-1 {
			builder.WriteString("\n")
		}
	}
	return builder.String()
}

func (s *SerialMatrixService) Invert(matrix [][]int) string {
	if len(matrix) == 0 {
		return ""
	}
	rows, cols := len(matrix), len(matrix[0])
	inverted := make([][]int, cols)
	for i := range inverted {
		inverted[i] = make([]int, rows)
	}
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			inverted[j][i] = matrix[i][j]
		}
	}
	return s.Echo(inverted)
}

// processStream is the helper for the Serial approach
func (s *SerialMatrixService) processStream(file io.Reader, process func(val int) error) error {
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	firstRow, err := reader.Read()
	if err == io.EOF {
		return domain.ErrEmpty
	}
	if err != nil {
		return domain.ErrParsingCSV
	}

	numCols := len(firstRow)
	numRows := 1

	for _, valStr := range firstRow {
		val, err := strconv.Atoi(valStr)
		if err != nil {
			return domain.ErrNotInt
		}
		if err := process(val); err != nil {
			return err
		}
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return domain.ErrParsingCSV
		}
		numRows++
		if len(record) != numCols {
			return domain.ErrNotAligned
		}
		for _, valStr := range record {
			val, err := strconv.Atoi(valStr)
			if err != nil {
				return domain.ErrNotInt
			}
			if err := process(val); err != nil {
				return err
			}
		}
	}

	if numRows != numCols {
		return domain.ErrNotSquare
	}
	return nil
}

func (s *SerialMatrixService) Flatten(file io.Reader) (string, error) {
	var builder strings.Builder
	isFirst := true
	processor := func(val int) error {
		if !isFirst {
			builder.WriteString(",")
		}
		builder.WriteString(strconv.Itoa(val))
		isFirst = false
		return nil
	}
	err := s.processStream(file, processor)
	return builder.String(), err
}

func (s *SerialMatrixService) Sum(file io.Reader) (*big.Int, error) {
	sum := big.NewInt(0)
	currentVal := new(big.Int)
	processor := func(val int) error {
		sum.Add(sum, currentVal.SetInt64(int64(val)))
		return nil
	}
	err := s.processStream(file, processor)
	return sum, err
}

func (s *SerialMatrixService) Multiply(file io.Reader) (*big.Int, error) {
	total := big.NewInt(1)
	currentVal := new(big.Int)
	isZero := false
	processor := func(val int) error {
		if val == 0 {
			isZero = true
		}
		if !isZero {
			total.Mul(total, currentVal.SetInt64(int64(val)))
		}
		return nil
	}
	err := s.processStream(file, processor)
	if isZero {
		return big.NewInt(0), nil
	}
	return total, err
}
