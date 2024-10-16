package cli

import (
	"fmt"

	"github.com/lla4u/Dude/app"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
// we can attach subcommands to this command
var rootCmd = &cobra.Command{
	Use:   "dude",
	Short: "cli to manage HDX datalogs into Influx Database",
	Long:  "cli to manage HDX datalogs into Influx Database",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		readGlobalConfig()
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := runRootJob()
		// On the most outside function we only log error
		if err != nil {
			fmt.Println(err)
		}
	},
}

// runRootJob is the actual job that is executed by the root command
func runRootJob() (err error) {
	// Print global config
	GlobalConfig.Print()

	// Create new app instance
	newApp := app.NewApplication()

	return newApp.Start()
}
