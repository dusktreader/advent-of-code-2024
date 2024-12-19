package cmd

import (
	"fmt"
	"log/slog"
	"math"
	"regexp"
	"strconv"

	"github.com/dusktreader/advent-of-code-2024/util"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(d13Cmd)
	d13Cmd.PersistentFlags().StringP("input-file", "i", "", "Get input from a file instead of stdin")
	d13Cmd.AddCommand(d13p1Cmd)
	d13Cmd.AddCommand(d13p2Cmd)
}

var d13Cmd = &cobra.Command{
	Use:   "day13",
	Short: "Day 13 Solutions",
	Long:  "The solutions for day 13 of Advent of Code 2025",
	Run:   d13Main,
}

var d13p1Cmd = &cobra.Command{
	Use:   "part1",
	Short: "Day 13, 1 Solution",
	Long:  "The solution for day 13, part 1 of Advent of Code 2025",
	Run:   d13p1Main,
}

var d13p2Cmd = &cobra.Command{
	Use:   "part2",
	Short: "Day 13, 2 Solution",
	Long:  "The solution for day 13, part 2 of Advent of Code 2025",
	Run:   d13p2Main,
}

func d13Main(cmd *cobra.Command, args []string){
	_ = cmd.Help()
}

func d13p1Main(cmd *cobra.Command, args []string) {
	inputStr, err := loadInput(cmd, args)
	MaybeDie(err)

	slog.Debug("Parsing buttons")
	buttons, err := ParseButtons(inputStr)
	MaybeDie(err)

	slog.Debug("Counting tokens")
	count := CountTokens(buttons)

	slog.Debug("Results:", "TokenCount", count)
	fmt.Printf("%v\n", count)
}

func d13p2Main(cmd *cobra.Command, args []string) {
	inputStr, err := loadInput(cmd, args)
	MaybeDie(err)

	slog.Debug("Parsing buttons")
	buttons, err := ParseButtons(inputStr)
	MaybeDie(err)

	slog.Debug("Moving prize")
	MovePrize(&buttons)

	slog.Debug("Counting tokens")
	count := CountTokens(buttons)

	slog.Debug("Results:", "TokenCount", count)
	fmt.Printf("%v\n", count)
}

type Button struct {
	Prize util.Point
	A util.Vector
	B util.Vector
}

func (b Button) CompCounts() (t int, u int, err error) {
	det := b.A.Cross(b.B.Neg())
	if det == 0 {
		err = fmt.Errorf("Rays are parallel or coincident")
		return
	}

	o := util.MakePoint(0, 0)
	v := b.Prize.Diff(o)
	slog.Debug("Diff", "V", v)

	t = v.Cross(b.B.Neg()) / det
	u = v.Cross(b.A) / det
	slog.Debug("t, u", "t", t, "u", u)

	if t < 0 || u < 0 {
		err = fmt.Errorf("Rays do not intersect")
		return
	}

	if o.Add(b.A.Mul(t)).Add(b.B.Mul(u)) != b.Prize  {
		err = fmt.Errorf("No solution found")
		return
	}

	return
}

func (b Button) Count() int {
	t, u, err := b.CompCounts()
	if err != nil {
		slog.Debug("Error computing counts", "Error", err)
		return math.MaxInt
	}

	return t * 3 + u
}

func CountTokens(buttons []Button) (count int) {
	for _, b := range buttons {
		bct := b.Count()
		if bct < math.MaxInt {
			count += bct
		}
	}
	return
}

func MovePrize(buttons *[]Button) {
	for i := 0; i < len(*buttons); i++ {
		(*buttons)[i].Prize.I += 10000000000000
		(*buttons)[i].Prize.J += 10000000000000
	}
}

func ParseButtons(inputStr string) ([]Button, error) {
	re, err := regexp.Compile(
		`Button\s+A:\s+X\+(\d+),\s+Y\+(\d+)\s+Button\s+B:\s+X\+(\d+),\s+Y\+(\d+)\s+Prize:\s+X=(\d+),\s+Y=(\d+)`,
	)
	if err != nil {
		return nil, util.ReErr(err, "Couldn't compile regex")
	}

	buttons := make([]Button, 0)

	for _, match := range re.FindAllStringSubmatch(inputStr, -1) {
		ints := [6]int{}
		for i := 0; i < 6; i++ {
			val, err := strconv.Atoi(match[i + 1])
			if err != nil {
				return nil, util.ReErr(err, "Failed to convert int: %v", match[i + 1])
			}
			ints[i] = val
		}

		buttons = append(
			buttons,
			Button{
				A: util.MakeVector(ints[0], ints[1]),
				B: util.MakeVector(ints[2], ints[3]),
				Prize: util.MakePoint(ints[4], ints[5]),
			},
		)
	}
	return buttons, nil
}

