package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsCgearProject(t *testing.T) {
	project := t.TempDir()
	for _, path := range []string{
		filepath.Join(project, "CMakeLists.txt"),
		filepath.Join(project, "src", "CMakeLists.txt"),
		filepath.Join(project, "test", "CMakeLists.txt"),
	} {
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, nil, 0644); err != nil {
			t.Fatal(err)
		}
	}

	if !IsCgearProject(project) {
		t.Fatal("IsCgearProject() = false, want true")
	}
	if err := os.Remove(filepath.Join(project, "test", "CMakeLists.txt")); err != nil {
		t.Fatal(err)
	}
	if IsCgearProject(project) {
		t.Fatal("IsCgearProject() = true with a required file missing")
	}
}
