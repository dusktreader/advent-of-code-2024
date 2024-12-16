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
	rootCmd.AddCommand(d10Cmd)
	d10Cmd.PersistentFlags().StringP("input-file", "i", "", "Get input from a file instead of stdin")
	d10Cmd.AddCommand(d10p1Cmd)
	d10Cmd.AddCommand(d10p2Cmd)
	d10Cmd.AddCommand(showCmd)
}

var d10Cmd = &cobra.Command{
	Use:   "day10",
	Short: "Day 10 Solutions",
	Long:  "The solutions for day 10 of Advent of Code 2025",
	Run:   d10Main,
}

var d10p1Cmd = &cobra.Command{
	Use:   "part1",
	Short: "Day 10, 1 Solution",
	Long:  "The solution for day 10, part 1 of Advent of Code 2025",
	Run:   d10p1Main,
}

var d10p2Cmd = &cobra.Command{
	Use:   "part2",
	Short: "Day 10, 2 Solution",
	Long:  "The solution for day 10, part 2 of Advent of Code 2025",
	Run:   d10p2Main,
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show",
	Long:  "Show",
	Run:   showMain,
}

func d10Main(cmd *cobra.Command, args []string){
	_ = cmd.Help()
}

func d10p1Main(cmd *cobra.Command, args []string) {
	inputStr, err := loadInput(cmd, args)
	MaybeDie(err)

	topo, err := ParseTopoMap(inputStr)
	MaybeDie(err)

	ct := topo.CountTrails()

	slog.Debug("Results:", "TrailCount", ct)
	fmt.Printf("%v\n", ct)
}

func d10p2Main(cmd *cobra.Command, args []string) {
	inputStr, err := loadInput(cmd, args)
	MaybeDie(err)

	topo, err := ParseTopoMap(inputStr)
	MaybeDie(err)

	rt := topo.RateTrails()

	slog.Debug("Results:", "Rating", rt)
	fmt.Printf("%v\n", rt)
}

func showMain(cmd *cobra.Command, args []string) {
	inputStr, err := loadInput(cmd, args)
	MaybeDie(err)

	topo, err := ParseTopoMap(inputStr)
	MaybeDie(err)

	fmt.Printf("%v\n", topo)
}

func (tm *TopoMap) CountTrails() int {
	count := 0
	for th := range tm.THs.Iter() {
		for top := range tm.Tops.Iter() {
			if tm.DAG.HasPath(th, top) {
				count++
			}
		}
	}
	return count
}

func (tm *TopoMap) RateTrails() int {
	rating := 0
	for th := range tm.THs.Iter() {
		for top := range tm.Tops.Iter() {
			paths := tm.DAG.Paths(th, top)
			rating += len(paths)
		}
	}
	return rating
}

type TopoMap struct {
	Size  util.Size
	Elevs []int
	DAG   *graph.Graph[util.Point]
	THs   util.Set[util.Point]
	Tops  util.Set[util.Point]
}

func (tm *TopoMap) GetElev(pt util.Point) (int, error) {
	idx, err := tm.Size.Idx(pt)
	if err != nil {
		return -1, util.ReErr(err, "Invalid location %v", pt)
	}
	return tm.Elevs[idx], nil
}

func (tm *TopoMap) String() string {
	offSz, err := util.MakeSize(tm.Size.W * 2 - 1, tm.Size.H * 2 - 1)
	if err != nil {
		return "ERROR RENDERING MAP"
	}
	runes := util.MakeFill(' ', offSz.Area())

	for pt := range tm.Size.Iter() {
		elev, err := tm.GetElev(pt)
		if err != nil {
			return fmt.Sprintf("ERROR RENDERING POINT: Couldn't fetch elevation for %v", pt)
		}

		off := pt.Mul(2)
		idx, err := offSz.Idx(off)
		if err != nil {
			return fmt.Sprintf("ERROR RENDERING POINT: Couldn't get offset at %v", off)
		}

		rn, err := util.ItoR(elev)
		if err != nil {
			return fmt.Sprintf("ERROR RENDERING Rune: Couldn't get rune at %v", off)
		}
		runes[idx] = rn
	}

	for pair := range tm.DAG.Edges().Iter() {
		from := pair.Left
		to   := pair.Right

		v := to.Diff(from)
		var rn rune
		switch v {
		case util.Vector{Di:  1, Dj:  0}:
			rn = '↓'
		case util.Vector{Di: -1, Dj:  0}:
			rn = '↑'
		case util.Vector{Di:  0, Dj:  1}:
			rn = '→'
		case util.Vector{Di:  0, Dj: -1}:
			rn = '←'
		default:
			rn = '?'
		}
		off := from.Mul(2).Add(v)

		idx, err := offSz.Idx(off)
		if err != nil {
			return fmt.Sprintf("ERROR RENDERING POINT: Couldn't get offset at %v", off)
		}
		runes[idx] = rn
	}

	out := ""
	for i := range offSz.H {
		out += string(runes[i * offSz.W:(i + 1) * offSz.W]) + "\n"
	}
	return out
}

func ParseTopoMap(inputStr string) (tm *TopoMap, err error) {
	inputStr = strings.TrimSpace(inputStr)

	lines := strings.Split(inputStr, "\n")
	west := util.MakeVector(0, -1)
	north := util.MakeVector(-1, 0)
	tm = &TopoMap{
		DAG:  graph.MakeGraph[util.Point](true),
		Size: util.Size{W: 0, H: len(lines)},
		THs:  util.MakeSet[util.Point](),
		Tops: util.MakeSet[util.Point](),
	}

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if tm.Size.W == 0 {
			tm.Size.W = len(line)
			tm.Elevs = make([]int, tm.Size.W * tm.Size.H)
		} else if len(line) != tm.Size.W {
			return nil, fmt.Errorf("Line #%v has a different size: %v", i, len(line))
		}
		for j, rn := range line {
			if rn == '.' {
				continue
			}

			pt := util.MakePoint(i, j)
			elev, err := util.RtoI(rn)
			if err != nil {
				return nil, util.ReErr(err, "Found invalid elevation at %+v", pt)
			}
			if elev == 0 {
				tm.THs.Add(pt)
			} else if elev == 9 {
				tm.Tops.Add(pt)
			}

			tm.Elevs[pt.I * tm.Size.W + pt.J] = elev

			nbor := pt.Add(west)
			if nelev, err := tm.GetElev(nbor); err == nil {
				d := nelev - elev
				if d == -1 {
					tm.DAG.AddEdge(nbor, pt)
				} else if d == 1 {
					tm.DAG.AddEdge(pt, nbor)
				}
			}

			nbor = pt.Add(north)
			if nelev, err := tm.GetElev(nbor); err == nil {
				d := nelev - elev
				if d == -1 {
					tm.DAG.AddEdge(nbor, pt)
				} else if d == 1 {
					tm.DAG.AddEdge(pt, nbor)
				}
			}
		}
	}
	return tm, nil
}
