package cmd

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/dusktreader/advent-of-code-2024/graph"
	"github.com/dusktreader/advent-of-code-2024/util"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(d16Cmd)
	d16Cmd.PersistentFlags().StringP("input-file", "i", "", "Get input from a file instead of stdin")
	d16Cmd.PersistentFlags().CountP("visualize", "V", "Visualize maze. Pass multiple to visualize more")
	d16Cmd.AddCommand(d16p1Cmd)
	d16Cmd.AddCommand(d16p2Cmd)
}

var d16Cmd = &cobra.Command{
	Use:   "day16",
	Short: "Day 16 Solutions",
	Long:  "The solutions for day 16 of Advent of Code 2025",
	Run:   d16Main,
}

var d16p1Cmd = &cobra.Command{
	Use:   "part1",
	Short: "Day 16, 1 Solution",
	Long:  "The solution for day 16, part 1 of Advent of Code 2025",
	Run:   d16p1Main,
}

var d16p2Cmd = &cobra.Command{
	Use:   "part2",
	Short: "Day 16, 2 Solution",
	Long:  "The solution for day 16, part 2 of Advent of Code 2025",
	Run:   d16p2Main,
}

func d16Main(cmd *cobra.Command, args []string){
	_ = cmd.Help()
}

func d16p1Main(cmd *cobra.Command, args []string) {
	//viz := getViz(cmd, args)

	inputStr, err := loadInput(cmd, args)
	MaybeDie(err)

	slog.Debug("Parsing maze")
	mz, err := ParseMaze(inputStr)
	MaybeDie(err)

	mz.Simplify()

	fmt.Printf("%v\n", mz)

	// slog.Debug("Results:", "Score", mz.BestDeer.Score)
	// fmt.Printf("%v\n", mz.BestDeer.Score)
}

func d16p2Main(cmd *cobra.Command, args []string) {
}

type Maze struct {
	Size      util.Size
	Walls     util.Set[util.Point]
	Start     util.Point
	End       util.Point
	Graph     graph.Graph[util.Point]
	Hilite    util.Point
}

func (mz *Maze) RenderPt(pt util.Point) rune {
	if mz.Walls.Has(pt) {
		return '█'
	} else if mz.Start == pt {
		return 'S'
	} else if mz.End == pt {
		return 'E'
	} else if mz.Hilite == pt {
		return 'X'
	} else if mz.Graph.Has(pt) {
		return 'o'
	}

	return ' '
}

func (mz *Maze) String() string {
	var sb strings.Builder
	//sb.WriteString(fmt.Sprintf("Best score: %v, Current score: %v\n", mz.BestDeer.Score, mz.Deer.Score))
	for i := 0; i < mz.Size.H; i++ {
		for j := 0; j < mz.Size.W; j++ {
			pt := util.MakePoint(i, j)
			sb.WriteRune(mz.RenderPt(pt))
		}
		sb.WriteRune('\n')
	}

	sb.WriteRune('\n')

	return sb.String()
}

func (mz *Maze) Simplify() {
	corners  := util.MakeSet[util.Point]()
	straights := util.MakeSet[util.Point]()
	deadEnds := util.MakeSet[util.Point]()
	for node := range mz.Graph.Nodes().Iter() {
		if node == mz.Start || node == mz.End {
			continue
		}

		nbors := mz.Graph.Nbors(node)
		if nbors.Size() == 1 {
			deadEnds.Add(node)
			continue
		}

		if nbors.Size() > 2 {
			corners.Add(node)
			continue
		}

		a := nbors.Pop()
		b := nbors.Pop()
		if a.I != b.I && a.J != b.J {
			corners.Add(node)
			continue
		}

		straights.Add(node)
	}

	for !corners.Empty() {
		node := corners.Pop()
		if node == mz.Start || node == mz.End {
			continue
		}

		nbors := mz.Graph.Nbors(node)
		for !nbors.Empty() {
			a := nbors.Pop()
			straights.Rem(a)
			for b := range nbors.Iter() {
				if a.I == b.I || a.J == b.J {
					mz.Graph.AddEdge(a, b, 2)
				} else {
					mz.Graph.AddEdge(a, b, 1002)
				}
			}
		}
		mz.Graph.Rem(node)
	}

	for !deadEnds.Empty() {
		node := deadEnds.Pop()
		if node == mz.Start || node == mz.End {
			continue
		}

		nbors := mz.Graph.Nbors(node)
		nbor := nbors.Pop()
		if mz.Graph.Nbors(nbor).Size() == 2 {
			deadEnds.Add(nbor)
		}
		mz.Graph.Rem(node)
	}

	for !straights.Empty() {
		b := straights.Pop()
		if !mz.Graph.Has(b) {
			continue
		}

		if b == mz.Start || b == mz.End {
			continue
		}

		nbors := mz.Graph.Nbors(b)
		a := nbors.Pop()
		c := nbors.Pop()
		ac, ok := mz.Graph.Edge(a, b)
		if !ok {
			panic("No edge from a to b. This should not happen")
		}

		bc, ok := mz.Graph.Edge(b, c)
		if !ok {
			panic("No edge from b to c. This should not happen")
		}

		mz.Graph.AddEdge(a, c, ac.Wt + bc.Wt)
		mz.Graph.Rem(b)
	}
}

func ParseMaze(inputStr string) (*Maze, error) {
	inputStr = strings.TrimSpace(inputStr)
	lines := strings.Split(inputStr, "\n")

	if len(lines) < 1 {
		return nil, fmt.Errorf("No lines in input")
	}

	mz := Maze{
		Size:      util.Size{W: len(lines[0]), H: len(lines)},
		Walls:     util.MakeSet[util.Point](),
		Graph:     *graph.MakeGraph[util.Point](false),
	}

	for i, line := range lines {
		if len(line) != mz.Size.W {
			return nil, fmt.Errorf("Line %d is not the same length", i)
		}
		for j, rn := range line {
			pt := util.MakePoint(i, j)
			if rn == '#' {
				mz.Walls.Add(pt)
				continue
			} else if rn == 'S' {
				mz.Start = pt
			} else if rn == 'E' {
				mz.End = pt
			}

			mz.Graph.Add(pt)

			nPt := pt.Add(util.NORTH)
			if mz.Graph.Has(nPt) {
				mz.Graph.AddEdge(pt, nPt)
			}
			wPt := pt.Add(util.WEST)
			if mz.Graph.Has(wPt) {
				mz.Graph.AddEdge(pt, wPt)
			}
		}
	}

	return &mz, nil
}

