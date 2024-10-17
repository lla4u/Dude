package cli

import (
	"fmt"

	"github.com/lla4u/Dude/app"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
// we can attach subcommands to this command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import datalog(s) into Influx database",
	Long:  "Import datalog(s) into Influx database",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		readGlobalConfig()
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := runimportJob()
		// On the most outside function we only log error
		if err != nil {
			fmt.Println(err)
		}
	},
}

// Add stats command
func init() {
	rootCmd.AddCommand(importCmd)
}

// runRootJob is the actual job that is executed by the root command
func runimportJob() (err error) {
	// Print global config
	if GlobalConfig.Verbose {
		GlobalConfig.Print()
	}

	// Create new app instance
	newApp := app.NewApplication(app.VersionInfo {
		Version: Version,
		Commit:  CommitHash,
		}
)

	err = newApp.Import(GlobalConfig.DatalogPath, GlobalConfig.Verbose, GlobalConfig.I_Url, GlobalConfig.I_Token)

	return err
}
