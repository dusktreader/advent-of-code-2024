package cmd

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Show verbose logging output")
}

var rootCmd = &cobra.Command{
	Use:   "aoc",
	Short: "Advent of Code - 2024",
	Long:  "The Advent of Code submission for dusktreader@github.com in 2024",
	PersistentPreRun: preRun,
	Run:   rootMain,
}

func preRun(cmd *cobra.Command, args []string) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	verbose, err := cmd.Flags().GetBool("verbose")
	MaybeDie(err)
	if verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	} else {
		slog.SetLogLoggerLevel(slog.LevelInfo)
	}
}

func rootMain(cmd *cobra.Command, args []string) {
	_ = cmd.Help()
}

func MaybeDie(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "There was an error:", err)
		os.Exit(1)
	}
}

func Die(msg string) {
	fmt.Fprintln(os.Stderr, "Aborting:", msg)
	os.Exit(1)
}


func Execute() {
	err := rootCmd.Execute()
	MaybeDie(err)
}
