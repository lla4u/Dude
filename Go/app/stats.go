package app

import (
	"fmt"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

// Imported yml structure helper
type Flight struct {
	Start    time.Time     `yaml:"start"`
	End      time.Time     `yaml:"stop"`
	Duration time.Duration `yaml:"duration"`
}

// Imported yml structure
type Imported struct {
	Datalogs []string `yaml:"datalogs"`
	Flights  []Flight `yaml:"flights"`
}

// Stats
func (a *app) Stats(d string, l string) error {
	// Read imported file
	imported := ReadImported(d)

	fmt.Printf("Datalogs count: %d \n", len(imported.Datalogs))
	fmt.Printf("Flights count: %d - Total Duration: %s\n", len(imported.Flights), totalFlightDuration(imported))
	fmt.Println("")

	DisplayFlights(imported, l)

	return nil
}

// Stats helper
func ReadImported(datalog string) Imported {
	const fileName = "/imported.yml"
	filePath := datalog + fileName

	var data Imported

	// Check file exist and create if not (needed at first run)
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		file, err := os.Create(filePath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
	}

	yamlFile, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(yamlFile, &data)
	if err != nil {
		log.Fatal(err)
	}

	return data
}

// Stats helper
func totalFlightDuration(i Imported) time.Duration {
	var sum time.Duration

	for _, v := range i.Flights {
		sum += v.Duration
	}
	return sum
}

// Stats helper
func DisplayFlights(i Imported, location string) {
	// Flights Start & End are in UTC and the timezone will be set to whatever your local timezone is
	loc, _ := time.LoadLocation(location)

	for j, f := range i.Flights {
		fmt.Printf("%d - %s -> %s : %s\n", j+1, f.Start.In(loc).String(), f.End.In(loc).String(), f.Duration)
		fmt.Printf("%d - http://localhost:3000/d/TxT6-pXSz/dude-dashboard?orgId=1&from=%d&to=%d&kiosk\n\n", j+1, f.Start.In(loc).UnixMilli(), f.End.In(loc).UnixMilli())
	}
}
