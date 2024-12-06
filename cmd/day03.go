package cmd

import (
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/dusktreader/advent-of-code-2024/util"
)

func init() {
	rootCmd.AddCommand(d3Cmd)
	d3Cmd.PersistentFlags().StringP("input-file", "i", "", "Get input from a file instead of stdin")
	d3Cmd.AddCommand(d3p1Cmd)
	d3Cmd.AddCommand(d3p2Cmd)
	d3p2Cmd.Flags().BoolP("annotate", "a", false, "Annotate redacted input string")
}

var d3Cmd = &cobra.Command{
	Use:   "day03",
	Short: "Day 3 Solutions",
	Long:  "The solutions for day 3 of Advent of Code 2024",
	Run:   d3Main,
}

var d3p1Cmd = &cobra.Command{
	Use:   "part1",
	Short: "Day 3, 1 Solution",
	Long:  "The solution for day 3, part 1 of Advent of Code 2024",
	Run:   d3p1Main,
}

var d3p2Cmd = &cobra.Command{
	Use:   "part2",
	Short: "Day 3, 2 Solution",
	Long:  "The solution for day 3, part 2 of Advent of Code 2024",
	Run:   d3p2Main,
}

func d3Main(cmd *cobra.Command, args []string){
	_ = cmd.Help()
}

func d3p1Main(cmd *cobra.Command, args []string) {
	inputStr, err := loadInput(cmd, args)
	MaybeDie(err)

	instructions := IsolatePairs(inputStr)
	total := ProcessPairs(instructions)

	slog.Debug("Results:", "total", total)
	fmt.Printf("%v\n", total)
}

func d3p2Main(cmd *cobra.Command, args []string) {
	annotate, err := cmd.Flags().GetBool("annotate")
	MaybeDie(err)

	inputStr, err := loadInput(cmd, args)
	MaybeDie(err)

	inputStr = Redact(inputStr, annotate)
	instructions := IsolatePairs(inputStr)
	total := ProcessPairs(instructions)

	slog.Debug("Results:", "total", total)
	fmt.Printf("%v\n", total)
}

func Redact(inputStr string, annotate bool) (outputStr string) {
	repl := ""
	if annotate {
		repl = "REDACTED"
	}
	re := regexp.MustCompile(`(?s)don't\(\).*?(?:do\(\)|$)`)
	outputStr = string(re.ReplaceAll([]byte(inputStr), []byte(repl)))
	slog.Debug("Redacted input:", "outputStr", outputStr)
	return

}

func IsolatePairs(inputStr string) (pairs []util.Pair[int]) {
	inputStr = strings.ReplaceAll(inputStr, "\n", "")
	// We could make the maximum digit count dynamic
	re := regexp.MustCompile(`mul\((\d{1,3}),(\d{1,3})\)`)
	matches := re.FindAllSubmatch([]byte(inputStr), -1)

	for _, match := range matches {
		left, err := strconv.Atoi(string(match[1]))
		if err != nil {
			slog.Debug("Skipping match due to failed integer conversion in left operand", "err", err)
			continue
		}

		right, err := strconv.Atoi(string(match[2]))
		if err != nil {
			slog.Debug("Skipping match due to failed integer conversion in right operand", "err", err)
			continue
		}
		pairs = append(pairs, util.Pair[int]{left, right})
	}
	return
}

func ProcessPairs(pairs []util.Pair[int]) (total int) {
	for i, pair := range pairs {
		slog.Debug("Processing pair:", "i", i, "pair", pair)
		val := pair.Left * pair.Right
		slog.Debug("Value is:", "val", val)
		total += val
		slog.Debug("New total is:", "total", total)
	}
	return
}

