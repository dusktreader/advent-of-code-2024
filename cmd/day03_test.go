package cmd_test

import (
	"log/slog"
	"reflect"
	"strings"
	"testing"

	"github.com/dusktreader/advent-of-code-2024/cmd"
	"github.com/dusktreader/advent-of-code-2024/util"
)

func TestIsolatePairs(t *testing.T) {
	inputs := []string{
	    "mul(1,2)",       // valid
		"mul(0,999)",     // valid
		"mul(1,9999)",    // invalid: too many digits in right operand
		"mul(-1,2)",      // invalid: negative number in left operand
		"mul(1, 2)",      // invalid: space before right operand
		"mul( 1,2)",      // invalid: space before left operand
		"mul(1,2 )",      // invalid: space after right operand
		"mul(1 ,2)",      // invalid: space after left operand
		"mul( 1 ,2 )",    // invalid: space around operands
		"mul (1,2)",      // invalid: space after operator
		"foo(1,2)",       // invalid: non-matching operator
		"mul[1,2]",       // invalid: wrong brackets
		" mul(1,2) ",     // valid
		"!#$mul(1,2)#$3", // valid
		"mul(1,2 ",       // invalid: No closing paren
		"mul1,2)",        // invalid: No opening paren
	}
	inputStr := strings.Join(inputs, "")

	want := []util.Pair[int]{
		{1, 2},
		{0, 999},
		{1, 2},
		{1, 2},
	}
	got := cmd.IsolatePairs(inputStr)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf(`IsolatePairs gave a bad answer. Wanted %+v, got %+v`, want, got)
	}
}

func TestProcesPairs(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	instructions := []util.Pair[int]{
		{1, 2},
		{0, 999},
		{1, 2},
		{1, 2},
	}
	want := (1 * 2) + (0 * 999) + (1 * 2) + (1 * 2)
	got := cmd.ProcessPairs(instructions)
	if got != want {
		t.Fatalf(`ProcessPairs gave a bad answer. Wanted %v, got %v`, want, got)
	}
}

func TestRedact(t *testing.T) {
	inputStr := strings.TrimSpace(`
		don't()
		should not include this text
		do()
		include this stuff
		don't()
		leave this out
		don't()
		do ()
		do( )
		doo()
		this also shouldn't go in
		do()
		dont()
		don't ()
		don't( )
		should include this text
		don't()
		this should not be included
	`)

	want := strings.TrimSpace(`
		REDACTED
		include this stuff
		REDACTED
		dont()
		don't ()
		don't( )
		should include this text
		REDACTED
	`)

	got := cmd.Redact(inputStr, true)
	if got != want {
		t.Fatalf(`Redact gave a bad answer. Wanted %v, got %v`, want, got)
	}
}
