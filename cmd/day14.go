package cmd

import (
	"fmt"
	"log/slog"
	"regexp"
	"strconv"

	"github.com/dusktreader/advent-of-code-2024/util"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(d14Cmd)
	d14Cmd.PersistentFlags().StringP("input-file", "i", "", "Get input from a file instead of stdin")
	d14Cmd.PersistentFlags().BoolP("small", "s", false, "Use smaller space")
	d14Cmd.PersistentFlags().BoolP("visualize", "V", false, "Visualize space")
	d14Cmd.AddCommand(d14p1Cmd)
	d14Cmd.AddCommand(d14p2Cmd)
}

var d14Cmd = &cobra.Command{
	Use:   "day14",
	Short: "Day 14 Solutions",
	Long:  "The solutions for day 14 of Advent of Code 2025",
	Run:   d14Main,
}

var d14p1Cmd = &cobra.Command{
	Use:   "part1",
	Short: "Day 14, 1 Solution",
	Long:  "The solution for day 14, part 1 of Advent of Code 2025",
	Run:   d14p1Main,
}

var d14p2Cmd = &cobra.Command{
	Use:   "part2",
	Short: "Day 14, 2 Solution",
	Long:  "The solution for day 14, part 2 of Advent of Code 2025",
	Run:   d14p2Main,
}

func d14Main(cmd *cobra.Command, args []string){
	_ = cmd.Help()
}

func d14p1Main(cmd *cobra.Command, args []string) {
	inputStr, err := loadInput(cmd, args)
	MaybeDie(err)

	slog.Debug("Parsing robots")
	ps, err := ParseRobots(inputStr)
	MaybeDie(err)

	small, err := cmd.Flags().GetBool("small")
	MaybeDie(err)

	if small {
		ps.Size = util.Size{W: 11, H: 7}
	} else {
		ps.Size = util.Size{W: 101, H: 103}
	}

	slog.Debug("Moving robots")
	ps.MoveRobots(100)

	visualize, err := cmd.Flags().GetBool("visualize")
	MaybeDie(err)

	if visualize {
		fmt.Printf("%v\n", ps.Viz())
	}

	slog.Debug("Computing safety")
	safety := ps.ComputeSafety()

	slog.Debug("Results:", "Safety", safety)
	fmt.Printf("%v\n", safety)
}

func d14p2Main(cmd *cobra.Command, args []string) {
	inputStr, err := loadInput(cmd, args)
	MaybeDie(err)

	slog.Debug("Parsing robots")
	ps, err := ParseRobots(inputStr)
	MaybeDie(err)

	ps.Size = util.Size{W: 101, H: 103}

	for i := range 10_000 {
		fmt.Printf("\nIteration %v\n", i)
		fmt.Printf("%v\n\n\n", ps.Viz())
		ps.MoveRobots(1)
	}
}

type Robot struct {
	Pos util.Point
	Vel util.Vector
}

type PissSpace struct {
	Size util.Size
	Robots []Robot
}

func (ps *PissSpace) MoveRobots(t int) {

	var mod func(int, int) int
	mod = func(a, b int) int {
		return (a % b + b) % b
	}
	for i := range ps.Robots {
		ps.Robots[i].Pos = ps.Robots[i].Pos.Add(ps.Robots[i].Vel.Mul(t))
		ps.Robots[i].Pos.I = mod(ps.Robots[i].Pos.I, ps.Size.H)
		ps.Robots[i].Pos.J = mod(ps.Robots[i].Pos.J, ps.Size.W)
	}
}

func (ps *PissSpace) Viz() string {
	grid := make([][]rune, ps.Size.H)
	for i := 0; i < ps.Size.H; i++ {
		grid[i] = make([]rune, ps.Size.W)
		for j := 0; j < ps.Size.W; j++ {
			grid[i][j] = '.'
		}
	}

	counter := util.MakeCounter[util.Point]()
	for _, r := range ps.Robots {
		slog.Debug("Counting robot at:", "Pos", r.Pos)
		counter.Incr(r.Pos)
	}

	for p, c := range counter.Iter() {
		v, err := util.ItoR(c)
		if err != nil {
			slog.Error("Couldn't convert int to rune", "Error", err)
			continue
		}
		grid[p.I][p.J] = v
	}

	out := ""
	for i := 0; i < ps.Size.H; i++ {
		for j := 0; j < ps.Size.W; j++ {
			out += string(grid[i][j])
		}
		out += "\n"
	}
	return out
}

func (ps *PissSpace) ComputeSafety() int {
	quads := make([]util.Rect, 4)
	quads[0] = util.Rect{
		O: util.MakePoint(0, 0),
		Sz: ps.Size.Div(2),
	}
	quads[1] = util.Rect{
		O: util.MakePoint(0, ps.Size.W / 2 + 1),
		Sz: ps.Size.Div(2),
	}
	quads[2] = util.Rect{
		O: util.MakePoint(ps.Size.H / 2 + 1, 0),
		Sz: ps.Size.Div(2),
	}
	quads[3] = util.Rect{
		O: util.MakePoint(ps.Size.H / 2 + 1, ps.Size.W / 2 + 1),
		Sz: ps.Size.Div(2),
	}
	slog.Debug("Quads", "Quads", quads)

	prod := 1
	for _, q := range quads {
		count := 0
		for _, r := range ps.Robots {
			if q.Has(r.Pos) {
				count++
			}
		}
		prod *= count
	}

	return prod
}

func ParseRobots(inputStr string) (*PissSpace, error) {
	re, err := regexp.Compile(
		`p=(\d+),(\d+) v=(-?\d+),(-?\d+)`,
	)
	if err != nil {
		return nil, util.ReErr(err, "Couldn't compile regex")
	}

	pisser := PissSpace{
		Size: util.Size{W: 0, H: 0},
		Robots: make([]Robot, 0),
	}

	for _, match := range re.FindAllStringSubmatch(inputStr, -1) {
		ints := [4]int{}
		for i := 0; i < 4; i++ {
			val, err := strconv.Atoi(match[i + 1])
			if err != nil {
				return nil, util.ReErr(err, "Failed to convert int: %v", match[i + 1])
			}
			ints[i] = val
		}

		rob := Robot{
			Pos: util.MakePoint(ints[1], ints[0]),
			Vel: util.MakeVector(ints[3], ints[2]),
		}

		pisser.Robots = append(pisser.Robots, rob)
	}
	return &pisser, nil
}

