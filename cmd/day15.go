package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/dusktreader/advent-of-code-2024/util"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(d15Cmd)
	d15Cmd.PersistentFlags().StringP("input-file", "i", "", "Get input from a file instead of stdin")
	d15Cmd.PersistentFlags().BoolP("small", "s", false, "Use smaller space")
	d15Cmd.PersistentFlags().CountP("visualize", "V", "Visualize space. Pass multiple to visualize more")
	d15Cmd.AddCommand(d15p1Cmd)
	d15Cmd.AddCommand(d15p2Cmd)
	d15p2Cmd.Flags().IntP("stretch-horizontal", "F", 2, "Horizontal stretch factor.")
	d15p2Cmd.Flags().IntP("stretch-vertical", "f", 1, "Vertical stretch factor.")
}

var d15Cmd = &cobra.Command{
	Use:   "day15",
	Short: "Day 15 Solutions",
	Long:  "The solutions for day 15 of Advent of Code 2025",
	Run:   d15Main,
}

var d15p1Cmd = &cobra.Command{
	Use:   "part1",
	Short: "Day 15, 1 Solution",
	Long:  "The solution for day 15, part 1 of Advent of Code 2025",
	Run:   d15p1Main,
}

var d15p2Cmd = &cobra.Command{
	Use:   "part2",
	Short: "Day 15, 2 Solution",
	Long:  "The solution for day 15, part 2 of Advent of Code 2025",
	Run:   d15p2Main,
}

func d15Main(cmd *cobra.Command, args []string){
	_ = cmd.Help()
}

func getViz(cmd *cobra.Command, args []string) int {
	viz, err := cmd.Flags().GetCount("visualize")
	MaybeDie(err)
	return viz
}

func getStretch(cmd *cobra.Command, args []string) (int, int) {
	wf, err := cmd.Flags().GetInt("stretch-horizontal")
	MaybeDie(err)
	if wf < 1 {
		Die("Invalid horizontal stretch factor: %v", wf)
	}

	hf, err := cmd.Flags().GetInt("stretch-vertical")
	MaybeDie(err)
	if hf < 1 {
		Die("Invalid vertical stretch factor: %v", hf)
	}
	return wf, hf
}

func d15p1Main(cmd *cobra.Command, args []string) {
	viz := getViz(cmd, args)

	inputStr, err := loadInput(cmd, args)
	MaybeDie(err)

	slog.Debug("Parsing warehouse")
	wh, err := ParseWarehouse(inputStr)
	MaybeDie(err)

	slog.Debug("Moving robot")
	wh.MoveAll(viz)

	if viz > 0 {
		slog.Debug("Visualizing final warehouse")
		fmt.Printf("%v\n", wh)
	}

	slog.Debug("Box count:", "Boxes", wh.Boxes.Size())
	gps := wh.GPS()

	slog.Debug("Results:", "GPS", gps)
	fmt.Printf("%v\n", gps)
}

func d15p2Main(cmd *cobra.Command, args []string) {
	wf, hf := getStretch(cmd, args)
	viz := getViz(cmd, args)

	inputStr, err := loadInput(cmd, args)
	MaybeDie(err)

	slog.Debug("Parsing warehouse")
	wh, err := ParseWarehouse(inputStr)
	MaybeDie(err)

	slog.Debug("Stretching warehouse")
	wh.Stretch(wf, hf)

	slog.Debug("Moving robot")
	wh.MoveAll(viz)

	if viz > 0 {
		slog.Debug("Visualizing final warehouse")
		fmt.Printf("%v\n", wh)
	}

	slog.Debug("Box count:", "Boxes", wh.Boxes.Size())
	gps := wh.GPS()

	slog.Debug("Results:", "GPS", gps)
	fmt.Printf("%v\n", gps)
}

type Warehouse struct {
	Sz    util.Size
	Walls util.Set[util.Point]
	Boxes util.Set[*util.Rect]
	Robot util.Point
	Moves util.Queue[util.Vector]
}

func (wh *Warehouse) Stretch(wf int, hf int) error {
	if wf < 1 || hf < 1 {
		return fmt.Errorf("Invalid stretch factors: %v, %v", wf, hf)
	}

	wh.Sz.W *= wf
	wh.Sz.H *= hf

	wh.Robot = util.MakePoint(wh.Robot.I * hf, wh.Robot.J * wf)

	newWalls := util.MakeSet[util.Point]()
	for wall := range wh.Walls.Iter() {
		r := util.Rect{
			O: util.MakePoint(
				wall.I * hf,
				wall.J * wf,
			),
			Sz: util.Size{W: wf, H: hf},
		}
		for newWall := range r.Iter() {
			newWalls.Add(newWall)
		}
	}
	wh.Walls = newWalls

	for box := range wh.Boxes.Iter() {
		box.O.I *= hf
		box.O.J *= wf
		box.Sz.H *= hf
		box.Sz.W *= wf
	}

	return nil
}

