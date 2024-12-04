package cmd

import (
	"fmt"
	"log/slog"
	"sort"
	"strconv"
	"strings"

	"github.com/dusktreader/advent-of-code-2024/util"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(d1Cmd)
	d1Cmd.PersistentFlags().StringP("input-file", "i", "", "Get input from a file instead of stdin")
	d1Cmd.AddCommand(d1p1Cmd)
	d1Cmd.AddCommand(d1p2Cmd)
}

var d1Cmd = &cobra.Command{
	Use:   "day01",
	Short: "Day 1 Solutions",
	Long:  "The solutions for day 1 of Advent of Code 2024",
	Run:   d1Main,
}

var d1p1Cmd = &cobra.Command{
	Use: "part1",
	Short: "Day 1, 1 Solution",
	Long:  "The solution for day 1, part 1 of Advent of Code 2024",
	Run:   d1p1Main,
}

var d1p2Cmd = &cobra.Command{
	Use: "part2",
	Short: "Day 1, 2 Solution",
	Long:  "The solution for day 1, part 2 of Advent of Code 2024",
	Run:   d1p2Main,
}

func d1Main(cmd *cobra.Command, args []string){
	_ = cmd.Help()
}

func d1GetInput(cmd *cobra.Command, args []string) (left []int, right []int, err error) {
	inputStr, err := loadInput(cmd, args)
	if err != nil {
		return left, right, fmt.Errorf("Couldn't load input: %#v", err)
	}

	inputStr = strings.TrimSpace(inputStr)
	inputStr = strings.ReplaceAll(inputStr, "\n", " ")

	isLeft := true

	tokens := strings.Split(inputStr, " ")
	for _, token := range tokens {
		if token != "" {
			number, err := strconv.Atoi(token)
			if err != nil {
				return left, right, fmt.Errorf("Failed to convert a token: %#v", err)
			}
			if isLeft {
				left = append(left, number)
			} else {
				right = append(right, number)
			}
			isLeft = !isLeft
		}
	}
	return left, right, nil
}

func d1p1Main(cmd *cobra.Command, args []string){
	left, right, err := d1GetInput(cmd, args)
	MaybeDie(err)

	distance, err := TotalDistance(left, right)
	MaybeDie(err)

	slog.Debug("Results:", "distance", distance, "left", left, "right", right)
	fmt.Printf("%v\n", distance)
}

func d1p2Main(cmd *cobra.Command, args []string){
	left, right, err := d1GetInput(cmd, args)
	MaybeDie(err)

	similarity, err := Similarity(left, right)
	MaybeDie(err)

	slog.Debug("Results:", "similarity", similarity, "left", left, "right", right)
	fmt.Printf("%v\n", similarity)
}

func TotalDistance(left []int, right []int) (int, error) {
	if len(left) != len(right) {
		return 0, fmt.Errorf("Need an even number of number inputs")
	}

	sort.Ints(left)
	sort.Ints(right)

	var total int

	for i := 0; i < len(left); i++ {
		total += util.AbsInt(left[i] - right[i])
	}
	return total, nil
}

func Similarity(left []int, right []int) (int, error) {
	if len(left) != len(right) {
		return 0, fmt.Errorf("Need an even number of number inputs")
	}

	leftMap := make(map[int]int)
	rightMap := make(map[int]int)
	allKeys := make(map[int]bool)

	for i := 0; i < len(left); i++ {
		leftMap[left[i]]++
		allKeys[left[i]] = true
		rightMap[right[i]]++
		allKeys[right[i]] = true
	}

	var total int
	for key, _ := range allKeys {
		total += key * leftMap[key] * rightMap[key]
	}

	return total, nil
}
