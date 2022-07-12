package cmd

import (
	"fmt"
	"github.com/chain710/awesome-home/cmd/values"
	"github.com/chain710/awesome-home/internal/log"
	"go.uber.org/zap/zapcore"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile  string
	logLevel = values.LogLevel{Level: zapcore.ErrorLevel}
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Short: "home commands collection",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.home.yaml)")
	rootCmd.PersistentFlags().VarP(&logLevel, "log-level", "L", "log level: error|warn|info|debug")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	log.Init(log.WithLogLevel(logLevel.Level))
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".home" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".home")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		_, _ = fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
