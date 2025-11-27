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
	cfgFile string
	quiet   bool
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
	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Search for config in standard locations
		viper.SetConfigName("pgedge-anonymizer")
		viper.SetConfigType("yaml")

		// Current directory
		viper.AddConfigPath(".")

		// /etc/pgedge
		viper.AddConfigPath("/etc/pgedge")

		// Directory containing the binary
		if exe, err := os.Executable(); err == nil {
			viper.AddConfigPath(filepath.Dir(exe))
		}
	}

	// Read environment variables with PGANON_ prefix
	viper.SetEnvPrefix("PGANON")
	viper.AutomaticEnv()

	// Read config file if found
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Fprintf(os.Stderr, "Error reading config file: %v\n", err)
		}
	}
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
