package cmd

import (
	"fmt"
	"io"
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

func loadInput(cmd *cobra.Command, args []string) (inputStr string, err error) {
	inputFile, err := cmd.Flags().GetString("input-file")
	if err != nil {
		return "", fmt.Errorf("Couldn't get input-file argument: %#v", err)
	}

	var input []byte

	if inputFile != "" {
		slog.Debug("Input file provided. Reading from file", "file", inputFile)
		input, err = os.ReadFile(inputFile)
	} else {
		slog.Debug("No input file provided. Reading from stdin")
		input, err = io.ReadAll(os.Stdin)
	}
	if err != nil {
		return "", fmt.Errorf("Couldn't read input: %#v", err)
	}

	inputStr = string(input)
	if inputStr == "" {
		return "", fmt.Errorf("Didn't get any input")
	}
	return inputStr, nil
}

func MaybeDie(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "There was an error:", err)
		os.Exit(1)
	}
}

func Die(msg string, flags ...interface{}) {
	msg = fmt.Sprintf(msg, flags...)
	fmt.Fprintln(os.Stderr, "Aborting:", msg)
	os.Exit(1)
}

func Execute() {
	err := rootCmd.Execute()
	MaybeDie(err)
}