func (wh *Warehouse) MovableBoxes(pt util.Point, vt util.Vector) (util.Set[*util.Rect], error) {
	inBoxes := util.MakeSet[*util.Rect]()

	if wh.Walls.Has(pt) {
		return inBoxes, fmt.Errorf("Hit a wall at %v", pt)
	}

	inBox := wh.Boxed(pt)
	if inBox == nil {
		return inBoxes, nil
	}

	var recPtA, recPtB util.Point
	switch vt {
	case util.NORTH:
		recPtA = inBox.Tl()
		recPtB = inBox.Tr()
	case util.SOUTH:
		recPtA = inBox.Bl()
		recPtB = inBox.Br()
	case util.EAST:
		recPtA = inBox.Tr()
		recPtB = inBox.Br()
	case util.WEST:
		recPtA = inBox.Tl()
		recPtB = inBox.Bl()
	}

	recBoxesA, err := wh.MovableBoxes(recPtA.Add(vt), vt)
	if err != nil {
		return inBoxes, err
	}

	recBoxesB, err := wh.MovableBoxes(recPtB.Add(vt), vt)
	if err != nil {
		return inBoxes, err
	}

	inBoxes = inBoxes.Un(recBoxesA).Un(recBoxesB)
	inBoxes.Add(inBox)
	return inBoxes, nil
}

func (wh *Warehouse) MoveRobot() error {
	vt, err := wh.Moves.Pop()
	if err != nil {
		return fmt.Errorf("No moves left")
	}

	boxes, err := wh.MovableBoxes(wh.Robot.Add(vt), vt)
	if err != nil {
		return nil
	}

	for box := range boxes.Iter() {
		box.O = box.O.Add(vt)
	}
	wh.Robot = wh.Robot.Add(vt)
	return nil
}

func Cls() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func Idle() {
	time.Sleep(100_000_000 * time.Nanosecond)
}

func (wh *Warehouse) MoveAll(viz int) {
	if viz == 3 {
		Cls()
		fmt.Printf("%v\n\n", wh)
		Idle()
		Cls()
	}
	for {
		err := wh.MoveRobot()
		if err != nil {
			break
		}
		if viz == 2 {
			fmt.Printf("%v\n\n", wh)
		} else if viz == 3 {
			fmt.Printf("%v\n\n", wh)
			Idle()
			Cls()
		}
	}
}

func (wh *Warehouse) Boxed(pt util.Point) *util.Rect {
	inBoxes := make([]*util.Rect, 0)
	for box := range wh.Boxes.Iter() {
		if box.Has(pt) {
			inBoxes = append(inBoxes, box)
		}
	}

	if len(inBoxes) == 0 {
		return nil
	} else if len(inBoxes) > 1 {
		message := fmt.Sprintf("Found multiple boxes at %v: %+v", pt, inBoxes)
		panic(message)
	}
	return inBoxes[0]
}

func (wh *Warehouse) RenderPt(pt util.Point) rune {
	if wh.Robot == pt {
		return '@'
	} else if wh.Walls.Has(pt) {
		return '▦'
	}

	inBox := wh.Boxed(pt)
	if inBox == nil {
		return ' '
	}

	if inBox.Sz.W == 1 && inBox.Sz.H == 1 {
		return '□'
	}

	if inBox.Sz.W == 1 {
		if pt.I == inBox.Tl().I {
			return '^'
		} else if pt.I == inBox.Bl().I {
			return 'v'
		} else {
			return '‖'
		}
	}

	if inBox.Sz.H == 1 {
		if pt.J == inBox.Tl().J {
			return '['
		} else if pt.J == inBox.Tr().J {
			return ']'
		} else {
			return '='
		}
	}


	if pt == inBox.Tl() {
		return '┌'
	} else if pt == inBox.Tr(){
		return '┐'
	} else if pt == inBox.Bl(){
		return '└'
	} else if pt == inBox.Br(){
		return '┘'
	} else if pt.I == inBox.Tl().I || pt.I == inBox.Bl().I {
		return '─'
	} else if pt.J == inBox.Tl().J || pt.J == inBox.Tr().J {
		return '│'
	} else {
		return '☷'
	}
}

func (wh *Warehouse) GPS() (gps int) {
	for box := range wh.Boxes.Iter() {
		gps += 100 * box.O.I + box.O.J
	}
	return
}

func (wh *Warehouse) String() string {
	var sb strings.Builder
	for i := 0; i < wh.Sz.H; i++ {
		for j := 0; j < wh.Sz.W; j++ {
			pt := util.MakePoint(i, j)
			sb.WriteRune(wh.RenderPt(pt))
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}

func ParseWarehouse(inputStr string) (*Warehouse, error) {
	inputStr = strings.TrimSpace(inputStr)
	lines := strings.Split(inputStr, "\n")

	parsingMap := true
	wh := Warehouse{
		Sz: util.Size{W: len(lines[0]), H: 0},
		Walls: util.MakeSet[util.Point](),
		Boxes: util.MakeSet[*util.Rect](),
		Moves: *util.MakeQueue[util.Vector](),
		Robot: util.MakePoint(0, 0),
	}

	for i, line := range lines {
		if line == "" {
			parsingMap = false
			continue
		}

		if parsingMap {
			for j, rn := range line {
				pt := util.MakePoint(i, j)
				switch rn {
				case '#':
					wh.Walls.Add(pt)
				case 'O':
					box := util.MakeUnitRect(pt)
					wh.Boxes.Add(&box)
				case '@':
					wh.Robot = pt
				}
			}
			wh.Sz.H++
		} else {
			for _, rn := range line {
				switch rn {
				case '^':
					wh.Moves.Push(util.NORTH)
				case '>':
					wh.Moves.Push(util.EAST)
				case 'v':
					wh.Moves.Push(util.SOUTH)
				case '<':
					wh.Moves.Push(util.WEST)
				default:
					return nil, fmt.Errorf("Unknown direction in moves: %v", rn)
				}
			}
		}
	}

	return &wh, nil
}

