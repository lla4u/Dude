package cmd

import (
	"fmt"

	common "github.com/lla4u/Dude/common"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import datalog(s) into Influx database.",
	Long:  `Import datalog(s) into Influx database.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("import called")

		if Debug {
			common.LogFlags()
		}

		// Look for datalog files to import
		datalogs, err := common.WalkMatch(DatalogsDir, "*USER_LOG_DATA*.csv")
		if err != nil {
			log.Fatal(err)
		}

		if Verbose {
			fmt.Println("Found:", len(datalogs), "datalog files to import")
		}

		// Read imported datalog files
		imported := common.ReadImported(DatalogsDir)

		if Verbose {
			fmt.Println("Got:", len(imported.Datalogs), "imported datalog files")
		}

		datalogsToImport := common.Diff(datalogs, imported.Datalogs)

		// Display missing datalogs if verbose
		if Verbose {
			for i, datatalog := range datalogsToImport {
				fmt.Println("Missing:", i+1, datatalog)
			}
		}

		// Import the missing datalogs if needed
		if len(datalogsToImport) == 0 {
			fmt.Println("Nothing to do!")
		} else {

			for i, datatalog := range datalogsToImport {
				fmt.Printf("Importing: %d %s", i+1, datatalog)
				common.Import(&imported, datatalog, Verbose, InfluxURL, InfluxToken)
				// Lastly add imported datalog & save into imported yml file
				imported.Datalogs = append(imported.Datalogs, datatalog)
				common.SaveImported(imported)
			}
			fmt.Printf("Datalogs: %d | Flights: %d\n", len(imported.Datalogs), len(imported.Flights))

		}

	},
}

func init() {
	rootCmd.AddCommand(importCmd)

	importCmd.PersistentFlags().StringVarP(&DatalogsDir, "datalog", "", "/Users/lla/Documents/Laurent/Aviation/P300 Dude", "Datalogs directory path.")
	viper.BindPFlag("datalog", rootCmd.PersistentFlags().Lookup("datalog"))
}
