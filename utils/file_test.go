package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteToFileOverwritesExistingContent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "generated.txt")
	WriteToFile(path, "long content")
	WriteToFile(path, "short")

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if got := string(content); got != "short" {
		t.Fatalf("WriteToFile() = %q, want %q", got, "short")
	}
}

func TestCopyFileTruncatesDestination(t *testing.T) {
	dir := t.TempDir()
	source := filepath.Join(dir, "source.txt")
	destination := filepath.Join(dir, "destination.txt")
	if err := os.WriteFile(source, []byte("new"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(destination, []byte("old trailing content"), 0644); err != nil {
		t.Fatal(err)
	}

	if _, err := CopyFile(source, destination); err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(destination)
	if err != nil {
		t.Fatal(err)
	}
	if got := string(content); got != "new" {
		t.Fatalf("CopyFile() = %q, want %q", got, "new")
	}
}

func TestReplaceFileContentDoesNotAppend(t *testing.T) {
	path := filepath.Join(t.TempDir(), "launch.json")
	if err := os.WriteFile(path, []byte("before TOKEN after"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := ReplaceFileContent(path, "TOKEN", "value"); err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if got := string(content); got != "before value after" {
		t.Fatalf("ReplaceFileContent() = %q, want %q", got, "before value after")
	}
}

func TestCopyDirCopiesNestedFiles(t *testing.T) {
	dir := t.TempDir()
	source := filepath.Join(dir, "source")
	destination := filepath.Join(dir, "destination")
	if err := os.MkdirAll(filepath.Join(source, "nested"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(source, "nested", "file.txt"), []byte("content"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := CopyDir(source, destination); err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(filepath.Join(destination, "nested", "file.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if got := string(content); got != "content" {
		t.Fatalf("copied content = %q, want %q", got, "content")
	}
}
