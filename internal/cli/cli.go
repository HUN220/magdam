package cli

import (
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"os"
)

type CliOptions struct {
	Command    string
	ConfigFile string
	Continue   bool   `json:"continue"`
	BaseUrl    string `json:"baseUrl"`
	ApiKeyId   string `json:"apiKeyId"`
	ApiKey     string `json:"apiKey"`
}

var (
	pullCmd = flag.NewFlagSet("pull", flag.ExitOnError)
	pushCmd = flag.NewFlagSet("push", flag.ExitOnError)
)

var subcommands = map[string]*flag.FlagSet{
	pullCmd.Name(): pullCmd,
	pushCmd.Name(): pushCmd,
}

// Attach global flags to each subcommand
func setupCommonFlags(opts *CliOptions) {
	for _, fs := range subcommands {
		fs.StringVar(
			&opts.ConfigFile,
			"config",
			"",
			"load configuration options from a JSON config file",
		)

		fs.StringVar(
			&opts.ConfigFile,
			"f",
			"",
			"load configuration options from a JSON config file  (shorthand)",
		)

		fs.StringVar(
			&opts.BaseUrl,
			"url",
			"",
			"base url for Magda API",
		)

		fs.StringVar(
			&opts.ApiKeyId,
			"api-key-id",
			"",
			"API Key ID for Magda API",
		)

		fs.StringVar(
			&opts.ApiKey,
			"api-key",
			"",
			"API Key for Magda API",
		)

		fs.BoolVar(
			&opts.Continue,
			"continue",
			false,
			"continue on error",
		)

		fs.BoolVar(
			&opts.Continue,
			"c",
			false,
			"continue on error (shorthand)",
		)
	}
}

// Collects subcommand and it's options from cli arguments and config files
func NewOptions(opts *CliOptions) error {
	// Set up subcommands
	cmdError := "expected 'pull' or 'push' subcommands"
	if len(os.Args) < 2 {
		return errors.New(cmdError)
	}

	cmd := subcommands[os.Args[1]]

	if cmd == nil {
		return errors.New(cmdError)
	}
	// Set up flags
	setupCommonFlags(opts)

	cmd.Parse(os.Args[2:])

	// Populate CliOptions
	opts.Command = cmd.Name()
	if opts.ConfigFile != "" {
		err := loadConfigFile(opts)
		if err != nil {
			return err
		}
	}

	// Check required options are set
	if opts.ApiKeyId == "" {
		return errors.New("magda API Key ID required")
	}

	if opts.ApiKey == "" {
		return errors.New("magda API Key required")
	}

	if opts.BaseUrl == "" {
		return errors.New("magda API url required")
	}

	return nil
}

// Loads configuration values from a config file.
// Path to the config file is defined in the passed CliOptions.
func loadConfigFile(opts *CliOptions) error {
	// Open our jsonFile
	jsonFile, err := os.Open(opts.ConfigFile)
	// if we os.Open returns an error then handle it
	if err != nil {
		return err
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &opts)

	return nil
}
