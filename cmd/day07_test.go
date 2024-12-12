package cmd_test

import (
	"log/slog"
	"testing"

	"github.com/dusktreader/advent-of-code-2024/cmd"
	"github.com/dusktreader/advent-of-code-2024/util"
	"github.com/go-test/deep"
)

func TestParseEquations(t *testing.T) {
	txt := `
		190: 10 19
		3267: 81 40 27
		83: 17 5
		156: 15 6
		7290: 6 8 6 15
		161011: 16 10 13
		192: 17 8 14
		21037: 9 7 18 13
		292: 11 6 16 20
	`
	got, err := cmd.ParseEquations(txt)
	util.Unexpect(t, err)

	want := []cmd.Equation{
		{Left: 190,    Right: []int{10, 19}        },
		{Left: 3267,   Right: []int{81, 40, 27}    },
		{Left: 83,     Right: []int{17, 5}         },
		{Left: 156,    Right: []int{15, 6}         },
		{Left: 7290,   Right: []int{6, 8, 6, 15}   },
		{Left: 161011, Right: []int{16, 10, 13}    },
		{Left: 192,    Right: []int{17, 8, 14}     },
		{Left: 21037,  Right: []int{9, 7, 18, 13}  },
		{Left: 292,    Right: []int{11, 6, 16, 20} },
	}

	if diff := deep.Equal(want, got); diff != nil {
		t.Error(diff)
	}
}

func TestCatBasic(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	want := 12345
	got  := cmd.Cat(123, 45)
	if want != got {
		t.Fatalf("Cat failed: wanted %v, got %v", want, got)
	}

	want = 1
	got  = cmd.Cat(0, 1)
	if want != got {
		t.Fatalf("Cat failed: wanted %v, got %v", want, got)
	}

	want = 10
	got  = cmd.Cat(1, 0)
	if want != got {
		t.Fatalf("Cat failed: wanted %v, got %v", want, got)
	}

	want = 1230
	got  = cmd.Cat(123, 0)
	if want != got {
		t.Fatalf("Cat failed: wanted %v, got %v", want, got)
	}

	want = 1231
	got  = cmd.Cat(123, 1)
	if want != got {
		t.Fatalf("Cat failed: wanted %v, got %v", want, got)
	}

	want = 1230456
	got  = cmd.Cat(1230, 456)
	if want != got {
		t.Fatalf("Cat failed: wanted %v, got %v", want, got)
	}

	want = 123456
	got  = cmd.Cat(123, 456)
	if want != got {
		t.Fatalf("Cat failed: wanted %v, got %v", want, got)
	}

	want = 5696124843
	got  = cmd.Cat(5696, 124843)
	if want != got {
		t.Fatalf("Cat failed: wanted %v, got %v", want, got)
	}

	want = 987654321098765432
	got  = cmd.Cat(9876543210, 98765432)
	if want != got {
		t.Fatalf("Cat failed: wanted %v, got %v", want, got)
	}
}

func TestOperate2(t *testing.T) {
	eq   := cmd.Equation{Left: 0, Right: []int{2, 3}}
	ops  := 1 // Multiply
	want := 6
	got  := eq.Operate2(ops)
	if want != got {
		t.Fatalf("Operate failed: wanted %v, got %v", want, got)
	}

	eq   = cmd.Equation{Left: 0, Right: []int{2, 3, 5}}
	ops  = 2 // Multiply then add
	want = 11
	got  = eq.Operate2(ops)
	if want != got {
		t.Fatalf("Operate failed: wanted %v, got %v", want, got)
	}

	eq   = cmd.Equation{Left: 0, Right: []int{2, 3, 5, 7}}
	ops  = 2 // Add, multiply, add
	want = 32
	got  = eq.Operate2(ops)
	if want != got {
		t.Fatalf("Operate failed: wanted %v, got %v", want, got)
	}
}

