package service_test

import (
	"strings"
	"testing"

	"matrix-api/internal/domain"
	"matrix-api/internal/service"
)

// Table-Driven Test Data
var matrixTests = []struct {
	name       string
	input      string
	expectSum  string
	expectMult string
	expectFlat string
	expectErr  error
}{
	{
		name:       "Valid Square Matrix",
		input:      "1,2,3\n4,5,6\n7,8,9",
		expectSum:  "45",     // 1+2+...+9
		expectMult: "362880", // 1*2*...*9
		expectFlat: "1,2,3,4,5,6,7,8,9",
		expectErr:  nil,
	},
	{
		name:       "Matrix with Zero (Multiply Optimization)",
		input:      "1,2\n0,4",
		expectSum:  "7",
		expectMult: "0",
		expectFlat: "1,2,0,4",
		expectErr:  nil,
	},
	{
		name:      "Invalid: Not Square",
		input:     "1,2,3\n4,5,6", // 2 rows, 3 cols
		expectErr: domain.ErrNotSquare,
	},
	{
		name:      "Invalid: Not Int",
		input:     "1,2\n3,a",
		expectErr: domain.ErrNotInt,
	},
	{
		name:      "Invalid: Empty File",
		input:     "",
		expectErr: domain.ErrEmpty,
	},
	{
		name:      "Invalid: Ragged Rows",
		input:     "1,2\n3,4,5",
		expectErr: domain.ErrNotAligned,
	},
}

func TestMatrixServices(t *testing.T) {
	// We test BOTH implementations to ensure they behave identically
	services := map[string]service.MatrixProcessor{
		"Serial":     service.NewSerialMatrixService(),
		"Concurrent": service.NewConcurrentMatrixService(),
	}

	for svcName, svc := range services {
		t.Run(svcName, func(t *testing.T) {
			for _, tc := range matrixTests {
				t.Run(tc.name, func(t *testing.T) {
					// Test Sum
					sum, err := svc.Sum(strings.NewReader(tc.input))
					assertError(t, tc.expectErr, err)
					if tc.expectErr == nil && sum.String() != tc.expectSum {
						t.Errorf("Sum expected %s, got %s", tc.expectSum, sum.String())
					}

					// Test Multiply
					mult, err := svc.Multiply(strings.NewReader(tc.input))
					assertError(t, tc.expectErr, err)
					if tc.expectErr == nil && mult.String() != tc.expectMult {
						t.Errorf("Multiply expected %s, got %s", tc.expectMult, mult.String())
					}

					// Test Flatten
					flat, err := svc.Flatten(strings.NewReader(tc.input))
					assertError(t, tc.expectErr, err)
					if tc.expectErr == nil && flat != tc.expectFlat {
						t.Errorf("Flatten expected %s, got %s", tc.expectFlat, flat)
					}
				})
			}
		})
	}
}

// Test Echo/Invert specifically (since they take [][]int, not io.Reader)
func TestMatrixOperations(t *testing.T) {
	svc := service.NewSerialMatrixService() // Logic is shared, so checking one is enough

	t.Run("Echo", func(t *testing.T) {
		input := [][]int{{1, 2}, {3, 4}}
		expected := "1,2\n3,4"
		if got := svc.Echo(input); got != expected {
			t.Errorf("expected %q, got %q", expected, got)
		}
	})

	t.Run("Invert", func(t *testing.T) {
		input := [][]int{{1, 2, 3}, {4, 5, 6}} // 2x3 matrix
		// Inverted should be 3x2:
		// 1, 4
		// 2, 5
		// 3, 6
		expected := "1,4\n2,5\n3,6"
		if got := svc.Invert(input); got != expected {
			t.Errorf("expected %q, got %q", expected, got)
		}
	})
}

// Helper to compare errors
func assertError(t *testing.T, expected, actual error) {
	if expected == nil && actual != nil {
		t.Errorf("Expected no error, got %v", actual)
	}
	if expected != nil && actual != expected {
		t.Errorf("Expected error %v, got %v", expected, actual)
	}
}
