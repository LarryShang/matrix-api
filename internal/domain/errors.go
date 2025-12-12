package domain

import "errors"

var (
	ErrFile       = errors.New("file upload error")
	ErrParsingCSV = errors.New("csv parsing error")
	ErrEmpty      = errors.New("empty file error")
	ErrNotSquare  = errors.New("invalid matrix: number of rows must be same as columns")
	ErrNotAligned = errors.New("invalid matrix: length of columns are not aligned")
	ErrNotInt     = errors.New("invalid value: contains a non-integer")
)
