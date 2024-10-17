package cli

import (
	"fmt"

	"github.com/lla4u/Dude/app"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
// we can attach subcommands to this command
var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Display recorded flight informations",
	Long:  "Display recorded flight informations",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		readGlobalConfig()
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := runstatsJob()
		// On the most outside function we only log error
		if err != nil {
			fmt.Println(err)
		}
	},
}

// Add stats command
func init() {
	rootCmd.AddCommand(statsCmd)
}

// runRootJob is the actual job that is executed by the root command
func runstatsJob() (err error) {
	// Print global config
	if GlobalConfig.Verbose {
		GlobalConfig.Print()
	}

	// Create new app instance
	newApp := app.NewApplication(app.VersionInfo{
		Version: Version,
		Commit:  CommitHash,
	})

	err = newApp.Stats(GlobalConfig.DatalogPath, GlobalConfig.Location)

	return err
}
