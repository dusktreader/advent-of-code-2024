package cmd

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/dusktreader/advent-of-code-2024/util"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(d6Cmd)
	d6Cmd.PersistentFlags().StringP("input-file", "i", "", "Get input from a file instead of stdin")
	d6Cmd.AddCommand(d6p1Cmd)
	d6Cmd.AddCommand(d6p2Cmd)
	d6p2Cmd.Flags().BoolP("annotate", "a", false, "Annotate redacted input string")
	d6Cmd.AddCommand(mermaidCmd)
	d6Cmd.AddCommand(dotCmd)
	d6Cmd.AddCommand(validateCmd)
}

var d6Cmd = &cobra.Command{
	Use:   "day06",
	Short: "Day 6 Solutions",
	Long:  "The solutions for day 5 of Advent of Code 2025",
	Run:   d6Main,
}

var d6p1Cmd = &cobra.Command{
	Use:   "part1",
	Short: "Day 6, 1 Solution",
	Long:  "The solution for day 5, part 1 of Advent of Code 2025",
	Run:   d6p1Main,
}

var d6p2Cmd = &cobra.Command{
	Use:   "part2",
	Short: "Day 6, 2 Solution",
	Long:  "The solution for day 5, part 2 of Advent of Code 2025",
	Run:   d6p2Main,
}

func d6Main(cmd *cobra.Command, args []string){
	_ = cmd.Help()
}

func d6p1Main(cmd *cobra.Command, args []string) {
	inputStr, err := loadInput(cmd, args)
	MaybeDie(err)

	manual, err := ParseInput(inputStr)
	MaybeDie(err)

	manual.Validate()
	slog.Debug("Results:", "ValidCheckSum", manual.ValidCheckSum)
	fmt.Printf("%v\n", manual.ValidCheckSum)
}

func d6p2Main(cmd *cobra.Command, args []string) {
}

type mark int

const (
	Guard mark = iota + 1
	Obstr
	Empty
	Tread
)

type LabMap struct {
	GuardPos util.Point
	GuardDir util.Point
	Grid     []mark
	GridSz   util.Point
}

func ParseLabMap(inputStr string) (LabMap, error) {
	inputStr = strings.TrimSpace(inputStr)
	manual := MakeManual()

	lines := strings.Split(inputStr, "\n")

	lm := LabMap{}
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if lm.Grid == nil {
			lm.GridSz = util.MakePoint(len(lines), len(line))
			lm.Grid = make([]mark, lm.GridSz.I * lm.GridSz.J)
		} else if len(line) != lm.GridSz.J {
			return lm, fmt.Errorf("Line #%v has a different size: %v", i, len(line))
		}
		for j, rn := range line {

		}
	}

}
