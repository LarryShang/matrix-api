package handler

import (
	"encoding/csv"
	"mime/multipart"
	"net/http"
	"strconv"

	"matrix-api/internal/domain"
	"matrix-api/internal/service"
)

type MatrixHandler struct {
	Service service.MatrixProcessor
}

func NewMatrixHandler(s service.MatrixProcessor) *MatrixHandler {
	return &MatrixHandler{Service: s}
}

func (h *MatrixHandler) Echo(w http.ResponseWriter, r *http.Request) {
	matrix, err := h.parseAndValidateMatrix(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte(h.Service.Echo(matrix)))
}

func (h *MatrixHandler) Invert(w http.ResponseWriter, r *http.Request) {
	matrix, err := h.parseAndValidateMatrix(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte(h.Service.Invert(matrix)))
}

// FileHandler is a reusable wrapper for endpoints that process a file stream.
// It handles the multipart form parsing so the logicFunc only worries about the file stream.
func (h *MatrixHandler) FileHandler(logicFunc func(file multipart.File) (string, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Parse the file from the request
		file, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, domain.ErrFile.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		// 2. Execute the passed service logic (Sum, Multiply, or Flatten)
		result, err := logicFunc(file)
		if err != nil {
			// We return 400 Bad Request for domain errors like "not square" or "not int"
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// 3. Write the result
		w.Write([]byte(result))
	}
}

// parseAndValidateMatrix reads the entire file into memory (for Echo/Invert only).
func (h *MatrixHandler) parseAndValidateMatrix(r *http.Request) ([][]int, error) {
	file, _, err := r.FormFile("file")
	if err != nil {
		return nil, domain.ErrFile
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()
	if err != nil {
		return nil, domain.ErrParsingCSV
	}
	if len(records) == 0 {
		return nil, domain.ErrEmpty
	}

	// Validate Square Matrix
	if len(records) != len(records[0]) {
		return nil, domain.ErrNotSquare
	}

	matrix := make([][]int, len(records))
	for i, row := range records {
		if len(row) != len(records) {
			return nil, domain.ErrNotAligned
		}
		matrix[i] = make([]int, len(row))
		for j, valStr := range row {
			val, err := strconv.Atoi(valStr)
			if err != nil {
				return nil, domain.ErrNotInt
			}
			matrix[i][j] = val
		}
	}
	return matrix, nil
}
