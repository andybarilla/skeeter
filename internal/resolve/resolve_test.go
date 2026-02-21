package resolve

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDir(t *testing.T) {
	t.Run("flag takes precedence", func(t *testing.T) {
		tmpDir := t.TempDir()
		skeeterDir := filepath.Join(tmpDir, ".skeeter")
		if err := os.Mkdir(skeeterDir, 0755); err != nil {
			t.Fatal(err)
		}

		os.Setenv("SKEETER_DIR", "/should/be/ignored")
		defer os.Unsetenv("SKEETER_DIR")

		result, err := Dir(skeeterDir)
		if err != nil {
			t.Fatalf("Dir: %v", err)
		}

		abs, _ := filepath.Abs(skeeterDir)
		if result != abs {
			t.Errorf("Dir = %q, want %q", result, abs)
		}
	})

	t.Run("environment variable when no flag", func(t *testing.T) {
		tmpDir := t.TempDir()
		skeeterDir := filepath.Join(tmpDir, ".skeeter")
		if err := os.Mkdir(skeeterDir, 0755); err != nil {
			t.Fatal(err)
		}

		os.Setenv("SKEETER_DIR", skeeterDir)
		defer os.Unsetenv("SKEETER_DIR")

		result, err := Dir("")
		if err != nil {
			t.Fatalf("Dir: %v", err)
		}

		abs, _ := filepath.Abs(skeeterDir)
		if result != abs {
			t.Errorf("Dir = %q, want %q", result, abs)
		}
	})

	t.Run("walk up to find .skeeter", func(t *testing.T) {
		tmpDir := t.TempDir()
		skeeterDir := filepath.Join(tmpDir, ".skeeter")
		subDir := filepath.Join(tmpDir, "subdir", "nested")
		if err := os.MkdirAll(skeeterDir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.MkdirAll(subDir, 0755); err != nil {
			t.Fatal(err)
		}

		oldWd, _ := os.Getwd()
		if err := os.Chdir(subDir); err != nil {
			t.Fatal(err)
		}
		defer os.Chdir(oldWd)

		os.Unsetenv("SKEETER_DIR")

		result, err := Dir("")
		if err != nil {
			t.Fatalf("Dir: %v", err)
		}

		abs, _ := filepath.Abs(skeeterDir)
		if result != abs {
			t.Errorf("Dir = %q, want %q", result, abs)
		}
	})

	t.Run("error when not found", func(t *testing.T) {
		tmpDir := t.TempDir()

		oldWd, _ := os.Getwd()
		if err := os.Chdir(tmpDir); err != nil {
			t.Fatal(err)
		}
		defer os.Chdir(oldWd)

		os.Unsetenv("SKEETER_DIR")

		_, err := Dir("")
		if err == nil {
			t.Error("expected error when .skeeter not found")
		}
	})
}

func TestDirForInit(t *testing.T) {
	t.Run("flag takes precedence", func(t *testing.T) {
		result, err := DirForInit("/custom/path")
		if err != nil {
			t.Fatalf("DirForInit: %v", err)
		}
		if result != "/custom/path" {
			t.Errorf("DirForInit = %q, want %q", result, "/custom/path")
		}
	})

	t.Run("environment variable when no flag", func(t *testing.T) {
		os.Setenv("SKEETER_DIR", "/env/path")
		defer os.Unsetenv("SKEETER_DIR")

		result, err := DirForInit("")
		if err != nil {
			t.Fatalf("DirForInit: %v", err)
		}
		if result != "/env/path" {
			t.Errorf("DirForInit = %q, want %q", result, "/env/path")
		}
	})

	t.Run("cwd/.skeeter when no flag or env", func(t *testing.T) {
		tmpDir := t.TempDir()

		oldWd, _ := os.Getwd()
		if err := os.Chdir(tmpDir); err != nil {
			t.Fatal(err)
		}
		defer os.Chdir(oldWd)

		os.Unsetenv("SKEETER_DIR")

		result, err := DirForInit("")
		if err != nil {
			t.Fatalf("DirForInit: %v", err)
		}

		expected := filepath.Join(tmpDir, ".skeeter")
		if result != expected {
			t.Errorf("DirForInit = %q, want %q", result, expected)
		}
	})
}

func TestFindRoot(t *testing.T) {
	t.Run("finds in current directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		skeeterDir := filepath.Join(tmpDir, ".skeeter")
		if err := os.Mkdir(skeeterDir, 0755); err != nil {
			t.Fatal(err)
		}

		result, err := findRoot(tmpDir)
		if err != nil {
			t.Fatalf("findRoot: %v", err)
		}
		if result != skeeterDir {
			t.Errorf("findRoot = %q, want %q", result, skeeterDir)
		}
	})

	t.Run("finds in parent directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		skeeterDir := filepath.Join(tmpDir, ".skeeter")
		subDir := filepath.Join(tmpDir, "sub", "dir")
		if err := os.MkdirAll(skeeterDir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.MkdirAll(subDir, 0755); err != nil {
			t.Fatal(err)
		}

		result, err := findRoot(subDir)
		if err != nil {
			t.Fatalf("findRoot: %v", err)
		}
		if result != skeeterDir {
			t.Errorf("findRoot = %q, want %q", result, skeeterDir)
		}
	})

	t.Run("returns error when not found", func(t *testing.T) {
		tmpDir := t.TempDir()

		_, err := findRoot(tmpDir)
		if err == nil {
			t.Error("expected error when .skeeter not found")
		}
	})

	t.Run("ignores .skeeter file (not directory)", func(t *testing.T) {
		tmpDir := t.TempDir()
		skeeterFile := filepath.Join(tmpDir, ".skeeter")
		if err := os.WriteFile(skeeterFile, []byte("test"), 0644); err != nil {
			t.Fatal(err)
		}

		_, err := findRoot(tmpDir)
		if err == nil {
			t.Error("expected error when .skeeter is a file, not directory")
		}
	})
}
