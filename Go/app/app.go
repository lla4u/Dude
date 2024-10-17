package app

import (
	"fmt"
)

// Application is the interface for the application
type Application interface {
	Start() error
	Stats(string, string) error
	Import(string, bool, string, string) error
	Version() string
}

// VersionInfo contains the version and commit hash
type VersionInfo struct {
	Version string
	Commit  string
}

// app is the implementation of the application
type app struct {
	versionInfo VersionInfo
}

// NewApplication creates a new application
func NewApplication(version VersionInfo) Application {
	return &app{
		versionInfo: version,
	}
}

func (a *app) Start() error {
	fmt.Println("Starting application")
	fmt.Println("Version: ", a.Version())

	return nil
}

// Version returns the version and commit hash
func (a *app) Version() string {
	return fmt.Sprintf("%s (%s)", a.versionInfo.Version, a.versionInfo.Commit)
}
