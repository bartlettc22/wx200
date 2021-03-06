package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "wx200",
	Short: "Prometheus exporter for the Radioshack WX200 Weather Station",
	Long:  `Prometheus exporter for the Radioshack WX200 Weather Station`,
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of wx200",
	Long:  `All software has versions. This is wx200's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version)
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&listenPort, "listen-port", "p", "9041", "Metrics server listener port")
	rootCmd.PersistentFlags().StringVarP(&comPortName, "com-port", "c", "/dev/ttyUSB0", "Com port that the WX200 is attached")
	rootCmd.PersistentFlags().StringVarP(&v, "verbosity", "v", log.InfoLevel.String(), "Log level (debug, info, warn")
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if err := setUpLogs(os.Stdout, v); err != nil {
			return err
		}
		return nil
	}
	rootCmd.AddCommand(versionCmd)
}

func setUpLogs(out io.Writer, level string) error {
	log.SetOutput(out)
	lvl, err := log.ParseLevel(level)
	if err != nil {
		return err
	}
	log.SetLevel(lvl)
	log.Debugf("Log level set to: %s", lvl)
	return nil
}
