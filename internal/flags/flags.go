package flags

import "github.com/coderaiser/go-subscriber/internal/flags/meta"

// Result holds the outcome of flag parsing.
// ExitCode == -1 means no flags matched — continue normal startup.
// ExitCode >= 0 means print Output and exit with ExitCode.
type Result struct {
	ExitCode int
	Output   string
}

// Parse inspects args (typically os.Args[1:]) and returns a Result.
// No I/O is performed — all output is in Result.Output.
func Parse(args []string) Result {
	for _, arg := range args {
		switch arg {
		case "-v", "--version":
			return Result{ExitCode: 0, Output: meta.VersionLine() + "\n"}
		case "-h", "--help":
			return Result{ExitCode: 0, Output: meta.Help()}
		default:
			return Result{ExitCode: 1, Output: "unknown flag: " + arg + "\n"}
		}
	}
	return Result{ExitCode: -1}
}
