package parser

import (
	"testing"
)

func TestParsePassEvent(t *testing.T) {
	line := `{"Action":"pass","Package":"mypkg","Test":"TestFoo","Elapsed":0.1}`
	e, err := Parse(line)
	if err != nil {
		t.Fatal(err)
	}
	if e.Action != "pass" {
		t.Errorf("want pass, got %s", e.Action)
	}
	if e.Package != "mypkg" {
		t.Errorf("want mypkg, got %s", e.Package)
	}
	if e.Test != "TestFoo" {
		t.Errorf("want TestFoo, got %s", e.Test)
	}
	if e.Elapsed != 0.1 {
		t.Errorf("want 0.1, got %f", e.Elapsed)
	}
}

func TestParseFailEvent(t *testing.T) {
	line := `{"Action":"fail","Package":"mypkg","Test":"TestBar"}`
	e, err := Parse(line)
	if err != nil {
		t.Fatal(err)
	}
	if e.Action != "fail" {
		t.Errorf("want fail, got %s", e.Action)
	}
}

func TestParseSkipEvent(t *testing.T) {
	line := `{"Action":"skip","Package":"mypkg","Test":"TestBaz"}`
	e, err := Parse(line)
	if err != nil {
		t.Fatal(err)
	}
	if e.Action != "skip" {
		t.Errorf("want skip, got %s", e.Action)
	}
}

func TestParseOutputEvent(t *testing.T) {
	line := `{"Action":"output","Package":"mypkg","Test":"TestFoo","Output":"ok\\n"}`
	e, err := Parse(line)
	if err != nil {
		t.Fatal(err)
	}
	if e.Action != "output" {
		t.Errorf("want output, got %s", e.Action)
	}
}

func TestParseRunEvent(t *testing.T) {
	line := `{"Action":"run","Package":"mypkg","Test":"TestFoo"}`
	e, err := Parse(line)
	if err != nil {
		t.Fatal(err)
	}
	if e.Action != "run" {
		t.Errorf("want run, got %s", e.Action)
	}
}

func TestParsePauseEvent(t *testing.T) {
	line := `{"Action":"pause","Package":"mypkg"}`
	e, err := Parse(line)
	if err != nil {
		t.Fatal(err)
	}
	if e.Action != "pause" {
		t.Errorf("want pause, got %s", e.Action)
	}
}

func TestParseContEvent(t *testing.T) {
	line := `{"Action":"cont","Package":"mypkg"}`
	e, err := Parse(line)
	if err != nil {
		t.Fatal(err)
	}
	if e.Action != "cont" {
		t.Errorf("want cont, got %s", e.Action)
	}
}

func TestParseInvalidJSON(t *testing.T) {
	_, err := Parse("not json")
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestParseEmptyLine(t *testing.T) {
	_, err := Parse("")
	if err == nil {
		t.Fatal("expected error for empty line")
	}
}
