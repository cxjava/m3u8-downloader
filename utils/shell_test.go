package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConcatFiles(t *testing.T) {
	dir := t.TempDir()
	check(os.WriteFile(filepath.Join(dir, "a.ts"), []byte("hello "), 0644))
	check(os.WriteFile(filepath.Join(dir, "b.ts"), []byte("world"), 0644))

	if got := concatFiles(dir, "merged.tmp", false); got != "merged.tmp" {
		t.Fatalf("concatFiles() = %q, want %q", got, "merged.tmp")
	}

	data, err := os.ReadFile(filepath.Join(dir, "merged.tmp"))
	check(err)
	if string(data) != "hello world" {
		t.Fatalf("merged contents = %q, want %q", string(data), "hello world")
	}
	if _, err := os.Stat(filepath.Join(dir, "a.ts")); err != nil {
		t.Fatalf("a.ts should exist: %v", err)
	}
}

func TestConcatFilesDelete(t *testing.T) {
	dir := t.TempDir()
	check(os.WriteFile(filepath.Join(dir, "a.ts"), []byte("x"), 0644))
	check(os.WriteFile(filepath.Join(dir, "b.ts"), []byte("y"), 0644))

	concatFiles(dir, "merged.tmp", true)

	if _, err := os.Stat(filepath.Join(dir, "a.ts")); !os.IsNotExist(err) {
		t.Fatalf("a.ts should be deleted, err=%v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "b.ts")); !os.IsNotExist(err) {
		t.Fatalf("b.ts should be deleted, err=%v", err)
	}
	data, err := os.ReadFile(filepath.Join(dir, "merged.tmp"))
	check(err)
	if string(data) != "xy" {
		t.Fatalf("merged contents = %q, want %q", string(data), "xy")
	}
}
