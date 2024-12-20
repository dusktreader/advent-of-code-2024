package cmd

import (
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"

	"github.com/dusktreader/advent-of-code-2024/util"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(d15Cmd)
	d15Cmd.PersistentFlags().StringP("input-file", "i", "", "Get input from a file instead of stdin")
	d15Cmd.PersistentFlags().BoolP("small", "s", false, "Use smaller space")
	d15Cmd.PersistentFlags().BoolP("visualize", "V", false, "Visualize space")
	d15Cmd.AddCommand(d15p1Cmd)
	d15Cmd.AddCommand(d15p2Cmd)
}

var d15Cmd = &cobra.Command{
	Use:   "day15",
	Short: "Day 15 Solutions",
	Long:  "The solutions for day 15 of Advent of Code 2025",
	Run:   d15Main,
}

var d15p1Cmd = &cobra.Command{
	Use:   "part1",
	Short: "Day 15, 1 Solution",
	Long:  "The solution for day 15, part 1 of Advent of Code 2025",
	Run:   d15p1Main,
}

var d15p2Cmd = &cobra.Command{
	Use:   "part2",
	Short: "Day 15, 2 Solution",
	Long:  "The solution for day 15, part 2 of Advent of Code 2025",
	Run:   d15p2Main,
}

func d15Main(cmd *cobra.Command, args []string){
	_ = cmd.Help()
}

func d15p1Main(cmd *cobra.Command, args []string) {
	inputStr, err := loadInput(cmd, args)
	MaybeDie(err)

	slog.Debug("Parsing warehouse")
	wh, err := ParseWarehouse(inputStr)
	MaybeDie(err)

	safety := 0

	slog.Debug("Results:", "Safety", safety)
	fmt.Printf("%v\n", safety)
}

func d15p2Main(cmd *cobra.Command, args []string) {
}

type Warehouse struct {
	Sz    util.Size
	Walls util.Set[util.Point]
	Boxes util.Set[util.Point]
	Robot util.Point
	Moves util.Queue[util.Vector]
}

func (wh *Warehouse) MoveRobot() error {
	v, err := wh.Moves.Pop()
	if err != nil {
		return fmt.Errorf("No moves left")
	}

	boxLine := util.MakeSet[util.Point]()
	f := wh.Robot.Add(v)
	for wh.Boxes.Has(f) {
		wh.Boxes.Rem(f)
		boxLine.Add(f)
		f.Add(v)
	}
	if wh.Walls.Has(f) {
		v = util.MakeVector(0, 0)
	}
	for pt := range boxLine.Iter() {
		wh.Boxes.Add(pt.Add(v))
	}
	wh.Robot = wh.Robot.Add(v)
	return nil
}

func ParseWarehouse(inputStr string) (*Warehouse, error) {
	inputStr = strings.TrimSpace(inputStr)
	lines := strings.Split(inputStr, "\n")

	parsingMap := true
	wh := Warehouse{
		Sz: util.Size{W: len(lines[0]), H: len(lines)},
		Walls: util.MakeSet[util.Point](),
		Boxes: util.MakeSet[util.Point](),
		Moves: *util.MakeQueue[util.Vector](),
		Robot: util.MakePoint(0, 0),
	}

	for i, line := range lines {
		if line == "" {
			parsingMap = false
			continue
		}

		if parsingMap {
			for j, rn := range line {
				pt := util.MakePoint(i, j)
				switch rn {
				case '#':
					wh.Walls.Add(pt)
				case 'O':
					wh.Boxes.Add(pt)
				case '@':
					wh.Robot = pt
				}
			}
		} else {
			for _, rn := range line {
				switch rn {
				case '^':
					wh.Moves.Push(util.NORTH)
				case '>':
					wh.Moves.Push(util.EAST)
				case 'v':
					wh.Moves.Push(util.SOUTH)
				case '<':
					wh.Moves.Push(util.WEST)
				default:
					return nil, fmt.Errorf("Unknown direction in moves: %v", rn)
				}
			}
		}
	}

	return &wh, nil
}

