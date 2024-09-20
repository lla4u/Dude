package cmd

import (
	"fmt"

	common "github.com/lla4u/Dude/common"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

// dbCheckCmd represents the dbCheck command
var dbCheckCmd = &cobra.Command{Use: "dbCheck", Short: "Check Influx database is ready.", Long: `Check Influx database is ready.`, Run: func(cmd *cobra.Command, args []string) {
	fmt.Println("dbCheck called")
	if Debug {
		common.LogFlags()
	}

	// Connecting de db check Health by default
	client, err := common.ConnectToInfluxDB(InfluxURL, InfluxToken)

	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Database Ok!")
	}
	// Always close client at the end
	defer client.Close()
}}

func init() {
	rootCmd.AddCommand(dbCheckCmd)
}
