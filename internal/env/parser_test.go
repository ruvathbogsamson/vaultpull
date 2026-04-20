package env

import (
	"strings"
	"testing"
)

func TestParse_BasicPairs(t *testing.T) {
	input := "FOO=bar\nBAZ=qux\n"
	m, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if m["FOO"] != "bar" {
		t.Errorf("FOO: got %q", m["FOO"])
	}
	if m["BAZ"] != "qux" {
		t.Errorf("BAZ: got %q", m["BAZ"])
	}
}

func TestParse_IgnoresComments(t *testing.T) {
	input := "# this is a comment\nFOO=bar\n"
	m, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := m["# this is a comment"]; ok {
		t.Error("comment line should not be parsed as key")
	}
	if m["FOO"] != "bar" {
		t.Errorf("FOO: got %q", m["FOO"])
	}
}

func TestParse_IgnoresBlankLines(t *testing.T) {
	input := "\n\nFOO=bar\n\n"
	m, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(m) != 1 {
		t.Errorf("expected 1 key, got %d", len(m))
	}
}

func TestParse_DoubleQuotedValue(t *testing.T) {
	input := `SECRET="hello world"`
	m, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if m["SECRET"] != "hello world" {
		t.Errorf("SECRET: got %q", m["SECRET"])
	}
}

func TestParse_SingleQuotedValue(t *testing.T) {
	input := `TOKEN='abc123'`
	m, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if m["TOKEN"] != "abc123" {
		t.Errorf("TOKEN: got %q", m["TOKEN"])
	}
}

func TestParse_MalformedLineSkipped(t *testing.T) {
	input := "NOEQUALS\nFOO=bar\n"
	m, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := m["NOEQUALS"]; ok {
		t.Error("malformed line should be skipped")
	}
	if m["FOO"] != "bar" {
		t.Errorf("FOO: got %q", m["FOO"])
	}
}

func TestParse_DuplicateKeyLastWins(t *testing.T) {
	input := "FOO=first\nFOO=second\n"
	m, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if m["FOO"] != "second" {
		t.Errorf("expected last value 'second', got %q", m["FOO"])
	}
}

func TestParse_EmptyValue(t *testing.T) {
	input := "EMPTY=\n"
	m, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if v, ok := m["EMPTY"]; !ok || v != "" {
		t.Errorf("expected empty string value, got %q (ok=%v)", v, ok)
	}
}
