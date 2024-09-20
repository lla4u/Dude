package common

import (
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"github.com/gocarina/gocsv"

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

// readLines reads a whole file into memory and returns a slice of its lines.
func ReadImported(datalog string) ([]string, error) {
	const fileName = "/imported.txt"
	imported := datalog + fileName

	// Create empty file if not found
	file, err := os.OpenFile(imported, os.O_CREATE|os.O_RDONLY, 0644)
	// file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// DatalogHistory writes the imported datalog to the history (imported.txt).
func DatalogHistory(datalogfile string) {
	const fileName = "/imported.txt"
	dir := filepath.Dir(datalogfile)
	file := dir + fileName

	importedfile, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatal("Could not open file", err)
	}

	defer importedfile.Close()

	_, err2 := importedfile.WriteString(datalogfile + "\n")

	if err2 != nil {
		log.Fatal("Could not write datalog to imported.txt", err)
	} else {
		fmt.Println("Datalog appended to imported.txt")
	}
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
func Import(file string, verbose bool, url string, token string) {

	// For import timing
	now := time.Now()

	readChannel := make(chan Datalog, 1)

	// Open the CSV readFile
	readFile, err := os.OpenFile(file, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	defer readFile.Close()

	csvCount := 0
	readFromCSV(readFile, readChannel)

	var gpsDateTime string
	influxCount := 0

	// Create temp file that will hold Line protocol data
	f, err := os.CreateTemp("", "InfluxLineProtocol.*.txt")
	if err != nil {
		panic(err)
	}
	defer os.Remove(f.Name())

	// c should be re-used for further calls
	c := httpClient()

	// Consume channel
	for r := range readChannel {
		//
		// Influxdb import:
		//
		// A valid datalog require:
		//   1: valid gps Fix,
		//   2: Number of satellites up to 6
		//   3: Ground speed up to 10 Kts
		//
		if (StringToInt(r.GpsFix, verbose) >= 1) && (StringToInt(r.NumSatellites, verbose) >= 6) && (StringToFloat(r.GroundSpeed_Knots, verbose) >= 10) {

			// Save only the first record
			if r.GpsDateTime != gpsDateTime {

				// Print filtered record if verbose on
				if verbose {
					// fmt.Printf("%+v\n", r)
				}
				// Save record data into the temp file
				fmt.Fprintf(f, "datalog lat=%f,lon=%f,alt=%s,GS=%s,IAS=%s,TAS=%s,VSpeed=%d,Volts=%s,Amps=%.2f,CHTR=%.2f,CHTL=%.2f,EGT1=%d,EGT2=%d,Pitch=%.2f,Roll=%.2f,Mag=%.2f,VertAccel=%.2f,LatAccel=%.2f,OAT=%d,OilTemp=%d,OilPress=%d,RPM=%d,MAP=%.2f,FuelPress=%.2f,FuelFlow=%.2f,FuelRemaining=%.2f %d\n",
					StringToFloat(r.Lat, verbose),
					StringToFloat(r.Lon, verbose),
					r.Alt,
					r.GroundSpeed_Knots,
					r.IndicatedAirspeed_Knots,
					r.TrueAirspeed_Knots,
					StringToInt(r.VerticalSpeed_ft_min, verbose),
					r.Volts,
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
				influxCount++
			}

		}
		// Update the csvCount
		csvCount++

	}

	// Flush Line protocol temp file
	response := sendRequest(c, f, url, token)

	// Print sendReques responce if not empty
	if len(response) >= 1 {
		log.Println("Response Body:", string(response))
	}

	// Append record information
	fmt.Println(" -", time.Since(now), csvCount, influxCount)

	DatalogHistory(file)
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

// Client for ILP
func httpClient() *http.Client {
	client := &http.Client{Timeout: 10 * time.Second}
	return client
}

// Influx send request using ILP
func sendRequest(client *http.Client, f *os.File, url string, token string) []byte {
	endpoint := url + "/api/v2/write?org=dude&bucket=dude&precision=s"
	Token := "Token " + token

	req, err := http.NewRequest("POST", endpoint, f)
	if err != nil {
		log.Fatalf("Error Occurred. %+v", err)
	}

	// 	req.Header.Set("Authorization", "Token my-super-secret-auth-token")
	req.Header.Set("Authorization", Token)
	// req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	req.Header.Set("Accept", "application/json")

	// Do the api call
	response, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request to API endpoint. %+v", err)
	}

	// Close the connection to reuse it
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Couldn't parse response body. %+v", err)
	}

	return body
}
