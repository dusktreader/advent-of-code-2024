package cmd

import (
	"fmt"
	"log/slog"
	"math"
	"strconv"
	"strings"

	"github.com/dusktreader/advent-of-code-2024/util"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(d7Cmd)
	d7Cmd.PersistentFlags().StringP("input-file", "i", "", "Get input from a file instead of stdin")
	d7Cmd.PersistentFlags().StringP("output-file", "o", "", "Draw the final map and store it in the file")
	d7Cmd.AddCommand(d7p1Cmd)
	d7Cmd.AddCommand(d7p2Cmd)
}

var d7Cmd = &cobra.Command{
	Use:   "day07",
	Short: "Day 7 Solutions",
	Long:  "The solutions for day 7 of Advent of Code 2025",
	Run:   d7Main,
}

var d7p1Cmd = &cobra.Command{
	Use:   "part1",
	Short: "Day 7, 1 Solution",
	Long:  "The solution for day 7, part 1 of Advent of Code 2025",
	Run:   d7p1Main,
}

var d7p2Cmd = &cobra.Command{
	Use:   "part2",
	Short: "Day 7, 2 Solution",
	Long:  "The solution for day 7, part 2 of Advent of Code 2025",
	Run:   d7p2Main,
}

func d7Main(cmd *cobra.Command, args []string){
	_ = cmd.Help()
}

func d7p1Main(cmd *cobra.Command, args []string) {
	inputStr, err := loadInput(cmd, args)
	MaybeDie(err)

	eqs, err := ParseEquations(inputStr)
	MaybeDie(err)

	total := EqTotal2(eqs)
	slog.Debug("Results:", "Total", total)
	fmt.Printf("%v\n", total)
}

func d7p2Main(cmd *cobra.Command, args []string) {
	inputStr, err := loadInput(cmd, args)
	MaybeDie(err)

	eqs, err := ParseEquations(inputStr)
	MaybeDie(err)

	total := EqTotal3(eqs)
	slog.Debug("Results:", "Total", total)
	fmt.Printf("%v\n", total)
}

type Equation struct {
	Left  int
	Right []int
	Sats  util.Set[int]
}

func Cat(l int, r int) int {
	if l == 0 {
		return r
	} else if r == 0 {
		return l * 10
	} else {
		return l * int(math.Pow10(int(math.Log10(float64(r)) + 1))) + r
	}
}

func (eq Equation) Operate2(ops int) int {
	mask := 1 << (len(eq.Right) - 2)
	total := eq.Right[0]
	for i := 1; i < len(eq.Right); i++ {
		if mask & ops != 0 {
			total *= eq.Right[i]
		} else {
			total += eq.Right[i]
		}
		mask >>= 1
	}
	return total
}

const (
	ADD = 0
	MUL = 1
	CAT = 2
)

func (eq Equation) Operate3(ops int) int {
	exp := len(eq.Right) - 2
	total := eq.Right[0]
	for i := 1; i < len(eq.Right); i++ {
		op := (ops / util.Pow3(exp)) % 3
		switch op {
		case ADD:
			total += eq.Right[i]
		case MUL:
			total *= eq.Right[i]
		case CAT:
			total = Cat(total, eq.Right[i])
		}
		exp--
	}
	return total
}

func (eq Equation) Process2() {
	opMax := util.Pow2(len(eq.Right) - 1)

	for i := 0; i < opMax; i++ {
		if eq.Left == eq.Operate2(i) {
			eq.Sats.Add(i)
		}
	}
}

func (eq Equation) Process3() {
	opMax := util.Pow3(len(eq.Right) - 1)

	for i := 0; i < opMax; i++ {
		if eq.Left == eq.Operate3(i) {
			eq.Sats.Add(i)
		}
	}
}

func EqTotal2(eqs []Equation) int {
	total := 0
	for _, eq := range eqs {
		eq.Process2()
		if !eq.Sats.Empty() {
			total += eq.Left
		}
	}
	return total
}

func EqTotal3(eqs []Equation) int {
	total := 0
	for _, eq := range eqs {
		eq.Process3()
		if !eq.Sats.Empty() {
			total += eq.Left
		}
	}
	return total
}

func ParseEquations(inputStr string) ([]Equation, error) {
	inputStr = strings.TrimSpace(inputStr)
	lines := strings.Split(inputStr, "\n")
	eqs := make([]Equation, len(lines))
	for i, line := range lines {
		line = strings.TrimSpace(line)
		parts := strings.Split(line, ":")

		left, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, util.ReErr(err, "Couldn't convert left side %v", parts[0])
		}

		parts = strings.Split(strings.TrimSpace(parts[1]), " ")
		right := make([]int, len(parts))
		for j, p := range parts {
			r, err := strconv.Atoi(p)
			if err != nil {
				return nil, util.ReErr(err, "Couldn't convert line %v right side %v at index %v", i, p, j)
			}
			right[j] = r
		}

		eqs[i] = Equation{
			Left: left,
			Right: right,
			Sats:  util.MakeSet[int](),
		}
	}
	return eqs, nil
}
