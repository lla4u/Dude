package common

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"time"

	"github.com/gocarina/gocsv"
	"gopkg.in/yaml.v3"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Datalog struct {
	GpsFix                  string `csv:"GPS Fix Quality"`
	NumSatellites           string `csv:"Number of Satellites"`
	GpsDateTime             string `csv:"GPS Date & Time"`
	Lat                     string `csv:"Latitude (deg)"`
	Lon                     string `csv:"Longitude (deg)"`
	Alt                     string `csv:"GPS Altitude (feet)"`
	GroundSpeed_Knots       string `csv:"Ground Speed (knots)"`
	Pitch_Deg               string `csv:"Pitch (deg)"`
	Roll_Deg                string `csv:"Roll (deg)"`
	MagneticHeading_Deg     string `csv:"Magnetic Heading (deg)"`
	IndicatedAirspeed_Knots string `csv:"Indicated Airspeed (knots)"`
	LateralAccel_G          string `csv:"Lateral Accel (g)"`
	VerticalAccel_G         string `csv:"Vertical Accel(g)"`
	VerticalSpeed_ft_min    string `csv:"Vertical Speed (ft/min)"`
	OAT_Deg_C               string `csv:"OAT (deg C)"`
	TrueAirspeed_Knots      string `csv:"True Airspeed (knots)"`
	WindDirection_Deg       string `csv:"Wind Direction (deg)"`
	WindSpeed_Knots         string `csv:"Wind Speed (knots)"`
	Oil_Pressure_PSI        string `csv:"Oil Pressure (PSI)"`
	OilTemp_Deg_C           string `csv:"Oil Temp (deg C)"`
	RPM                     string `csv:"RPM L"`
	ManifoldPressure_inHg   string `csv:"Manifold Pressure (inHg)"`
	FuelFlow1_Gal_hr        string `csv:"Fuel Flow 1 (gal/hr)"`
	FuelPressure_PSI        string `csv:"Fuel Pressure (PSI)"`
	FuelRemaining_Gal       string `csv:"Fuel Remaining (gal)"`
	Volts                   string `csv:"Volts 1"`
	Amps                    string `csv:"Amps"`
	EGT1_Deg_C              string `csv:"EGT 1 (deg C)"`
	EGT2_Deg_C              string `csv:"EGT 2 (deg C)"`
	CHTL_Deg_C              string `csv:"CHTL TEMPERATURE (deg C)"`
	CHTR_Deg_C              string `csv:"CHTR TEMPERATURE (deg C)"`
}

type Flight struct {
	Start    time.Time     `yaml:"start"`
	End      time.Time     `yaml:"stop"`
	Duration time.Duration `yaml:"duration"`
}

type Imported struct {
	Datalogs []string `yaml:"datalogs"`
	Flights  []Flight `yaml:"flights"`
}

// Read imported datalogs & flights
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

