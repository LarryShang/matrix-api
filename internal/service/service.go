package service

import (
	"io"
	"math/big"
)

// MatrixProcessor defines the contract for matrix operations.
type MatrixProcessor interface {
	Echo(matrix [][]int) string
	Invert(matrix [][]int) string
	Flatten(file io.Reader) (string, error)
	Sum(file io.Reader) (*big.Int, error)
	Multiply(file io.Reader) (*big.Int, error)
}
