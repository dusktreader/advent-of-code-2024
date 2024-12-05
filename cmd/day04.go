package cmd

import (
	"fmt"
	"log/slog"
	"math/rand/v2"
	"slices"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(d4Cmd)
	d4Cmd.PersistentFlags().StringP("input-file", "i", "", "Get input from a file instead of stdin")
	d4Cmd.AddCommand(d4p1Cmd)
	d4Cmd.AddCommand(d4p2Cmd)
	d4p2Cmd.Flags().BoolP("annotate", "a", false, "Annotate redacted input string")
}

var d4Cmd = &cobra.Command{
	Use:   "day04",
	Short: "Day 4 Solutions",
	Long:  "The solutions for day 4 of Advent of Code 2024",
	Run:   d4Main,
}

var d4p1Cmd = &cobra.Command{
	Use:   "part1",
	Short: "Day 3, 1 Solution",
	Long:  "The solution for day 4, part 1 of Advent of Code 2024",
	Run:   d4p1Main,
}

var d4p2Cmd = &cobra.Command{
	Use:   "part2",
	Short: "Day 3, 2 Solution",
	Long:  "The solution for day 4, part 2 of Advent of Code 2024",
	Run:   d4p2Main,
}

func d4Main(cmd *cobra.Command, args []string){
	_ = cmd.Help()
}

func d4p1Main(cmd *cobra.Command, args []string) {
	inputStr, err := loadInput(cmd, args)
	MaybeDie(err)

	board, err := Runify(inputStr)
	MaybeDie(err)

	count := CountMatches([]rune("XMAS"), board)

	slog.Debug("Results:", "count", count)
	fmt.Printf("%v\n", count)
}

func d4p2Main(cmd *cobra.Command, args []string) {
	inputStr, err := loadInput(cmd, args)
	MaybeDie(err)

	board, err := Runify(inputStr)
	MaybeDie(err)

	count := CountCrossWords([]rune("MAS"), board)

	slog.Debug("Results:", "count", count)
	fmt.Printf("%v\n", count)
}

func Runify(inputStr string) (runes [][]rune, err error) {
	inputStr = strings.TrimSpace(inputStr)
	lines := strings.Split(inputStr, "\n")
	var lineLen int
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			return nil, fmt.Errorf("Line #%v was empty", i)
		} else if lineLen == 0 {
			lineLen = len(line)
		} else if lineLen != len(line) {
			return nil, fmt.Errorf("Line #%v didn't match previous line lengths", i)
		}
		runes = append(runes, []rune(line))
	}
	return
}

type Point struct {
	I int
	J int
}

var vectors = []Point{
	{ 0,  1}, // right
	{ 1,  1}, // down-right
	{ 1,  0}, // down
	{ 1, -1}, // down-1eft
	{ 0, -1}, // 1eft
	{-1, -1}, // up-1eft
	{-1,  0}, // up
	{-1,  1}, // up-right
}

func add(a Point, b Point) (Point) {
	return Point{a.I + b.I, a.J + b.J}
}

func mul(a Point, m int) (Point) {
	return Point{a.I * m, a.J * m}
}

func inBounds(board [][]rune, pt Point) (bool) {
	return pt.I >= 0 && pt.J >= 0 && pt.I < len(board) && pt.J < len(board[0])
}

func fetch(board [][]rune, pt Point) (rune) {
	return board[pt.I][pt.J]
}

func CountMatches(word []rune, board [][]rune) (count int) {
	l := len(word)
	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[0]); j++ {
			pt := Point{i, j}
			slog.Debug("-------------")
			slog.Debug("Checking at point:", "i", pt.I, "j", pt.J)
			for _, v := range vectors {
				slog.Debug("Checking with vector:", "di", v.I, "dj", v.J)
				for m := l - 1; m >= 0; m-- {
					off := add(pt, mul(v, m))
					slog.Debug("Checking at offset:", "i", off.I, "j", off.J)
					if m == l - 1 && !inBounds(board, off) {
						slog.Debug("Out of bounds!")
						break
					} else if fetch(board, off) != word[m] {
						slog.Debug("Doesn't match!")
						break
					} else if m == 0 {
						slog.Debug("Count it!")
						count++
					}
				}
			}
			slog.Debug("-------------")
		}
	}
	return
}

func MakePatch(l int) (newPatch [][]rune) {
	newPatch = make([][]rune, l)
	for i := 0; i < l; i++ {
		newPatch[i] = slices.Repeat([]rune{'.'}, l)
	}
	return
}

func RotatePatch(patch [][]rune) ([][]rune) {
	// Assume square
	l := len(patch)
	newPatch := MakePatch(l)
	for j := 0; j < l; j++ {
		for i := l - 1; i >= 0; i-- {
			newPatch[j][l - 1 - i] = patch[i][j]
		}
	}
	return newPatch
}

func PrettyPrintPatch(patch [][]rune) (pretty string) {
	for i := 0; i < len(patch); i++ {
		pretty += string(patch[i]) + "\n"
	}
	return
}

func MakeCrossPatch(word []rune) ([][]rune) {
	l := len(word)
	patch := MakePatch(l)
	for i := 0; i < l; i++ {
		patch[i][i] = word[i]
		patch[i][l - i - 1] = word[i]
	}
	return patch
}

func PatchMatch(patch [][]rune, board [][]rune, pt Point) (bool) {
	// assume pt is within board and buffered from edge
	l := len(patch)
	r := l / 2
	for i := 0; i < l; i++ {
		for j := 0; j < l; j++ {
			if patch[i][j] == '.' {
				continue
			}

			bi := pt.I - r + i
			bj := pt.J - r + j
			if patch[i][j] != board[bi][bj] {
				return false
			}
		}
	}
	return true
}

func PatchRandFill(patch [][]rune) {
	l := len(patch)
	for i := 0; i < l; i++ {
		for j := 0; j < l; j++ {
			if patch[i][j] == '.' {
				patch[i][j] = 'A' + rune(rand.IntN(26))
			}
		}
	}
}

func CountPatchMatch(patch [][]rune, board [][]rune) (count int) {
	l := len(patch)
	r := l / 2
	for i := r; i < len(board) - r; i++ {
		for j := r; j < len(board[0]) - r; j++ {
			pt := Point{i, j}
			if PatchMatch(patch, board, pt) {
				count++
			}
		}
	}
	return
}

func CountCrossWords(word []rune, board [][]rune) (count int) {
	patch000 := MakeCrossPatch(word)
	count += CountPatchMatch(patch000, board)

	patch090 := RotatePatch(patch000)
	count += CountPatchMatch(patch090, board)

	patch180 := RotatePatch(patch090)
	count += CountPatchMatch(patch180, board)

	patch270 := RotatePatch(patch180)
	count += CountPatchMatch(patch270, board)

	return
}
