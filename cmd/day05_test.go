package cmd_test

import (
	"log/slog"
	"reflect"
	"strings"
	"testing"

	"github.com/dusktreader/advent-of-code-2024/cmd"
	"github.com/dusktreader/advent-of-code-2024/util"
)

func TestAddRule(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	got := cmd.MakeManual()
	got.AddRule(47, 53)
	got.AddRule(97, 13)
	got.AddRule(97, 61)
	got.AddRule(97, 47)

	want := cmd.Manual{
		Rules: util.MakeDag(
			util.MakePair(47, 53),
			util.MakePair(97, 13),
			util.MakePair(97, 61),
			util.MakePair(97, 47),
		),
		Updates: []*cmd.Update{},
		ValidCheckSum: 0,
	}

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("Add rule failed: wanted %+v, got %+v", want, got)
	}
}

func TestAddPages(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	got := cmd.MakeManual()
	got.AddUpdate([]int{75, 47, 61, 53})

	want := cmd.Manual{
		Rules:   util.MakeDag[int](),
		Updates: []*cmd.Update{
			{
				Pages:       []int{75, 47, 61, 53},
			},
		},
	}

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("AddPages failed: wanted %+v, got %+v", want, got)
	}
}

func TestValidate(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	m := cmd.MakeManual()
	m.AddRule(47, 53)
	m.AddRule(97, 13)
	m.AddRule(97, 61)
	m.AddRule(97, 47)
	m.AddRule(75, 29)
	m.AddRule(61, 13)
	m.AddRule(75, 53)
	m.AddRule(29, 13)
	m.AddRule(97, 29)
	m.AddRule(53, 29)
	m.AddRule(61, 53)
	m.AddRule(97, 53)
	m.AddRule(61, 29)
	m.AddRule(47, 13)
	m.AddRule(75, 47)
	m.AddRule(97, 75)
	m.AddRule(47, 61)
	m.AddRule(75, 61)
	m.AddRule(47, 29)
	m.AddRule(75, 13)
	m.AddRule(53, 13)

	m.AddUpdate([]int{75, 47, 61, 53, 29})
	m.AddUpdate([]int{97, 61, 53, 29, 13})
	m.AddUpdate([]int{75, 29, 13})
	m.AddUpdate([]int{75, 97, 47, 61, 53})
	m.AddUpdate([]int{61, 13, 29})
	m.AddUpdate([]int{97, 13, 75, 29, 47})

	m.Validate()

	want := []bool{true, true, true, false, false, false}
	got := make([]bool, len(m.Updates))
	for i, u := range m.Updates {
		got[i] = u.Valid
	}

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("Validate failed to identify invalid updates: wanted %+v, got %+v", want, got)
	}

	if m.ValidCheckSum != 143 {
		t.Fatalf("Validate failed computing ValidCheckSum: wanted 143, got %v", m.ValidCheckSum)
	}
}

func TestAmend(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	m := cmd.MakeManual()
	m.AddRule(47, 53)
	m.AddRule(97, 13)
	m.AddRule(97, 61)
	m.AddRule(97, 47)
	m.AddRule(75, 29)
	m.AddRule(61, 13)
	m.AddRule(75, 53)
	m.AddRule(29, 13)
	m.AddRule(97, 29)
	m.AddRule(53, 29)
	m.AddRule(61, 53)
	m.AddRule(97, 53)
	m.AddRule(61, 29)
	m.AddRule(47, 13)
	m.AddRule(75, 47)
	m.AddRule(97, 75)
	m.AddRule(47, 61)
	m.AddRule(75, 61)
	m.AddRule(47, 29)
	m.AddRule(75, 13)
	m.AddRule(53, 13)

	m.AddUpdate([]int{75, 47, 61, 53, 29})
	m.AddUpdate([]int{97, 61, 53, 29, 13})
	m.AddUpdate([]int{75, 29, 13})
	m.AddUpdate([]int{75, 97, 47, 61, 53})
	m.AddUpdate([]int{61, 13, 29})
	m.AddUpdate([]int{97, 13, 75, 29, 47})

	m.Validate()
	m.Amend()

	want := []bool{false, false, false, true, true, true}
	got := make([]bool, len(m.Updates))
	for i, u := range m.Updates {
		got[i] = u.Amended
	}

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("Validate failed to identify amended updates: wanted %+v, got %+v", want, got)
	}

	if m.AmendCheckSum != 123 {
		t.Fatalf("Validate failed computing AmendCheckSum: wanted 143, got %v", m.ValidCheckSum)
	}
}