func TestOperate3(t *testing.T) {
	eq   := cmd.Equation{Left: 0, Right: []int{2, 3}}
	ops  := (1 * 1) // MUL
	want := 6
	got  := eq.Operate3(ops)
	if want != got {
		t.Fatalf("Operate failed: wanted %v, got %v", want, got)
	}

	eq   = cmd.Equation{Left: 0, Right: []int{2, 3, 5}}
	ops  = (2 * 3) + (0 * 1) // CAT, MUL
	want = 28
	got  = eq.Operate3(ops)
	if want != got {
		t.Fatalf("Operate failed: wanted %v, got %v", want, got)
	}

	eq   = cmd.Equation{Left: 0, Right: []int{2, 3, 5, 7}}
	ops  = (0 * 9) + (2 * 3) + (1 * 1) // ADD, CAT, MUL
	want = 385
	got  = eq.Operate3(ops)
	if want != got {
		t.Fatalf("Operate failed: wanted %v, got %v", want, got)
	}
}

func TestProcess2(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	eq := cmd.Equation{
		Left:  3267,
		Right: []int{81, 40, 27},
		Sats:  util.MakeSet[int](),
	}
	eq.Process2()
	want := util.MakeSet(2, 1)
	got  := eq.Sats
	if !want.Eq(got) {
		t.Fatalf("Operate failed: wanted %v, got %v", want, got)
	}

	eq = cmd.Equation{
		Left:  83,
		Right: []int{17, 5},
		Sats:  util.MakeSet[int](),
	}
	eq.Process2()
	want = util.MakeSet[int]()
	got  = eq.Sats
	if !want.Eq(got) {
		t.Fatalf("Operate failed: wanted %v, got %v", want, got)
	}

	eq = cmd.Equation{
		Left:  292,
		Right: []int{11, 6, 16, 20},
		Sats:  util.MakeSet[int](),
	}
	eq.Process2()
	want = util.MakeSet(2)
	got  = eq.Sats
	if !want.Eq(got) {
		t.Fatalf("Operate failed: wanted %v, got %v", want, got)
	}

	eq = cmd.Equation{
		Left:  21037,
		Right: []int{9, 7, 18, 13},
		Sats:  util.MakeSet[int](),
	}
	eq.Process2()
	want = util.MakeSet[int]()
	got  = eq.Sats
	if !want.Eq(got) {
		t.Fatalf("Operate failed: wanted %v, got %v", want, got)
	}
}

func TestProcess3(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	eq := cmd.Equation{
		Left:  156,
		Right: []int{15, 6},
		Sats:  util.MakeSet[int](),
	}
	eq.Process3()
	want := util.MakeSet(cmd.CAT)
	got  := eq.Sats
	if !want.Eq(got) {
		t.Fatalf("Operate failed: wanted %v, got %v", want, got)
	}

	eq = cmd.Equation{
		Left:  7290,
		Right: []int{6, 8, 6, 15},
		Sats:  util.MakeSet[int](),
	}
	eq.Process3()
	want = util.MakeSet(
		(cmd.MUL * 9) + (cmd.CAT * 3) + (cmd.MUL * 1),
	)
	got  = eq.Sats
	if !want.Eq(got) {
		t.Fatalf("Operate failed: wanted %v, got %v", want, got)
	}
}

func TestEqTotal2(t *testing.T) {
	eqs, err := cmd.ParseEquations(`
		190: 10 19
		3267: 81 40 27
		83: 17 5
		156: 15 6
		7290: 6 8 6 15
		161011: 16 10 13
		192: 17 8 14
		21037: 9 7 18 13
		292: 11 6 16 20
	`)
	util.Unexpect(t, err)

	want := 3749
	got  := cmd.EqTotal2(eqs)
	if want != got {
		t.Fatalf("EqTotal failed: wanted %v, got %v", want, got)
	}
}

func TestEqTotal3(t *testing.T) {
	eqs, err := cmd.ParseEquations(`
		190: 10 19
		3267: 81 40 27
		83: 17 5
		156: 15 6
		7290: 6 8 6 15
		161011: 16 10 13
		192: 17 8 14
		21037: 9 7 18 13
		292: 11 6 16 20
	`)
	util.Unexpect(t, err)

	want := 11387
	got  := cmd.EqTotal3(eqs)
	if want != got {
		t.Fatalf("EqTotal failed: wanted %v, got %v", want, got)
	}
}
