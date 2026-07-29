package scope

import "testing"

func TestValidScopeMessage(t *testing.T) {
	if !Valid("config: loads defaults") {
		t.Error("expected valid")
	}
}

func TestValidSingleWordScope(t *testing.T) {
	if !Valid("a: message") {
		t.Error("expected valid")
	}
}

func TestInvalidNoColon(t *testing.T) {
	if Valid("just a message") {
		t.Error("expected invalid")
	}
}

func TestInvalidEmptyScope(t *testing.T) {
	if Valid(": message") {
		t.Error("expected invalid")
	}
}

func TestInvalidEmptyMessage(t *testing.T) {
	if Valid("scope: ") {
		t.Error("expected invalid")
	}
}

func TestInvalidEmptyString(t *testing.T) {
	if Valid("") {
		t.Error("expected invalid")
	}
}

func TestInvalidNoSpaceAfterColon(t *testing.T) {
	if Valid("scope:message") {
		t.Error("expected invalid")
	}
}