func TestParseRuleSuccess(t *testing.T) {
	wantLeft  := 47
	wantRight := 53
	gotLeft, gotRight, err := cmd.ParseRule("47|53")
	if err != nil {
		t.Fatalf("Unexpected error from ParseRule: %#v", err)
	} else if wantLeft != gotLeft || wantRight != gotRight {
		t.Fatalf("ParseRule failed: wantLeft=%v, gotLeft=%v; wantRight=%v, gotRight=%v", wantLeft, gotLeft, wantRight, gotRight)
	}
}

func TestParseRuleFail(t *testing.T) {
	_, _, err := cmd.ParseRule("47")
	if err == nil {
		t.Fatalf("Did not get expected error from ParseRule: %#v", err)
	} else if !strings.Contains(err.Error(), "Found a malformed rule") {
		t.Fatalf(`ParseRule did not fail as expected on split`)
	}

	_, _, err = cmd.ParseRule("47|53|99")
	if err == nil {
		t.Fatalf("Did not get expected error from ParseRule: %#v", err)
	} else if !strings.Contains(err.Error(), "Found a malformed rule") {
		t.Fatalf(`ParseRule did not fail as expected on split`)
	}

	_, _, err = cmd.ParseRule("forty-seven|53")
	if err == nil {
		t.Fatalf("Did not get expected error from ParseRule: %#v", err)
	} else if !strings.Contains(err.Error(), "Left side was not a number") {
		t.Fatalf(`ParseRule did not fail as expected on left side`)
	}

	_, _, err = cmd.ParseRule("47|fifty-three")
	if err == nil {
		t.Fatalf("Did not get expected error from ParseRule: %#v", err)
	} else if !strings.Contains(err.Error(), "Right side was not a number") {
		t.Fatalf(`ParseRule did not fail as expected on right side`)
	}
}

func TestParsePagesSuccess(t *testing.T) {
	want := []int{75, 57, 61, 53, 29}
	got, err := cmd.ParsePages("75,57,61,53,29")

	if err != nil {
		t.Fatalf("Unexpected error from ParsePages: %#v", err)
	} else if !reflect.DeepEqual(want, got) {
		t.Fatalf("ParsePages failed: want=%+v, got=%+v", want, got)
	}
}

func TestParsePagesFail(t *testing.T) {
	_, err := cmd.ParsePages("75,57,sixty-one,53,29")
	if err == nil {
		t.Fatalf("Did not get expected error from ParsePages: %#v", err)
	} else if !strings.Contains(err.Error(), "Location was not a number at position 2: sixty-one") {
		t.Fatalf(`ParsePages did not fail as expected on split`)
	}
}

func TestParseInput(t *testing.T) {
	inputStr := `
		47|53
		97|13
		97|61
		97|47
		75|29
		61|13
		75|53
		29|13
		97|29
		53|29
		61|53
		97|53
		61|29
		47|13
		75|47
		97|75
		47|61
		75|61
		47|29
		75|13
		53|13

		75,47,61,53,29
		97,61,53,29,13
		75,29,13
		75,97,47,61,53
		61,13,29
		97,13,75,29,47
	`
	got, err := cmd.ParseInput(inputStr)
	if err != nil {
		t.Fatalf("Unexpected error from ParseInput: %#v", err)
	}

	want := cmd.MakeManual()
	want.AddRule(47, 53)
	want.AddRule(97, 13)
	want.AddRule(97, 61)
	want.AddRule(97, 47)
	want.AddRule(75, 29)
	want.AddRule(61, 13)
	want.AddRule(75, 53)
	want.AddRule(29, 13)
	want.AddRule(97, 29)
	want.AddRule(53, 29)
	want.AddRule(61, 53)
	want.AddRule(97, 53)
	want.AddRule(61, 29)
	want.AddRule(47, 13)
	want.AddRule(75, 47)
	want.AddRule(97, 75)
	want.AddRule(47, 61)
	want.AddRule(75, 61)
	want.AddRule(47, 29)
	want.AddRule(75, 13)
	want.AddRule(53, 13)

	want.AddUpdate([]int{75, 47, 61, 53, 29})
	want.AddUpdate([]int{97, 61, 53, 29, 13})
	want.AddUpdate([]int{75, 29, 13})
	want.AddUpdate([]int{75, 97, 47, 61, 53})
	want.AddUpdate([]int{61, 13, 29})
	want.AddUpdate([]int{97, 13, 75, 29, 47})

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("ParseInput failed: wanted %+v, got %+v", want, got)
	}
}
