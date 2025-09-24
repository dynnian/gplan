//go:build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// -----------------------------------------------------------------------------
// Config
// -----------------------------------------------------------------------------

var (
	BinaryName = "gplan"
	Package    = "github.com/dynnian/gplan"
	BuildDir   = "build"

	Platforms = []string{
		"windows/amd64",
		"windows/arm64",
		"darwin/amd64",
		"darwin/arm64",
		"linux/amd64",
		"linux/arm64",
	}

	ToolGolangCILint = "github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
	ToolGofumpt      = "mvdan.cc/gofumpt@latest"
	ToolGoimports    = "golang.org/x/tools/cmd/goimports@latest"
	ToolGolines      = "github.com/segmentio/golines@latest"
)

// Default target
var Default = All

// -----------------------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------------------

func gitCommit() string {
	out, _ := sh.Output("git", "rev-parse", "HEAD")
	return strings.TrimSpace(out)
}

func buildDateUTC() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05Z")
}

func ldflags(version, date, commit string) (string, string) {
	value := fmt.Sprintf(`-X %s/internal/version.Version=%s `+
		`-X %s/internal/version.BuildDate=%s `+
		`-X %s/internal/version.GitCommit=%s`,
		Package, version, Package, date, Package, commit)
	return "-ldflags", value
}

func goPathBin() (string, error) {
	out, err := sh.Output("go", "env", "GOPATH")
	if err != nil {
		return "", err
	}
	return filepath.Join(strings.TrimSpace(out), "bin"), nil
}

func hasCmd(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func getVersion() string {
	if v := os.Getenv("VERSION"); v != "" {
		return v
	}
	return "dev"
}

// -----------------------------------------------------------------------------
// Targets
// -----------------------------------------------------------------------------

// All: deps, lint, format, build
func All() {
	mg.SerialDeps(Deps, Lint, Format, Build)
}

// Build for current platform
func Build() error {
	version := getVersion()
	commit := gitCommit()
	date := buildDateUTC()

	flag, val := ldflags(version, date, commit)

	out := BinaryName
	if runtime.GOOS == "windows" {
		out += ".exe"
	}

	args := []string{"build", flag, val, "-o", out, "."}
	fmt.Println("Building", out, "for", runtime.GOOS, runtime.GOARCH)
	return sh.RunV("go", args...)
}

// Clean artifacts
func Clean() error {
	fmt.Println("Cleaning...")
	_ = sh.RunV("go", "clean")
	return os.RemoveAll(BuildDir)
}

// Deps: tidy + install tooling if missing
func Deps() error {
	fmt.Println("Checking and updating dependencies...")
	if err := sh.RunV("go", "mod", "tidy"); err != nil {
		return err
	}
	if out, _ := sh.Output("git", "status", "--porcelain", "go.mod", "go.sum"); strings.TrimSpace(out) == "" {
		fmt.Println("No missing dependencies. All modules are up to date.")
	} else {
		fmt.Println("Dependencies updated. Please review changes in go.mod and go.sum.")
	}
	if !hasCmd("golangci-lint") {
		fmt.Println("Installing golangci-lint...")
		if err := sh.RunV("go", "install", ToolGolangCILint); err != nil {
			return err
		}
	}
	if !hasCmd("gofumpt") {
		fmt.Println("Installing gofumpt...")
		if err := sh.RunV("go", "install", ToolGofumpt); err != nil {
			return err
		}
	}
	if !hasCmd("goimports") {
		fmt.Println("Installing goimports...")
		if err := sh.RunV("go", "install", ToolGoimports); err != nil {
			return err
		}
	}
	if !hasCmd("golines") {
		fmt.Println("Installing golines...")
		if err := sh.RunV("go", "install", ToolGolines); err != nil {
			return err
		}
	}
	return nil
}

// Lint with golangci-lint
func Lint() error {
	fmt.Println("Running linter...")
	return sh.RunV("golangci-lint", "run", "--fix", "-c", ".golangci.yml", "./...")
}

// Format with gofumpt, goimports, golines
func Format() error {
	fmt.Println("Formatting code...")
	if err := sh.RunV("gofumpt", "-l", "-w", "."); err != nil {
		return err
	}
	if err := sh.RunV("goimports", "-w", "."); err != nil {
		return err
	}
	if err := sh.RunV("golines", "-l", "-m", "120", "-t", "4", "-w", "."); err != nil {
		return err
	}
	fmt.Println("Code formatted.")
	return nil
}

// BuildAll for matrix
func BuildAll() error {
	version := getVersion()
	commit := gitCommit()
	date := buildDateUTC()

	if err := os.MkdirAll(BuildDir, 0o755); err != nil {
		return err
	}

	flag, val := ldflags(version, date, commit)

	for _, p := range Platforms {
		parts := strings.SplitN(p, "/", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid platform %q", p)
		}
		goos, goarch := parts[0], parts[1]
		ext := ""
		if goos == "windows" {
			ext = ".exe"
		}
		out := filepath.Join(BuildDir, fmt.Sprintf("%s-%s-%s%s", BinaryName, goos, goarch, ext))
		fmt.Printf("Building %s (%s/%s)\n", out, goos, goarch)

		env := map[string]string{
			"GOOS":        goos,
			"GOARCH":      goarch,
			"CGO_ENABLED": "0",
		}
		if err := sh.RunWithV(env, "go", "build", flag, val, "-o", out, "."); err != nil {
			return err
		}
	}
	return nil
}

// Version info
func Version() {
	fmt.Println("Git Commit:", gitCommit())
	fmt.Println("Build Date:", buildDateUTC())
}

// Install to GOPATH/bin
func Install() error {
	bin, err := goPathBin()
	if err != nil {
		return err
	}
	version := getVersion()
	commit := gitCommit()
	date := buildDateUTC()

	flag, val := ldflags(version, date, commit)

	out := filepath.Join(bin, BinaryName)
	if runtime.GOOS == "windows" {
		out += ".exe"
	}

	fmt.Println("Installing to", out)
	return sh.RunV("go", "build", flag, val, "-o", out, ".")
}

// Uninstall from GOPATH/bin
func Uninstall() error {
	bin, err := goPathBin()
	if err != nil {
		return err
	}
	target := filepath.Join(bin, BinaryName)
	fmt.Println("Removing", target)
	return os.Remove(target)
}

// Help (optional; mage -l also works)
func Help() {
	fmt.Println("Available targets:")
	fmt.Println("  mage                 - Run deps, lint, format, build")
	fmt.Println("  mage build           - Build for current platform")
	fmt.Println("  mage clean           - Remove build artifacts")
	fmt.Println("  mage deps            - go mod tidy + install tools")
	fmt.Println("  mage lint            - Run golangci-lint")
	fmt.Println("  mage format          - Run gofumpt, goimports, golines")
	fmt.Println("  mage buildAll        - Cross-compile for all platforms")
	fmt.Println("  mage version         - Print git commit and build date")
	fmt.Println("  mage install         - Install to GOPATH/bin")
	fmt.Println("  mage uninstall       - Remove from GOPATH/bin")
}
