package cmd

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/dusktreader/advent-of-code-2024/util"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(day1Cmd)
	day1Cmd.Flags().StringP("input-file", "i", "", "Get input from a file instead of stdin")
}

var day1Cmd = &cobra.Command{
	Use: "day01-1",
	Short: "Day 1 Solution",
	Long:  "The solution for day 1, part 1 of Advent of Code 2024",
	Run:   day1Main,
}

func day1Main(cmd *cobra.Command, args []string){
	inputFile, err := cmd.Flags().GetString("input-file")
	MaybeDie(err)

	var input []byte

	if inputFile != "" {
		slog.Debug("Input file provided. Reading from file", "file", inputFile)
		input, err = os.ReadFile(inputFile)
	} else {
		slog.Debug("No input file provided. Reading from stdin")
		input, err = io.ReadAll(os.Stdin)
	}
	MaybeDie(err)

	inputStr := string(input)
	if inputStr == "" {
		Die("Didn't get any input")
	}

	inputStr = strings.TrimSpace(inputStr)
	inputStr = strings.ReplaceAll(inputStr, "\n", " ")

	var leftNumbers []int
	var rightNumbers []int
	isLeft := true

	tokens := strings.Split(inputStr, " ")
	for _, token := range tokens {
		if token != "" {
			number, err := strconv.Atoi(token)
			MaybeDie(err)
			if isLeft {
				leftNumbers = append(leftNumbers, number)
			} else {
				rightNumbers = append(rightNumbers, number)
			}
			isLeft = !isLeft
		}
	}

	distance, err := TotalDistance(leftNumbers, rightNumbers)
	MaybeDie(err)

	slog.Debug("Results:", "distance", distance, "left", leftNumbers, "right", rightNumbers)
	fmt.Printf("%v\n", distance)
}

func TotalDistance(leftNumbers []int, rightNumbers []int) (int, error) {
	if len(leftNumbers) != len(rightNumbers) {
		return 0, fmt.Errorf("Need an even number of number inputs")
	}

	sort.Ints(leftNumbers)
	sort.Ints(rightNumbers)

	var total int

	for i := 0; i < len(leftNumbers); i++ {
		total += util.AbsInt(leftNumbers[i] - rightNumbers[i])
	}
	return total, nil
}