// Save imported datalogs & flights
func SaveImported(data Imported) {
	const fileName = "/imported.yml"
	file := filepath.Dir(data.Datalogs[0]) + fileName

	// Marshal the structure
	newYamlFile, err := yaml.Marshal(&data)
	if err != nil {
		panic(err)
	}

	f, err := os.Create(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = io.Writer.Write(f, newYamlFile)
	if err != nil {
		panic(err)
	}
}

func LogFlags() {
	for key, value := range viper.GetViper().AllSettings() {
		log.WithFields(log.Fields{
			key: value,
		}).Info("Command Flag")
	}
}

func ConnectToInfluxDB(url string, token string) (influxdb2.Client, error) {

	// client := influxdb2.NewClient(dbURL, dbToken)
	client := influxdb2.NewClientWithOptions(url, token,
		influxdb2.DefaultOptions().SetBatchSize(2000))

	// validate client connection health
	_, err := client.Health(context.Background())

	return client, err
}

// Usage: files, err := WalkMatch("/root/", "*.md")
func WalkMatch(root, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}

// Strictly compare one slice against the other
func Diff(a []string, b []string) []string {
	// Turn b into a map
	m := make(map[string]bool, len(b))
	for _, s := range b {
		m[s] = false
	}
	// Append values from the longest slice that don't exist in the map
	var diff []string
	for _, s := range a {
		if _, ok := m[s]; !ok {
			diff = append(diff, s)
			continue
		}
		m[s] = true
	}
	// Sort the resulting slice
	sort.Strings(diff)
	return diff
}

// Import datalog file into Influx database
func Import(imported *Imported, file string, verbose bool, url string, token string) {

	var gpsDateTime string
	var influxCount int
	var csvCount int
	var currentTime time.Time
	var lastTime time.Time
	var startTime time.Time
	var skiped int

	// For import timing
	now := time.Now()

	readChannel := make(chan Datalog, 1)

	// Open the CSV readFile
	readFile, err := os.OpenFile(file, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	defer readFile.Close()

	readFromCSV(readFile, readChannel)

	// Create / replace file that will hold Influx Line Protocol data
	ilpFile, err := os.Create("InfluxLineProtocol.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer ilpFile.Close()

	// Consume CSV channel
	for r := range readChannel {
		//
		// Influxdb import:
		//
		// Update the csvCount
		csvCount++
		// A valid datalog require a valid gps Fix, Number of satellites up to 6, Ground speed up to 10 Kts
		if (StringToInt(r.GpsFix, verbose) >= 1) && (StringToInt(r.NumSatellites, verbose) >= 6) && (StringToFloat(r.GroundSpeed_Knots, verbose) >= 10) {

			currentTime, err = time.Parse("2006-01-02 15:04:05", r.GpsDateTime)
			if err != nil {
				log.Fatal(csvCount, err)
			}

			// Skip already recorded flight data
			if existingFlight(imported, currentTime) {
				skiped++
				continue
			}

			// One record sample per second at today
			if r.GpsDateTime != gpsDateTime {

				// A new flight if record stop duration is up to 10 minutes
				if currentTime.Sub(lastTime).Minutes() >= 10 {

					var flight Flight

					flight.Start = startTime
					flight.End = lastTime
					flight.Duration = lastTime.Sub(startTime)

					if flight.Duration != 0 {
						imported.Flights = append(imported.Flights, flight)

						// Clean
						imported.Flights = uniqueFlight(imported.Flights)

						slices.SortFunc(imported.Flights,
							func(a, b Flight) int {
								return a.Start.Compare(b.Start)
							})

					}

					startTime = currentTime

				}

				// Print filtered record if verbose on
				if verbose {
					// fmt.Printf("%+v\n", r)
				}

				// Save record data into the influxdb line protocol file
				fmt.Fprintf(ilpFile, "datalog lat=%f,lon=%f,alt=%d,GS=%.2f,IAS=%.2f,TAS=%.2f,VSpeed=%d,Volts=%.2f,Amps=%.2f,CHTR=%.2f,CHTL=%.2f,EGT1=%d,EGT2=%d,Pitch=%.2f,Roll=%.2f,Mag=%.2f,VertAccel=%.2f,LatAccel=%.2f,OAT=%d,OilTemp=%d,OilPress=%d,RPM=%d,MAP=%.2f,FuelPress=%.2f,FuelFlow=%.2f,FuelRemaining=%.2f %d\n",
					StringToFloat(r.Lat, verbose),
					StringToFloat(r.Lon, verbose),
					StringToInt(r.Alt, verbose),
					StringToFloat(r.GroundSpeed_Knots, verbose),
					StringToFloat(r.IndicatedAirspeed_Knots, verbose),
					StringToFloat(r.TrueAirspeed_Knots, verbose),
					StringToInt(r.VerticalSpeed_ft_min, verbose),
					StringToFloat(r.Volts, verbose),
					StringToFloat(r.Amps, verbose),
					StringToFloat(r.CHTR_Deg_C, verbose),
					StringToFloat(r.CHTL_Deg_C, verbose),
					StringToInt(r.EGT1_Deg_C, verbose),
					StringToInt(r.EGT2_Deg_C, verbose),
					StringToFloat(r.Pitch_Deg, verbose),
					StringToFloat(r.Roll_Deg, verbose),
					StringToFloat(r.MagneticHeading_Deg, verbose),
					StringToFloat(r.VerticalAccel_G, verbose),
					StringToFloat(r.LateralAccel_G, verbose),
					StringToInt(r.OAT_Deg_C, verbose),
					StringToInt(r.OilTemp_Deg_C, verbose),
					StringToInt(r.Oil_Pressure_PSI, verbose),
					StringToInt(r.RPM, verbose),
					StringToFloat(r.ManifoldPressure_inHg, verbose),
					StringToFloat(r.FuelPressure_PSI, verbose),
					StringToFloat(r.FuelFlow1_Gal_hr, verbose),
					StringToFloat(r.FuelRemaining_Gal, verbose),
					dateStringToUnix(r.GpsDateTime).Unix())

				// Update the gpsDateTime and influx record count
				gpsDateTime = r.GpsDateTime
				lastTime = currentTime
				influxCount++
			}

		}

	}
	ilpFile.Close()

	// Flush Line protocol temp file
	response := sendRequest(url, token)

	// Print sendRequest responce if not empty
	if len(response) >= 1 {
		log.Println("Response Body:", string(response))
	}

	// Append record information
	fmt.Println(" -", time.Since(now), csvCount, influxCount, skiped)
}

// Read CSV
func readFromCSV(file *os.File, c chan Datalog) {
	gocsv.SetCSVReader(func(r io.Reader) gocsv.CSVReader {
		reader := csv.NewReader(r)
		reader.LazyQuotes = true
		reader.FieldsPerRecord = -1
		return reader
	})

	// Read the CSV file into a slice of Record structs
	go func() {
		err := gocsv.UnmarshalToChan(file, c)
		if err != nil {
			panic(err)
		}
	}()
}

// Convert string to Float
func StringToFloat(s string, verbose bool) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		// log only in verbose mode
		if verbose {
			log.Println(s, err)
		}
	}
	return f
}

// Convert string to Int
func StringToInt(s string, verbose bool) int {
	f, err := strconv.Atoi(s)
	if err != nil {
		if verbose {
			log.Println(s, err)
		}
	}
	return f
}

// Influx expected time format
func dateStringToUnix(s string) time.Time {
	layout := "2006-01-02 15:04:05"
	date, _ := time.Parse(layout, s)
	return date
}

// Influx send request using ILP
func sendRequest(url string, token string) []byte {
	const ilpFileName = "InfluxLineProtocol.txt"

	endpoint := url + "/api/v2/write?org=dude&bucket=dude&precision=s"
	Token := "Token " + token

	// Client for ILP
	client := &http.Client{Timeout: 10 * time.Second}

	// Open ILP file
	ilpFile, err := os.Open(ilpFileName)
	if err != nil {
		log.Fatalf("Error Occurred opening ILP file. %+v", err)
	}
	defer ilpFile.Close()

	req, err := http.NewRequest("POST", endpoint, ilpFile)
	if err != nil {
		log.Fatalf("Error Occurred. %+v", err)
	}

	req.Header.Set("Authorization", Token)
	// req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	req.Header.Set("Accept", "application/json")

	// Do the api call
	response, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request to API endpoint. %+v", err)
	}

	// Close the connection
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Couldn't parse response body. %+v", err)
	}

	// Remove ilp file
	os.Remove(ilpFileName)

	return body
}

// One flight at a time
func uniqueFlight(stringSlice []Flight) []Flight {
	keys := make(map[time.Time]bool)
	list := []Flight{}
	for _, entry := range stringSlice {
		if _, value := keys[entry.Start]; !value {
			keys[entry.Start] = true
			list = append(list, entry)
		}
	}
	return list
}

// Helper function for skip
func timeIsBetween(t, min, max time.Time) bool {
	if min.After(max) {
		min, max = max, min
	}
	return (t.Equal(min) || t.After(min)) && (t.Equal(max) || t.Before(max))
}

// Skip logic on CSV record already onboarded
func existingFlight(i *Imported, currentTime time.Time) bool {
	for _, flight := range i.Flights {
		//fmt.Println("checking:", currentTime, " with:", flight.Start, flight.End)
		if timeIsBetween(currentTime, flight.Start, flight.End) {
			return true
		}
	}
	return false
}
