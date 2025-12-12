package handler_test

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"matrix-api/internal/handler"
	"matrix-api/internal/service"
)

// Helper to create a multipart request
func createUploadRequest(t *testing.T, url, fieldName, content string) *http.Request {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Create the form file part
	part, err := writer.CreateFormFile(fieldName, "matrix.csv")
	if err != nil {
		t.Fatal(err)
	}
	// Copy content into the part
	_, err = io.Copy(part, strings.NewReader(content))
	if err != nil {
		t.Fatal(err)
	}
	writer.Close() // Important: Close to write the boundary

	req := httptest.NewRequest("POST", url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

func TestMatrixHandler_Endpoints(t *testing.T) {
	// Setup
	svc := service.NewConcurrentMatrixService()
	h := handler.NewMatrixHandler(svc)

	// Define the wrapper functions exactly as they appear in main.go
	flattenHandler := h.FileHandler(func(f multipart.File) (string, error) {
		return svc.Flatten(f)
	})
	sumHandler := h.FileHandler(func(f multipart.File) (string, error) {
		res, err := svc.Sum(f)
		if err != nil {
			return "", err
		}
		return res.String(), nil
	})

	tests := []struct {
		name           string
		endpoint       string
		handlerFunc    http.HandlerFunc
		fileContent    string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Flatten Success",
			endpoint:       "/flatten",
			handlerFunc:    flattenHandler,
			fileContent:    "1,2\n3,4",
			expectedStatus: http.StatusOK,
			expectedBody:   "1,2,3,4",
		},
		{
			name:           "Sum Success",
			endpoint:       "/sum",
			handlerFunc:    sumHandler,
			fileContent:    "1,2\n3,4",
			expectedStatus: http.StatusOK,
			expectedBody:   "10",
		},
		{
			name:           "Sum Invalid CSV",
			endpoint:       "/sum",
			handlerFunc:    sumHandler,
			fileContent:    "1,2\n3,a", // Invalid int
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid value: contains a non-integer",
		},
		{
			name:           "Echo Success",
			endpoint:       "/echo",
			handlerFunc:    h.Echo,
			fileContent:    "1,2\n3,4",
			expectedStatus: http.StatusOK,
			expectedBody:   "1,2\n3,4",
		},
		{
			name:           "Invert Success",
			endpoint:       "/invert",
			handlerFunc:    h.Invert,
			fileContent:    "1,2\n4,5",
			expectedStatus: http.StatusOK,
			expectedBody:   "1,4\n2,5",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := createUploadRequest(t, tc.endpoint, "file", tc.fileContent)
			rr := httptest.NewRecorder()

			tc.handlerFunc.ServeHTTP(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					rr.Code, tc.expectedStatus)
			}

			// Trim spaces/newlines for easier comparison
			gotBody := strings.TrimSpace(rr.Body.String())
			if tc.expectedBody != "" && gotBody != tc.expectedBody {
				t.Errorf("handler returned unexpected body: got %v want %v",
					gotBody, tc.expectedBody)
			}
		})
	}
}

func TestMatrixHandler_InvalidUpload(t *testing.T) {
	svc := service.NewSerialMatrixService()
	h := handler.NewMatrixHandler(svc)

	// Test case: Sending a request WITHOUT a file
	req := httptest.NewRequest("POST", "/echo", nil) // No body
	rr := httptest.NewRecorder()

	h.Echo(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing file, got %d", rr.Code)
	}
}
