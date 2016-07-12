// +build ignore

package main

import (
	"go/build"
	"log"
	"os"
	"path/filepath"
	"strings"

	gps "github.com/sdboyer/vsolver"
)

// This is probably the simplest possible implementation of gps. It does the
// substantive work that `go get` does, except it drops the resulting tree into
// vendor/, and prefers semver tags (if available) over branches.
func main() {
	// Operate on the current directory
	root, _ := os.Getwd()
	// Assume the current directory is correctly placed on a GOPATH, and derive
	// the ProjectRoot from it
	srcprefix := filepath.Join(build.Default.GOPATH, "src") + string(filepath.Separator)
	importroot := filepath.ToSlash(strings.TrimPrefix(root, srcprefix))

	// Set up params, including tracing
	params := gps.SolveParameters{
		RootDir:     root,
		ImportRoot:  gps.ProjectRoot(importroot),
		Trace:       true,
		TraceLogger: log.New(os.Stdout, "", 0),
	}

	// Set up a SourceManager with the NaiveAnalyzer
	sourcemgr, _ := gps.NewSourceManager(NaiveAnalyzer{}, ".repocache", false)
	defer sourcemgr.Release()

	// Prep and run the solver
	solver, _ := gps.Prepare(params, sourcemgr)
	solution, err := solver.Solve()
	if err == nil {
		// If no failure, blow away the vendor dir and write a new one out,
		// stripping nested vendor directories as we go.
		os.RemoveAll(filepath.Join(root, "vendor"))
		gps.CreateVendorTree(filepath.Join(root, "vendor"), solution, sourcemgr, true)
	}
}

type NaiveAnalyzer struct{}

func (a NaiveAnalyzer) GetInfo(path string, n gps.ProjectRoot) (gps.Manifest, gps.Lock, error) {
	return nil, nil, nil
}
