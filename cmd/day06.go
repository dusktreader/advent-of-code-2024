package cmd

import (
	"fmt"
"log/slog"
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

type Visit struct {
	Pt   util.Point
	Dir  util.Vector
}

type LabMap struct {
	Size     util.Size
	GuardPos util.Point
	GuardDir util.Vector
	Visits   util.Set[Visit]
	Obstr    util.Set[util.Point]

}

func MakeLabMap(w int, h int) (*LabMap, error) {
	var err error
	lm := LabMap{}
	lm.Size, err = util.MakeSize(w, h)
	if err != nil {
		return nil, util.ReErr(err, "Couldn't make LabMap")
	}
	lm.Visits = util.MakeSet[Visit]()
	lm.Obstr  = util.MakeSet[util.Point]()
	return &lm, nil
}

func (lm *LabMap) Walk() (bool, error) {
	var newPos util.Point
	var newDir util.Vector

	newPos = lm.GuardPos.Move(lm.GuardDir)
	cell, err := lm.Grid.Get(newPos)
	if err != nil {
		// Error indicates we've exited the grid
		return false, nil
	}

	if cell.HasObs {
		newDir = lm.GuardDir.Rot()
		newPos = lm.GuardPos
	} else {
		newDir = lm.GuardDir
	}

	if cell.PrevDirs.Has(lm.GuardDir) {
		return false, fmt.Errorf("Walk cycle detected at %v facing %v", newPos, newDir)
	}

	lm.GuardDir = newDir
	lm.GuardPos = newPos
	cell.PrevDirs.Add(newDir)
	return true, nil
}


func ParseLabMap(inputStr string) (lm *LabMap, err error) {
	inputStr = strings.TrimSpace(inputStr)

	lines := strings.Split(inputStr, "\n")

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if lm == nil {
			lm, err = MakeLabMap(len(lines), len(line))
			if err != nil {
				return nil, util.ReErr(err, "Couldn't parse lab map")
			}
		} else if len(line) != lm.Size.W {
			return nil, fmt.Errorf("Line #%v has a different size: %v", i, len(line))
		}
		for j, rn := range line {
			pt := util.MakePoint(i, j)
			if err != nil {
				return nil, util.ReErr(err, "Bad point %v", pt)
			}
			switch rn {
			case '#':
				v.HasObs = true
			case '^':
				lm.GuardPos = pt
				lm.GuardDir = util.MakeVector(-1, 0)
				lm.Visits.Add(lm.GuardPos)
				v.PrevDirs.Add(lm.GuardDir)
			case '>':
				lm.GuardPos = pt
				lm.GuardDir = util.MakeVector(0, 1)
				lm.Visits.Add(lm.GuardPos)
				v.PrevDirs.Add(lm.GuardDir)
			case 'v':
				lm.GuardPos = pt
				lm.GuardDir = util.MakeVector(1, 0)
				lm.Visits.Add(lm.GuardPos)
				v.PrevDirs.Add(lm.GuardDir)
			case '<':
				lm.GuardPos = pt
				lm.GuardDir = util.MakeVector(0, -1)
				lm.Visits.Add(lm.GuardPos)
				v.PrevDirs.Add(lm.GuardDir)
			}
		}
	}
	return &lm, err
}

func (lm *LabMap) String() string {
	lines := make([]string, lm.Grid.Size().H)
	for i := range lm.Grid.Size().H {
		line := make([]rune, lm.Grid.Size().W)
		for j := range lm.Grid.Size().W {
			pt := util.MakePoint(i, j)
			c, err := lm.Grid.Get(pt)
			if err != nil {
				return "INVALID MAP"
			}
			if c.HasObs {
				line[j] = '#'
			} else if pt == lm.GuardPos {
				line[j] = lm.GuardDir.Pretty()
			} else if c.TreadCt > 0 {
				line[j] = 'X'
			} else {
				line[j] = '.'
			}

		}
		lines[i] = string(line)
	}
	return strings.Join(lines, "\n")
}
