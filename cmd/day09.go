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
	rootCmd.AddCommand(d9Cmd)
	d9Cmd.PersistentFlags().StringP("input-file", "i", "", "Get input from a file instead of stdin")
	d9Cmd.PersistentFlags().StringP("output-file", "o", "", "Draw the final map and store it in the file")
	d9Cmd.AddCommand(d9p1Cmd)
	d9Cmd.AddCommand(d9p2Cmd)
}

var d9Cmd = &cobra.Command{
	Use:   "day09",
	Short: "Day 9 Solutions",
	Long:  "The solutions for day 9 of Advent of Code 2025",
	Run:   d9Main,
}

var d9p1Cmd = &cobra.Command{
	Use:   "part1",
	Short: "Day 9, 1 Solution",
	Long:  "The solution for day 9, part 1 of Advent of Code 2025",
	Run:   d9p1Main,
}

var d9p2Cmd = &cobra.Command{
	Use:   "part2",
	Short: "Day 9, 2 Solution",
	Long:  "The solution for day 9, part 2 of Advent of Code 2025",
	Run:   d9p2Main,
}

func d9Main(cmd *cobra.Command, args []string){
	_ = cmd.Help()
}

func d9p1Main(cmd *cobra.Command, args []string) {
	inputStr, err := loadInput(cmd, args)
	MaybeDie(err)

	files, err := ParseFiles(inputStr)
	MaybeDie(err)

	compact := CompactFilesDense(files)
	checksum := ComputeChecksumCompact(compact)

	slog.Debug("Results:", "Checksum", checksum)
	fmt.Printf("%v\n", checksum)
}

func d9p2Main(cmd *cobra.Command, args []string) {
	inputStr, err := loadInput(cmd, args)
	MaybeDie(err)

	files, err := ParseFiles(inputStr)
	MaybeDie(err)

	slots := ExpandSlots(files)
	slots = CompactFilesSparse(slots)
	checksum := ComputeChecksumSparse(slots)

	slog.Debug("Results:", "Checksum", checksum)
	fmt.Printf("%v\n", checksum)
}

func ParseFiles(inputStr string) ([]int, error) {
	inputStr = strings.TrimSpace(inputStr)
	l := len(inputStr)
	files := make([]int, l)
	for i := 0; i < l; i++ {
		v := int(inputStr[i] - '0')
		if v < 0 || v > 9 {
			return []int{}, fmt.Errorf("Couldn't convert token at index %v", i)
		}
		files[i] = v
	}
	return files, nil
}

func ShoveFiles(compact []int, id int, ct int) []int {
	for range ct {
		compact = append(compact, id)
	}
	return compact
}

func CompactFilesDense(files []int) []int {
	compact := make([]int, 0, len(files))

	l, lId := 0, 0
	r := len(files) - 1
	rId := r / 2

	rCt := files[r]
	empty := 0
	for {
		if empty == 0 {
			lCt := files[l]
			compact = ShoveFiles(compact, lId, lCt)
			l += 2
			if l >= r {
				break
			}
			empty = files[l - 1]
			lId++
		} else if rCt == 0 {
			r -= 2
			if r <= l {
				break
			}
			rId--
			rCt = files[r]
		} else {
			m := util.MinI(empty, rCt)
			compact = ShoveFiles(compact, rId, m)
			empty -= m
			rCt -= m
		}
	}
	if rCt > 0 {
		compact = ShoveFiles(compact, rId, rCt)
	}
	return compact
}

const EMPTY = -1

func ExpandSlots(files []int) []FileSlot {
	slots := make([]FileSlot, len(files))
	id := 0
	for i := 0; ; i += 2 {
		slots[i] = FileSlot{Id: id, Ct: files[i]}
		if i + 1 >= len(files) {
			break
		}
		slots[i + 1] = FileSlot{Id: EMPTY, Ct: files[i + 1]}
		id++
	}
	return slots
}

func CondenseSlots(slots []FileSlot) []int {
	files := make([]int, 0, len(slots) * 2)
	for _, s := range slots {
		files = ShoveFiles(files, s.Id, s.Ct)
	}
	return files
}

type FileSlot struct {
	Id int
	Ct int
}

func CompactFilesSparse(slots []FileSlot) []FileSlot {
	for r := len(slots) - 1; r > 0; r-- {
		if slots[r].Id == EMPTY {
			continue
		}
		for l := 0; l < r; l++ {
			if slots[l].Id != EMPTY {
				continue
			}

			d := slots[l].Ct - slots[r].Ct
			if d >= 0 {
				slog.Debug("Found a slot!", "l", l, "l.Ct", slots[l].Ct)
				slots[l].Id = slots[r].Id
				slots[l].Ct = slots[r].Ct
				slots[r].Id = EMPTY
				if d > 0 {
					slog.Debug("Inserting a smaller empty!", "l+1", l + 1, "Ct", d)
					slots = util.Insert(slots, l + 1, FileSlot{Id: EMPTY, Ct: d})
				}
				slog.Debug("Now slots look like:", "slots", slots)
				break
			}
		}
	}
	return slots
}

func PrintSlots(slots []FileSlot) string {
	var builder strings.Builder
	for _, s := range slots {
		var c string
		if s.Id == EMPTY {
			c = "."
		} else {
			c = strconv.Itoa(s.Id)
		}
		builder.WriteString(strings.Repeat(c, s.Ct))
	}
	return builder.String()
}

func ParseDots(dots string) []int {
	files := make([]int, len(dots))
	for i := 0; i < len(dots); i++ {
		if dots[i] == '.' {
			files[i] = EMPTY
		} else {
			files[i] = int(dots[i] - '0')
		}
	}
	return files
}


func ComputeChecksumCompact(compact []int) int {
	sum := 0
	for i := range len(compact) {
		sum += i * compact[i]
	}
	return sum
}

func ComputeChecksumSparse(slots []FileSlot) int {
	sum := 0
	i := 0
	for _, s := range slots {
		if s.Id == EMPTY {
			i += s.Ct
			continue
		}
		for range s.Ct {
			sum += i * s.Id
			i++
		}
	}
	return sum
}
