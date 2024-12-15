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
	rootCmd.AddCommand(d11Cmd)
	d11Cmd.PersistentFlags().StringP("input-file", "i", "", "Get input from a file instead of stdin")
	d11Cmd.AddCommand(d11p1Cmd)
	d11Cmd.AddCommand(d11p2Cmd)
	d11Cmd.AddCommand(blinkCmd)
	blinkCmd.Flags().IntP("count", "n", 75, "Blink n times")
}

var d11Cmd = &cobra.Command{
	Use:   "day11",
	Short: "Day 11 Solutions",
	Long:  "The solutions for day 11 of Advent of Code 2025",
	Run:   d11Main,
}

var d11p1Cmd = &cobra.Command{
	Use:   "part1",
	Short: "Day 11, 1 Solution",
	Long:  "The solution for day 11, part 1 of Advent of Code 2025",
	Run:   d11p1Main,
}

var d11p2Cmd = &cobra.Command{
	Use:   "part2",
	Short: "Day 11, 2 Solution",
	Long:  "The solution for day 11, part 2 of Advent of Code 2025",
	Run:   d11p2Main,
}

var blinkCmd = &cobra.Command{
	Use:   "blink",
	Short: "Blink n times",
	Long:  "Blink a custom number of times",
	Run:   blinkMain,
}

func d11Main(cmd *cobra.Command, args []string){
	_ = cmd.Help()
}

func d11p1Main(cmd *cobra.Command, args []string) {
	inputStr, err := loadInput(cmd, args)
	MaybeDie(err)

	stones, err := ParseStones(inputStr)
	MaybeDie(err)

	ct := CountStones(stones, 25)

	slog.Debug("Results:", "StoneCount", ct)
	fmt.Printf("%v\n", ct)
}

func d11p2Main(cmd *cobra.Command, args []string) {
	inputStr, err := loadInput(cmd, args)
	MaybeDie(err)

	stones, err := ParseStones(inputStr)
	MaybeDie(err)

	ct := CountStones(stones, 75)

	slog.Debug("Results:", "StoneCount", ct)
	fmt.Printf("%v\n", ct)
}

func blinkMain(cmd *cobra.Command, args []string) {
	blinks, err := cmd.Flags().GetInt("count")
	MaybeDie(err)

	inputStr, err := loadInput(cmd, args)
	MaybeDie(err)

	stones, err := ParseStones(inputStr)
	MaybeDie(err)

	ct := CountStones(stones, blinks)

	slog.Debug("Results:", "StoneCount", ct)
	fmt.Printf("%v\n", ct)
}

func Split(v int) (int, int, bool) {
	ct    := util.DigiCt(v)
	if ct % 2 != 0 {
		return -1, -1, false
	}
	div   := int(math.Pow10(ct / 2))
	left  := v / div
	right := v % div
	return left, right, true
}

func Blinker(stone int, blinks int, cache map[util.Pair[int]]int) int {
	var count int

	if blinks == 0 {
		return 1
	}

	p := util.MakePair(stone, blinks)
	prev, ok := cache[p]
	if ok {
		return prev
	}

	if stone == 0 {
		count = Blinker(1, blinks - 1, cache)
		cache[p] = count
		return count
	}

	left, right, ok := Split(stone)
	if ok {
		count = Blinker(left, blinks - 1, cache) + Blinker(right, blinks - 1, cache)
		cache[p] = count
		return count
	}

	count = Blinker(stone * 2024, blinks - 1, cache)
	cache[p] = count
	return count
}

func CountStones(stones []int, blinks int) int {
	cache := make(map[util.Pair[int]]int)
	count := 0
	for _, stone := range stones {
		count += Blinker(stone, blinks, cache)
	}
	return count
}

func ParseStones(inputStr string) ([]int, error) {
	inputStr = strings.TrimSpace(inputStr)
	tokens := strings.Split(inputStr, " ")
	stones := make([]int, len(tokens))
	for i, t := range tokens {
		s, err := strconv.Atoi(t)
		if err != nil {
			return []int{}, fmt.Errorf("Couldn't convert stone value at index %v", i)
		}
		stones[i] = s
	}
	return stones, nil
}
