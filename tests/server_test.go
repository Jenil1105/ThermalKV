package tests

import (
	"testing"
	"thermalkv/internal/server"
)

func TestExecuteCommandLifecycle(t *testing.T) {
	env := NewTestEnv(t)
	db := env.Store

	lines, shouldClose := server.ExecuteCommand(db, "SET user hello world")
	if shouldClose {
		t.Fatal("SET should not close the connection")
	}
	if len(lines) != 1 || lines[0] != "OK :)" {
		t.Fatalf("unexpected SET response: %#v", lines)
	}

	lines, shouldClose = server.ExecuteCommand(db, "GET user")
	if shouldClose {
		t.Fatal("GET should not close the connection")
	}
	if len(lines) != 1 || lines[0] != "hello world" {
		t.Fatalf("unexpected GET response: %#v", lines)
	}

	lines, shouldClose = server.ExecuteCommand(db, "KEYS")
	if shouldClose {
		t.Fatal("KEYS should not close the connection")
	}
	if len(lines) != 1 || lines[0] != "user" {
		t.Fatalf("unexpected KEYS response: %#v", lines)
	}

	lines, shouldClose = server.ExecuteCommand(db, "COUNT")
	if shouldClose {
		t.Fatal("COUNT should not close the connection")
	}
	if len(lines) != 1 || lines[0] != "1" {
		t.Fatalf("unexpected COUNT response: %#v", lines)
	}
}

func TestExecuteCommandValidationAndExit(t *testing.T) {
	env := NewTestEnv(t)
	db := env.Store

	lines, shouldClose := server.ExecuteCommand(db, "TTL missing")
	if shouldClose {
		t.Fatal("invalid TTL should not close the connection")
	}
	if len(lines) != 1 || lines[0] != "Usage: TTL key seconds" {
		t.Fatalf("unexpected TTL validation response: %#v", lines)
	}

	lines, shouldClose = server.ExecuteCommand(db, "   ")
	if shouldClose {
		t.Fatal("empty input should not close the connection")
	}
	if len(lines) != 1 || lines[0] != "Empty command" {
		t.Fatalf("unexpected empty-command response: %#v", lines)
	}

	lines, shouldClose = server.ExecuteCommand(db, "EXIT")
	if !shouldClose {
		t.Fatal("EXIT should signal the connection to close")
	}
	if len(lines) != 1 || lines[0] != "bye... ;) " {
		t.Fatalf("unexpected EXIT response: %#v", lines)
	}
}
