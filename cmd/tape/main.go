package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	tapeast "github.com/coderaiser/go-subscriber/cmd/tape/ast"
	"github.com/coderaiser/go-subscriber/formatter"
	"github.com/coderaiser/go-subscriber/runner"
	"github.com/coderaiser/go-subscriber/state"
)

const version = "1.0.0"

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("go-subscriber", flag.ExitOnError)
	format            := flags.String("f", "", "output format: tap|progress-bar|short|fail|time|json-lines")
	help              := flags.Bool("h", false, "display this help and exit")
	ver               := flags.Bool("v", false, "output version information and exit")
	noCheckScopes     := flags.Bool("no-check-scopes", false, "do not check scope format")
	noCheckAssertions := flags.Bool("no-check-assertions-count", false, "do not check assertion count")
	noCheckDuplicates := flags.Bool("no-check-duplicates", false, "do not check for duplicates")
	flags.Parse(args)

	if *help {
		fmt.Fprint(stdout, usage)
		return 0
	}
	if *ver {
		fmt.Fprintln(stdout, version)
		return 0
	}
	if *noCheckScopes {
		os.Setenv("TAPE_CHECK_SCOPES", "0")
	}
	if *noCheckAssertions {
		os.Setenv("TAPE_CHECK_ASSERTIONS_COUNT", "0")
	}

	path := "./..."
	if flags.NArg() > 0 {
		path = flags.Arg(0)
	}
	dir := "."

	// check duplicates before running
	if !*noCheckDuplicates {
		dups, err := tapeast.FindDuplicates(dir)
		if err != nil {
			fmt.Fprintf(stderr, "tape: scan duplicates: %v\n", err)
			return 1
		}
		if len(dups) > 0 {
			fmt.Fprintf(stderr, "tape: duplicate test names found:\n")
			for _, d := range dups {
				fmt.Fprintf(stderr, "  %s\n", d)
			}
			return 1
		}
	}

	// count tests for progress bar total
	total, err := tapeast.CountTests(dir)
	if err != nil {
		fmt.Fprintf(stderr, "tape: count tests: %v\n", err)
		return 1
	}

	// find Only calls — restrict run if any found
	onlyCalls, err := tapeast.FindOnlyCalls(dir)
	if err != nil {
		fmt.Fprintf(stderr, "tape: scan Only: %v\n", err)
		return 1
	}

	goArgs := []string{"test", "-json", "-v", path}
	if pattern := tapeast.BuildRunPattern(onlyCalls); pattern != "" {
		fmt.Fprintf(stderr, "tape: Only mode — %s\n", pattern)
		goArgs = append(goArgs, "-run", pattern)
	}

	f := formatter.New(*format, stdout, total)
	store := state.New()
	r := runner.New(runner.NewOSExecutor())

	ch, err := r.Run(goArgs...)
	if err != nil {
		fmt.Fprintf(stderr, "tape: %v\n", err)
		return 1
	}

	for event := range ch {
		store.Apply(event)
		f.FromEvent(event)
	}

	passed, failed, skipped := store.Summary()
	f.End(len(passed), len(failed), len(skipped))

	if len(failed) > 0 {
		return 1
	}
	return 0
}

const usage = `Usage: go-subscriber [options] [path]

Options:
  -h                           display this help and exit
  -v                           output version information and exit
  -f format                    use a specific output format
                               default: progress-bar (tap on CI)
                               values: tap|progress-bar|short|fail|time|json-lines
  --no-check-scopes            do not check that messages contain scope: 'scope: message'
  --no-check-assertions-count  do not check that assertion count is no more than 1
  --no-check-duplicates        do not check messages for duplicates
`
