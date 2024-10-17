package cli

import (
	"fmt"
	"os"

	"github.com/fatih/structs"
	"github.com/leebenson/conform"
	"github.com/sanity-io/litter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// handle global configuration through a config file, environment vars  cli parameters.

// GlobalConfig the global config object
var GlobalConfig *config

func readGlobalConfig() {
	// Priority of configuration options
	// 1: CLI Parameters
	// 2: environment
	// 3: config.yaml
	// 4: defaults
	config, err := readConfig()
	if err != nil {
		panic(err.Error())
	}

	// Set config object for main package
	GlobalConfig = config
}

var defaultConfig = &config{
	I_Url:   "http://localhost:8086",
	I_Token: "my-super-secret-auth-token",
	Verbose: false,
}

// configInit must be called from the packages' init() func
func configInit() error {
	cliFlags()
	return bindFlagsAndEnv()
}

// Create private data struct to hold config options.
// `mapstructure` => viper tags
// `struct` => fatih structs tag
// `env` => environment variable name
type config struct {
	DatalogPath string `mapstructure:"datalogpath" structs:"datalogpath" env:"DATALOGPATH"`
	I_Url       string `mapstructure:"iurl" structs:"iurl" env:"IURL"`
	I_Token     string `mapstructure:"itoken" structs:"itoken" env:"ITOKEN" conform:"redact"`
	Verbose     bool   `mapstructure:"verbose" structs:"verbose" env:"VERBOSE"`
	Location    string `mapstructure:"location" structs:"location" env:"LOCATION"`
}

// cliFlags defines cli parameters for all config options
func cliFlags() {
	// Keep cli parameters in sync with the config struct

	// Example params
	rootCmd.PersistentFlags().String("datalogpath", defaultConfig.DatalogPath, "Datalog directory path")
	rootCmd.PersistentFlags().String("iurl", defaultConfig.I_Url, "Grafana URL")
	rootCmd.PersistentFlags().String("itoken", defaultConfig.I_Token, "Grafana Token")
	rootCmd.PersistentFlags().Bool("verbose", defaultConfig.Verbose, "Enable verbose mode")
	rootCmd.PersistentFlags().String("location", defaultConfig.Location, "Location IE: Europe/Paris")
}

// bindFlagsAndEnv will assign the environment variables to the cli parameters
func bindFlagsAndEnv() (err error) {
	for _, field := range structs.Fields(&config{}) {
		// Get the struct tag values
		key := field.Tag("structs")
		env := field.Tag("env")

		// Bind cobra flags to viper
		err = viper.BindPFlag(key, rootCmd.PersistentFlags().Lookup(key))
		if err != nil {
			return err
		}
		err = viper.BindEnv(key, env)
		if err != nil {
			return err
		}
	}
	return nil
}

// Print the config object
// but remove sensitive data
func (c *config) Print() {
	cp := *c
	_ = conform.Strings(&cp)
	litter.Dump(cp)
}

// String the config object
// but remove sensitive data
func (c *config) String() string {
	cp := *c
	_ = conform.Strings(&cp)
	return litter.Sdump(cp)
}

// readConfig a helper to read default from a default config object.
func readConfig() (*config, error) {
	// Create a map of the default config
	defaultsAsMap := structs.Map(defaultConfig)

	// Set defaults
	for key, value := range defaultsAsMap {
		viper.SetDefault(key, value)
	}

	// Read config from file
	// Find user HOME directory
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	// Set Viper config params
	viper.AddConfigPath(home)
	viper.SetConfigName(".Dude")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	// Unmarshal config into struct
	c := &config{}
	err = viper.Unmarshal(c)
	cobra.CheckErr(err)
	return c, nil
}
