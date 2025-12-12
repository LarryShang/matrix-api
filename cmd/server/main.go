package main

import (
	"fmt"
	"mime/multipart"
	"net/http"

	"matrix-api/internal/handler"
	"matrix-api/internal/service"
)

func main() {
	// OPTION A: Use the Original (Serial) Approach
	// svc := service.NewSerialMatrixService()

	// OPTION B: Use the New (Concurrent) Approach
	svc := service.NewConcurrentMatrixService()

	matrixHandler := handler.NewMatrixHandler(svc)

	http.HandleFunc("/echo", matrixHandler.Echo)
	http.HandleFunc("/invert", matrixHandler.Invert)

	http.HandleFunc("/flatten", matrixHandler.FileHandler(func(f multipart.File) (string, error) {
		return svc.Flatten(f)
	}))
	http.HandleFunc("/sum", matrixHandler.FileHandler(func(f multipart.File) (string, error) {
		res, err := svc.Sum(f)
		if err != nil {
			return "", err
		}
		return res.String(), nil // Convert BigInt to String
	}))
	http.HandleFunc("/multiply", matrixHandler.FileHandler(func(f multipart.File) (string, error) {
		res, err := svc.Multiply(f)
		if err != nil {
			return "", err
		}
		return res.String(), nil
	}))

	fmt.Println("Server is starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server failed: %s\n", err)
	}
}
