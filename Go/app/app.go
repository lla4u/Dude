package app

import (
	"fmt"
)

// Application is the interface for the application
type Application interface {
	Start() error
	Stats(string, string) error
	Import(string, bool, string, string) error
}

// app is the implementation of the application
type app struct {
}

// NewApplication creates a new application
func NewApplication() Application {
	return &app{}
}

func (a *app) Start() error {
	fmt.Println("Starting application")
	return nil
}
