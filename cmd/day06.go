package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/dusktreader/advent-of-code-2024/util"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(d6Cmd)
	d6Cmd.PersistentFlags().StringP("input-file", "i", "", "Get input from a file instead of stdin")
	d6Cmd.PersistentFlags().StringP("output-file", "o", "", "Draw the final map and store it in the file")
	d6Cmd.AddCommand(d6p1Cmd)
	d6Cmd.AddCommand(d6p2Cmd)
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

	lm, err := ParseLabMap(inputStr)
	MaybeDie(err)

	err = lm.Patrol()
	MaybeDie(err)

	outputFile, err := cmd.Flags().GetString("output-file")
	if err == nil && outputFile != "" {
		err = os.WriteFile(outputFile, []byte(lm.String()), 0644)
		if err != nil {
			slog.Error("Couldn't write to output-file:", "file", outputFile, "err", err)
		}
	}

	ct := lm.CountVisits()

	slog.Debug("Results:", "VisitCount", ct)
	fmt.Printf("%v\n", ct)
}

func d6p2Main(cmd *cobra.Command, args []string) {
	inputStr, err := loadInput(cmd, args)
	MaybeDie(err)

	lm, err := ParseLabMap(inputStr)
	MaybeDie(err)

	err = lm.Loopify()
	MaybeDie(err)

	outputFile, err := cmd.Flags().GetString("output-file")
	if err == nil && outputFile != "" {
		lm.ClearVisits()
		slog.Debug("Dumping output as requested:", "file", outputFile)
		err = os.WriteFile(outputFile, []byte(lm.String()), 0644)
		if err != nil {
			slog.Error("Couldn't write to output-file:", "file", outputFile, "err", err)
		}
	}

	ct := lm.CountLoopers()

	slog.Debug("Results:", "NewObstacleCount", ct)
	fmt.Printf("%v\n", ct)
}

type LabMap struct {
	Size     util.Size
	GuardPos util.Point
	GuardDir util.Vector
	Visits   util.SetMap[util.Point, util.Vector]
	Obstr    util.Set[util.Point]
	Loopers  util.Set[util.Point]
}

func (lm *LabMap) Eq(om *LabMap) bool {
	if lm.Size != om.Size {
		slog.Debug("Sizes didn't match", "us", lm.Size, "them", om.Size)
		return false
	} else if lm.GuardPos != om.GuardPos {
		slog.Debug("Guard positions didn't match", "us", lm.GuardPos, "them", om.GuardPos)
		return false
	} else if lm.GuardDir != om.GuardDir {
		slog.Debug("Guard directions didn't match", "us", lm.GuardDir, "them", om.GuardDir)
		return false
	} else if !lm.Obstr.Eq(om.Obstr) {
		slog.Debug("Obstructions didn't match", "us", lm.Obstr, "them", om.Obstr)
		return false
	} else if !lm.Loopers.Eq(om.Loopers) {
		slog.Debug("Looping obstuctions didn't match", "us", lm.Loopers, "them", om.Loopers)
		return false
	} else if !lm.Visits.Eq(&om.Visits) {
		slog.Debug("Visits didn't match", "us", lm.Visits, "them", om.Visits)
	}
	return true
}

func (lm *LabMap) ClearVisits() {
	lm.Visits.Clear()
}

func (lm *LabMap) Visit(pt util.Point, dir util.Vector) error {
	visits := lm.Visits.Get(pt)
	if !visits.Has(dir) {
		visits.Add(dir)
	} else {
		return fmt.Errorf("Walk cycle detected at %v facing %v", pt, dir)
	}
	lm.GuardPos = pt
	lm.GuardDir = dir
	return nil
}

func MakeLabMap(w int, h int) (*LabMap, error) {
	var err error
	lm := LabMap{}
	lm.Size, err = util.MakeSize(w, h)
	if err != nil {
		return nil, util.ReErr(err, "Couldn't make LabMap")
	}
	lm.Visits  = util.MakeSetMap[util.Point, util.Vector]()
	lm.Obstr   = util.MakeSet[util.Point]()
	lm.Loopers = util.MakeSet[util.Point]()
	return &lm, nil
}

func (lm *LabMap) Clone() *LabMap {
	om := LabMap{}
	om.Size     = lm.Size
	om.GuardPos = lm.GuardPos
	om.GuardDir = lm.GuardDir
	om.Visits   = lm.Visits.Clone()
	om.Obstr    = lm.Obstr.Clone()
	om.Loopers  = lm.Loopers.Clone()
	return &om
}

func (lm *LabMap) Walk() (bool, error) {
	lastPos := lm.GuardPos
	lm.GuardPos = lm.GuardPos.Add(lm.GuardDir)
	if lm.Size.Out(lm.GuardPos) {
		return false, nil
	}

	if lm.Obstr.Has(lm.GuardPos) {
		lm.GuardPos = lastPos
		lm.GuardDir = lm.GuardDir.RotCW()
	}
	err := lm.Visit(lm.GuardPos, lm.GuardDir)
	if err != nil {
		return false, fmt.Errorf("Couldn't walk: %#v", err)
	}
	return true, nil
}

func (lm *LabMap) Patrol() error {
	for {
		stillIn, err := lm.Walk()
		if err != nil {
			return util.ReErr(err, "Patrol failed")
		} else if !stillIn {
			return nil
		}
	}
}

func (lm *LabMap) CountVisits() int {
	return lm.Visits.Size()
}

func (lm *LabMap) Loopify() error {
	for {
		oPos := lm.GuardPos.Add(lm.GuardDir)
		if !lm.Size.Out(oPos) && !lm.Obstr.Has(oPos) && !lm.Visits.Has(oPos) {
			om := lm.Clone()
			om.Obstr.Add(oPos)
			err := om.Patrol()
			if err != nil {
				lm.Loopers.Add(oPos)
			}
		}

		stillIn, err := lm.Walk()
		if err != nil {
			return util.ReErr(err, "Map already had an infinite loop")
		} else if !stillIn {
			return nil
		}
	}
}

func (lm *LabMap) CountLoopers() int {
	return lm.Loopers.Size()
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
				lm.Obstr.Add(pt)
			case '^':
				lm.Visit(pt, util.MakeVector(-1, 0))
			case '>':
				lm.Visit(pt, util.MakeVector(0, 1))
			case 'v':
				lm.Visit(pt, util.MakeVector(1, 0))
			case '<':
				lm.Visit(pt, util.MakeVector(0, -1))
			}
		}
	}
	return lm, nil
}

func (lm *LabMap) String() string {
	lines := make([]string, lm.Size.H)
	for i := range lm.Size.H {
		line := make([]rune, lm.Size.W)
		for j := range lm.Size.W {
			pt := util.MakePoint(i, j)
			if lm.Obstr.Has(pt) {
				line[j] = '#'
			} else if pt == lm.GuardPos {
				line[j] = lm.GuardDir.Pretty()
			} else if lm.Loopers.Has(pt) {
				line[j] = 'O'
			} else if lm.Visits.Has(pt) {
				line[j] = 'X'
			} else {
				line[j] = '.'
			}

		}
		lines[i] = string(line)
	}
	return strings.Join(lines, "\n")
}
