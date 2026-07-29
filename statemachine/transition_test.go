package statemachine

import "testing"

func TestValidateEmptyFrom(t *testing.T) {
	err := Validate([]TransitionDef{
		{From: "", Event: "run", To: "running"},
	})
	if err == nil {
		t.Fatal("expected error for empty From")
	}
}

func TestValidateEmptyEvent(t *testing.T) {
	err := Validate([]TransitionDef{
		{From: "idle", Event: "", To: "running"},
	})
	if err == nil {
		t.Fatal("expected error for empty Event")
	}
}

func TestValidateEmptyTo(t *testing.T) {
	err := Validate([]TransitionDef{
		{From: "idle", Event: "run", To: ""},
	})
	if err == nil {
		t.Fatal("expected error for empty To")
	}
}

func TestValidatePass(t *testing.T) {
	err := Validate([]TransitionDef{
		{From: "idle", Event: "run", To: "running"},
	})
	if err != nil {
		t.Fatal(err)
	}
}
