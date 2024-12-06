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
	rootCmd.AddCommand(d5Cmd)
	d5Cmd.PersistentFlags().StringP("input-file", "i", "", "Get input from a file instead of stdin")
	d5Cmd.AddCommand(d5p1Cmd)
	d5Cmd.AddCommand(d5p2Cmd)
	d5p2Cmd.Flags().BoolP("annotate", "a", false, "Annotate redacted input string")
	d5Cmd.AddCommand(mermaidCmd)
	d5Cmd.AddCommand(dotCmd)
	d5Cmd.AddCommand(validateCmd)
}

var d5Cmd = &cobra.Command{
	Use:   "day05",
	Short: "Day 5 Solutions",
	Long:  "The solutions for day 5 of Advent of Code 2025",
	Run:   d5Main,
}

var d5p1Cmd = &cobra.Command{
	Use:   "part1",
	Short: "Day 5, 1 Solution",
	Long:  "The solution for day 5, part 1 of Advent of Code 2025",
	Run:   d5p1Main,
}

var d5p2Cmd = &cobra.Command{
	Use:   "part2",
	Short: "Day 5, 2 Solution",
	Long:  "The solution for day 5, part 2 of Advent of Code 2025",
	Run:   d5p2Main,
}

var mermaidCmd = &cobra.Command{
	Use:   "mermaid",
	Short: "Mermaid",
	Long:  "Mermaid",
	Run:   mermaidMain,
}

var dotCmd = &cobra.Command{
	Use:   "dot",
	Short: "Dot",
	Long:  "Dot",
	Run:   dotMain,
}

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate",
	Long:  "Validate",
	Run:   validateMain,
}

func d5Main(cmd *cobra.Command, args []string){
	_ = cmd.Help()
}

func mermaidMain(cmd *cobra.Command, args []string) {
	inputStr, err := loadInput(cmd, args)
	MaybeDie(err)

	manual, err := ParseInput(inputStr)
	MaybeDie(err)

	fmt.Printf("%v\n", manual.Rules.Mermaid())
}

func dotMain(cmd *cobra.Command, args []string) {
	inputStr, err := loadInput(cmd, args)
	MaybeDie(err)

	manual, err := ParseInput(inputStr)
	MaybeDie(err)

	fmt.Printf("%v\n", manual.Rules.Dot())
}

func validateMain(cmd *cobra.Command, args []string) {
	inputStr, err := loadInput(cmd, args)
	MaybeDie(err)

	manual, err := ParseInput(inputStr)
	MaybeDie(err)

	if !manual.Rules.HasCycle() {
		Die("Rules DAG has at least one cycle")
	}
}

func d5p1Main(cmd *cobra.Command, args []string) {
	inputStr, err := loadInput(cmd, args)
	MaybeDie(err)

	manual, err := ParseInput(inputStr)
	MaybeDie(err)

	manual.Validate()
	slog.Debug("Results:", "ValidCheckSum", manual.ValidCheckSum)
	fmt.Printf("%v\n", manual.ValidCheckSum)
}

func d5p2Main(cmd *cobra.Command, args []string) {
	inputStr, err := loadInput(cmd, args)
	MaybeDie(err)

	manual, err := ParseInput(inputStr)
	MaybeDie(err)

	manual.Validate()
	manual.Amend()
	slog.Debug("Results:", "AmendCheckSum", manual.AmendCheckSum)
	fmt.Printf("%v\n", manual.AmendCheckSum)
}

type Update struct {
	Pages       []int
	Valid       bool
	Amended     bool
}

func MakeUpdate() (Update) {
	return Update{
		Pages:       make([]int, 0, 10),
	}
}

type RuleSet struct {
	contents map[int]util.Set[int]
}

type Manual struct {
	Rules          util.DAG[int]
	Updates        []*Update
	ValidCheckSum  int
	AmendCheckSum  int
}

func MakeManual() (Manual) {
	return Manual{
		Rules:       util.MakeDag[int](),
		Updates:     make([]*Update, 0, 10),
	}
}

func (m *Manual) AddRule(left int, right int) {
	m.Rules.AddEdge(left, right)
}

func (m *Manual) AddUpdate(pages []int) {
	u := MakeUpdate()
	u.Pages = make([]int, len(pages))
	for i, v := range pages {
		u.Pages[i] = v
	}
	m.Updates = append(m.Updates, &u)
}

func (u *Update) Validate(rules util.DAG[int]) int {
	u.Valid = rules.IsSortedDumb(u.Pages)
	if !u.Valid {
		return 0
	}
	return u.Pages[len(u.Pages) / 2]
}

func (u *Update) Amend(rules util.DAG[int]) int {
	if u.Valid {
		return 0
	}
	newPages := rules.SortDumb(u.Pages)
	u.Pages = newPages
	u.Amended = true
	return u.Pages[len(u.Pages) / 2]
}

func (m *Manual) Validate() {
	for _, u := range m.Updates {
		m.ValidCheckSum += u.Validate(m.Rules)
	}
}

func (m *Manual) Amend() {
	for _, u := range m.Updates {
		m.AmendCheckSum += u.Amend(m.Rules)
	}
}

func ParseRule(line string) (int, int, error) {
	parts := strings.Split(line, "|")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("Found a malformed rule: %v", line)
	}

	left, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("Left side was not a number: %v", parts[0])
	}

	right, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("Right side was not a number: %v", parts[1])
	}
	return left, right, nil
}

func ParsePages(line string) ([]int, error) {
	parts := strings.Split(line, ",")
	pages := make([]int, len(parts))
	for i, part := range parts {
		val, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("Location was not a number at position %v: %v", i, part)
		}
		pages[i] = val
	}
	return pages, nil
}

func ParseInput(inputStr string) (Manual, error) {
	inputStr = strings.TrimSpace(inputStr)
	manual := MakeManual()

	lines := strings.Split(inputStr, "\n")
	onSinai := true

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			onSinai = false
		} else if onSinai {
			left, right, err := ParseRule(line)
			if err != nil {
				return manual, fmt.Errorf("Bad rule on line %v: %#v", i, err)
			}
			manual.AddRule(left, right)
		} else {
			pages, err := ParsePages(line)
			if err != nil {
				return manual, fmt.Errorf("Bad page list on line %v: %#v", i, err)
			}
			manual.AddUpdate(pages)
		}
	}

	return manual, nil
}
