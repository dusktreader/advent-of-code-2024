package cmd

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/dusktreader/advent-of-code-2024/util"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(d2Cmd)
	d2Cmd.PersistentFlags().StringP("input-file", "i", "", "Get input from a file instead of stdin")
	d2Cmd.AddCommand(d2p1Cmd)
	d2Cmd.AddCommand(d2p2Cmd)
}

var d2Cmd = &cobra.Command{
	Use:   "day02",
	Short: "Day 2 Solutions",
	Long:  "The solutions for day 2 of Advent of Code 2024",
	Run:   d2Main,
}

var d2p1Cmd = &cobra.Command{
	Use:   "part1",
	Short: "Day 2, 1 Solution",
	Long:  "The solution for day 2, part 1 of Advent of Code 2024",
	Run:   d2p1Main,
}

var d2p2Cmd = &cobra.Command{
	Use:   "part2",
	Short: "Day 2, 2 Solution",
	Long:  "The solution for day 2, part 2 of Advent of Code 2024",
	Run:   d2p2Main,
}

func d2Main(cmd *cobra.Command, args []string){
	_ = cmd.Help()
}

func d2p1Main(cmd *cobra.Command, args []string) {
	inputStr, err := loadInput(cmd, args)
	MaybeDie(err)

	reports, err := ParseReport(inputStr)
	MaybeDie(err)

	safeCount := CountSafe(reports, 3, 0)

	slog.Debug("Results (undampened):", "safeCount", safeCount)
	fmt.Printf("%v\n", safeCount)
}

func d2p2Main(cmd *cobra.Command, args []string) {
	inputStr, err := loadInput(cmd, args)
	MaybeDie(err)

	reports, err := ParseReport(inputStr)
	MaybeDie(err)

	safeCount := CountSafe(reports, 3, 1)

	slog.Debug("Results (dampened):", "safeCount", safeCount)
	fmt.Printf("%v\n", safeCount)
}

func ParseReport(inputStr string) (reports [][]int, err error) {
	lines := strings.Split(inputStr, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			slog.Debug("Skipping empty line")
			continue
		}

		tokens := strings.Split(line, " ")
		var numbers []int
		for _, token := range tokens {
			if token != "" {
				number, err := strconv.Atoi(token)
				if err != nil {
					return reports, fmt.Errorf("Failed to convert a token: %#v", err)
				}
				numbers = append(numbers, number)
			}
		}
		reports = append(reports, make([]int, len(numbers)))
		copy(reports[len(reports) - 1], numbers)
	}
	return
}

type JLimit struct {
	start int
	stop  int
	incr  int
}

func IsSafe(report []int, diffMax int, dampMax int) (safe bool) {
	slog.Debug("-----------------------------")
	slog.Debug("Checking report:", "report", report, "diffMax", diffMax, "dampMax", dampMax)

	var dir int
	var diff int
	var absDiff int
	var dampCount int
	var lastValue int

	l := len(report)
	if l < 2 + dampMax {
		return true
	}

	limits := []JLimit{
		{start:   1, stop:  l, incr:  1},
		{start: l-2, stop: -1, incr: -1},
	}

	for _, lim := range limits {

		if lim.incr == 1 {
			slog.Debug("Sweeping forward--->")
		} else {
			slog.Debug("<---Sweeping backward")
		}

		dir = 0
		dampCount = 0
		lastValue = report[lim.start - lim.incr]
		safe = true
		for j := lim.start; j != lim.stop; j += lim.incr {
			diff = report[j] - lastValue
			absDiff = util.AbsInt(diff)
			if dir == 0 && absDiff > 0 {
				dir = diff / absDiff
				slog.Debug("Direction is:", "dir", dir)
			}
			slog.Debug("Checking index:", "j", j, "diff", diff)

			if absDiff < 1 || absDiff > diffMax {
				if dampCount < dampMax {
					slog.Debug("Dampened (diff) at:", "j", j)
					dampCount++
					continue
				} else {
					slog.Debug("Unsafe (diff) at:", "j", j)
					safe = false
					break
				}
			} else if dir != diff / absDiff {
				if dampCount < dampMax {
					slog.Debug("Dampened (dir) at:", "j", j)
					dampCount++
					continue
				} else {
					slog.Debug("Unsafe (dir) at:", "j", j)
					safe = false
					break
				}
			}
			lastValue = report[j]
		}
		if safe {
			slog.Debug("Report was safe")
			return
		}
	}
	slog.Debug("Report was unsafe")
	return false
}

func CountSafe(reports [][]int, diffMax int, dampMax int) (safeCount int) {
	if len(reports) == 0 {
		return 0
	}

	slog.Debug("Counting safe:", "diffMax", diffMax, "dampMax", dampMax)
	for _, report := range reports {
		if IsSafe(report, diffMax, dampMax) {
			safeCount++
		}
	}
	return
}
