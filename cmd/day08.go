package cmd

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/dusktreader/advent-of-code-2024/util"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(d8Cmd)
	d8Cmd.PersistentFlags().StringP("input-file", "i", "", "Get input from a file instead of stdin")
	d8Cmd.PersistentFlags().StringP("output-file", "o", "", "Draw the final map and store it in the file")
	d8Cmd.AddCommand(d8p1Cmd)
	d8Cmd.AddCommand(d8p2Cmd)
}

var d8Cmd = &cobra.Command{
	Use:   "day08",
	Short: "Day 8 Solutions",
	Long:  "The solutions for day 8 of Advent of Code 2025",
	Run:   d8Main,
}

var d8p1Cmd = &cobra.Command{
	Use:   "part1",
	Short: "Day 8, 1 Solution",
	Long:  "The solution for day 8, part 1 of Advent of Code 2025",
	Run:   d8p1Main,
}

var d8p2Cmd = &cobra.Command{
	Use:   "part2",
	Short: "Day 8, 2 Solution",
	Long:  "The solution for day 8, part 2 of Advent of Code 2025",
	Run:   d8p2Main,
}

func d8Main(cmd *cobra.Command, args []string){
	_ = cmd.Help()
}

func d8p1Main(cmd *cobra.Command, args []string) {
	inputStr, err := loadInput(cmd, args)
	MaybeDie(err)

	am, err := ParseAntMap(inputStr)
	MaybeDie(err)

	am.FindAll(false)

	ct := am.CountAns()

	slog.Debug("Results:", "AntinodeCount", ct)
	fmt.Printf("%v\n", ct)
}

func d8p2Main(cmd *cobra.Command, args []string) {
	inputStr, err := loadInput(cmd, args)
	MaybeDie(err)

	am, err := ParseAntMap(inputStr)
	MaybeDie(err)

	am.FindAll(true)

	ct := am.CountAns()

	slog.Debug("Results:", "AntinodeCount", ct)
	fmt.Printf("%v\n", ct)
}

type AntMap struct {
	Size util.Size
	Ants util.SetMap[rune, util.Point]
	Ans  util.SetMap[rune, util.Point]
}

func MakeAntMap(w int, h int) (*AntMap, error) {
	var err error
	am := AntMap{}
	am.Size, err = util.MakeSize(w, h)
	if err != nil {
		return nil, util.ReErr(err, "Couldn't make AntMap")
	}
	am.Ants = util.MakeSetMap[rune, util.Point]()
	return &am, nil
}

func (am *AntMap) Eq(om *AntMap) bool {
	if am.Size != om.Size {
		slog.Debug("Sizes didn't match", "us", am.Size, "them", om.Size)
		return false
	} else if !am.Ants.Eq(&om.Ants) {
		slog.Debug("Antenna didn't match", "us", am.Ants, "them", om.Ants)
		return false
	}
	return true
}

func FindAnsSimp(sz util.Size, pts util.Set[util.Point]) util.Set[util.Point] {
	ans := util.MakeSet[util.Point]()
	for a := range pts.Iter() {
		for b := range pts.Iter() {
			if a == b {
				continue
			}
			d := a.Diff(b)

			t := a.Add(d)
			if !sz.Out(t) {
				ans.Add(t)
			}

			t = b.Add(d.Neg())
			if !sz.Out(t) {
				ans.Add(t)
			}
		}
	}
	return ans
}

func FindAnsRes(sz util.Size, pts util.Set[util.Point]) util.Set[util.Point] {
	anodes := util.MakeSet[util.Point]()
	for a := range pts.Iter() {
		for b := range pts.Iter() {
			if a == b {
				continue
			}

			d := a.Diff(b)
			d = d.Div(util.GCD(d.Di, d.Dj))
			anodes.Add(a)
			i := 1
			for {
				t := a.Add(d.Mul(i))
				if sz.Out(t) {
					break
				}
				anodes.Add(t)
				i++
			}

			d = d.Neg()
			i = 1
			for {
				t := a.Add(d.Mul(i))
				if sz.Out(t) {
					break
				}
				anodes.Add(t)
				i++
			}
		}
	}
	return anodes
}

func (am *AntMap) FindAll(res bool) {
	am.Ans.Clear()
	for rn, pts := range am.Ants.Iter() {
		var f func(util.Size, util.Set[util.Point]) util.Set[util.Point]

		if res {
			f = FindAnsRes
		} else {
			f = FindAnsSimp
		}

		am.Ans.Add(rn, f(am.Size, pts).Items()...)
	}
}

func (am *AntMap) CountAns() int {
	all := util.MakeSet[util.Point]()
	for _, pts := range am.Ans.Iter() {
		all.Add(pts.Items()...)
	}
	return all.Size()
}


func ParseAntMap(inputStr string) (am *AntMap, err error) {
	inputStr = strings.TrimSpace(inputStr)

	lines := strings.Split(inputStr, "\n")

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if am == nil {
			am, err = MakeAntMap(len(lines), len(line))
			if err != nil {
				return nil, util.ReErr(err, "Couldn't parse antenna map")
			}
		} else if len(line) != am.Size.W {
			return nil, fmt.Errorf("Line #%v has a different size: %v", i, len(line))
		}
		for j, rn := range line {
			if rn == '.' {
				continue
			}
			pt := util.MakePoint(i, j)
			if err != nil {
				return nil, util.ReErr(err, "Bad point %v", pt)
			}
			am.Ants.Add(rn, pt)
		}
	}
	return am, nil
}

func (am *AntMap) Clone() *AntMap {
	om := AntMap{}
	om.Size = am.Size
	om.Ants = am.Ants.Clone()
	return &om
}

func (am *AntMap) String() string {
	l := am.Size.W * am.Size.H
	grid := make([]rune, l)
	for i := range l {
		grid[i] = '.'
	}

	for rn, as := range am.Ants.Iter() {
		for pt := range as.Iter() {
			grid[pt.I * am.Size.W + pt.J] = rn
		}
	}

	lines := make([]string, am.Size.H)
	for i := range am.Size.H {
		s := i * am.Size.W
		e := s + am.Size.W
		lines[i] = string(grid[s:e])
	}
	return strings.Join(lines, "\n")
}
