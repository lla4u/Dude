package cli

import (
	"github.com/leebenson/conform"
)

// init is called before main
func init() {
	// A custom sanitizer to redact sensitive data by defining a struct tag= named "redact".
	conform.AddSanitizer("redact", func(_ string) string { return "*****" })

	// Initialize the config and panic on failure
	if err := configInit(); err != nil {
		panic(err.Error())
	}
}

// Execute is the entry point for the cli
// called from main
func Execute() error {
	return rootCmd.Execute()
}
