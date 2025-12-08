// Package cmd implements the CLI commands for pgedge-anonymizer.
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/pgedge/pgedge-anonymizer/internal/version"
)

var (
	cfgFile       string
	quiet         bool
	configLoadErr error
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pgedge-anonymizer",
	Short: "Anonymize PII data in PostgreSQL databases",
	Long: `pgEdge Anonymizer is a tool for anonymizing Personally Identifiable
Information (PII) in PostgreSQL databases to meet GDPR and other regulatory
requirements when cloning production data for development purposes.

The tool processes columns specified in a configuration file, applying
pattern-based anonymization while maintaining data consistency across tables.`,
	SilenceUsage: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "",
		"config file (default: ./pgedge-anonymizer.yaml or /etc/pgedge/pgedge-anonymizer.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false,
		"suppress progress output")

	// Bind flags to viper
	_ = viper.BindPFlag("quiet", rootCmd.PersistentFlags().Lookup("quiet"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Skip config loading for commands that don't need it
	// The version and help commands should work without any config file
	if len(os.Args) > 1 {
		cmd := os.Args[1]
		if cmd == "version" || cmd == "help" || cmd == "--help" || cmd == "-h" {
			return
		}
	}

	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Search for config file explicitly with .yaml extension
		// to avoid viper matching other files (like the binary itself)
		configName := "pgedge-anonymizer.yaml"
		searchPaths := []string{"."}

		// Add /etc/pgedge
		searchPaths = append(searchPaths, "/etc/pgedge")

		// Add directory containing the binary
		if exe, err := os.Executable(); err == nil {
			searchPaths = append(searchPaths, filepath.Dir(exe))
		}

		// Search for config file in each path
		var foundConfig string
		for _, dir := range searchPaths {
			path := filepath.Join(dir, configName)
			if _, err := os.Stat(path); err == nil {
				foundConfig = path
				break
			}
		}

		if foundConfig != "" {
			viper.SetConfigFile(foundConfig)
		} else {
			// No config found - set up viper for error reporting
			viper.SetConfigName("pgedge-anonymizer")
			viper.SetConfigType("yaml")
			viper.AddConfigPath(".")
		}
	}

	// Read environment variables with PGANON_ prefix
	viper.SetEnvPrefix("PGANON")
	viper.AutomaticEnv()

	// Read config file if found
	if err := viper.ReadInConfig(); err != nil {
		configLoadErr = err
	}
}

// CheckConfigLoaded returns an error if no config file was loaded.
// Commands that require configuration should call this.
func CheckConfigLoaded() error {
	if configLoadErr != nil {
		if _, ok := configLoadErr.(viper.ConfigFileNotFoundError); ok {
			if cfgFile != "" {
				return fmt.Errorf("config file not found: %s", cfgFile)
			}
			return fmt.Errorf("no config file found. Create pgedge-anonymizer.yaml or specify one with --config")
		}
		// Include which file caused the error if available
		if file := viper.ConfigFileUsed(); file != "" {
			return fmt.Errorf("error reading config file %s: %w", file, configLoadErr)
		}
		return fmt.Errorf("error reading config file: %w", configLoadErr)
	}
	return nil
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("pgedge-anonymizer %s (built %s)\n",
			version.Version, version.BuildTime)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
