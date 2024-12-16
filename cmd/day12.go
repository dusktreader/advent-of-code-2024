package cmd

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/dusktreader/advent-of-code-2024/util"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(d12Cmd)
	d12Cmd.PersistentFlags().StringP("input-file", "i", "", "Get input from a file instead of stdin")
	d12Cmd.AddCommand(d12p1Cmd)
	d12Cmd.AddCommand(d12p2Cmd)
}

var d12Cmd = &cobra.Command{
	Use:   "day12",
	Short: "Day 12 Solutions",
	Long:  "The solutions for day 12 of Advent of Code 2025",
	Run:   d12Main,
}

var d12p1Cmd = &cobra.Command{
	Use:   "part1",
	Short: "Day 12, 1 Solution",
	Long:  "The solution for day 12, part 1 of Advent of Code 2025",
	Run:   d12p1Main,
}

var d12p2Cmd = &cobra.Command{
	Use:   "part2",
	Short: "Day 12, 2 Solution",
	Long:  "The solution for day 12, part 2 of Advent of Code 2025",
	Run:   d12p2Main,
}

func d12Main(cmd *cobra.Command, args []string){
	_ = cmd.Help()
}

func d12p1Main(cmd *cobra.Command, args []string) {
	inputStr, err := loadInput(cmd, args)
	MaybeDie(err)

	slog.Debug("Parsing gardens")
	garden, err := ParseGarden(inputStr)
	MaybeDie(err)

	slog.Debug("Finding regions")
	garden.FindRegions()

	slog.Debug("Computing price")
	price := garden.Price(false)

	slog.Debug("Results:", "GardenPrice", price)
	fmt.Printf("%v\n", price)
}

func d12p2Main(cmd *cobra.Command, args []string) {
	inputStr, err := loadInput(cmd, args)
	MaybeDie(err)

	slog.Debug("Parsing gardens")
	garden, err := ParseGarden(inputStr)
	MaybeDie(err)

	slog.Debug("Finding regions")
	garden.FindRegions()

	slog.Debug("Computing price")
	price := garden.Price(true)

	slog.Debug("Results:", "GardenPrice", price)
	fmt.Printf("%v\n", price)
}

type Garden struct {
	Size util.Size
	Plots []Plot
	Regions []util.Set[util.Point]
}

type Plot struct {
	Loc    util.Point
	Label  rune
}

func ParseGarden(inputStr string) (*Garden, error) {
	inputStr = strings.TrimSpace(inputStr)
	lines := strings.Split(inputStr, "\n")
	gardens := &Garden{
		Size: util.Size{W: 0, H: len(lines)},
		Regions: []util.Set[util.Point]{},
	}

	for i, line := range lines {
		if gardens.Size.W == 0 {
			gardens.Size.W = len(line)
			gardens.Plots = make([]Plot, gardens.Size.Area())
		} else if gardens.Size.W != len(line) {
			return nil, fmt.Errorf("Line length didn't match at index %v", i)
		}

		for j, rn := range line {
			pt := util.MakePoint(i, j)
			idx, err := gardens.Size.Idx(pt)
			if err != nil {
				return nil, util.ReErr(err, "Couldn't compute an index")
			}
			gardens.Plots[idx] = Plot{
				Loc: pt,
				Label: rn,
			}
		}
	}
	return gardens, nil
}

var rose = [4]util.Vector{
	{Di: -1, Dj:  0},
	{Di:  0, Dj:  1},
	{Di:  1, Dj:  0},
	{Di:  0, Dj: -1},
}

func (grd *Garden) FindRegions() {
	visited := util.MakeSet[Plot]()
	for _, plot := range grd.Plots {
		if visited.Has(plot) {
			continue
		}

		nbors := util.MakeSet(plot)
		region := util.MakeSet[util.Point]()

		for !nbors.Empty() {
			p := nbors.Pop()
			visited.Add(p)
			region.Add(p.Loc)

			for _, dir := range rose {
				idx, err := grd.Size.Idx(p.Loc.Add(dir))
				if err != nil {
					continue
				}
				o := grd.Plots[idx]
				if p.Label == o.Label {
					if !visited.Has(o) {
						nbors.Add(o)
					}
				}
			}
		}
		slog.Debug("Found region", "r", region)
		grd.Regions = append(grd.Regions, region)
	}
}

func PriceRegion(region util.Set[util.Point], discount bool) (price int) {
	area := region.Size()
	fences := 0

	for pt := range region.Iter() {

		for _, dir := range rose {
			if !region.Has(pt.Add(dir)) {
				if !discount {
					fences++
				} else {
					left := pt.Add(dir.RotCCW())
					diag := left.Add(dir)
					if !region.Has(left) || region.Has(diag) {
						fences++
					}

				}
			}
		}
	}

	return fences * area
}

func (grd *Garden) Price(discount bool) (price int) {
	for _, region := range grd.Regions {
		price += PriceRegion(region, discount)
	}
	return
}
