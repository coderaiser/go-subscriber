package formatter

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

const (
	barWidth    = 40
	barComplete = '█'
	barEmpty    = '░'
	okEmoji     = "👌"
	failEmoji   = "❌"
	skipEmoji   = "⚠️"
	okMark      = "✅"
	YELLOW      = "\033[33m" 
)

// ProgressBarFormatter outputs a progress bar to stderr and final output to stdout.
type ProgressBarFormatter struct {
	total    int
	color    string
	stackEnv string
	out      strings.Builder // stdout output buffered
}

func NewProgressBar(total int) *ProgressBarFormatter {
	color := os.Getenv("TAPE_PROGRESS_BAR_COLOR")
	if color == "" {
		color = YELLOW
	}
	return &ProgressBarFormatter{
		total:    total,
		color:    color,
		stackEnv: os.Getenv("TAPE_PROGRESS_BAR_STACK"),
	}
}

func (f *ProgressBarFormatter) Start(total int) string {
	return ""
}

func (f *ProgressBarFormatter) Test(name string) string { return "" }

func (f *ProgressBarFormatter) TestEnd(count, total, failed int, name string) string {
	failStr := okEmoji
	if failed > 0 {
		failStr = fmt.Sprintf("\033[31m%d\033[0m", failed) // red
	}
	bar := RenderBar(count, total, f.color)
	pct := 0
	if total > 0 {
		pct = count * 100 / total
		if pct > 100 {
			pct = 100
		}
	}
	displayTotal := total
	if count > displayTotal {
		displayTotal = count
	}
	truncName := Truncate(name, 40)
	line := fmt.Sprintf("%s %d%% | %s | %d/%d | %s", bar, pct, failStr, count, displayTotal, truncName)
	width := termWidth()
	if visibleLen(line) > width {
		line = truncateANSI(line, width)
	}
	fmt.Fprintf(os.Stderr, "\r%s", line)
	return ""
}

func (f *ProgressBarFormatter) Success(count int, message string) string { return "" }

func (f *ProgressBarFormatter) Fail(count int, message, operator string, result, expected any, output, at, errorStack string) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "\n# %s\n", message) // will be flushed after bar clears
	fmt.Fprintf(&sb, "%s not ok %d %s\n", failEmoji, count, message)
	sb.WriteString("  ---\n")
	if operator != "" {
		fmt.Fprintf(&sb, "    operator: %s\n", operator)
	}
	if output != "" && operator == "" {
		sb.WriteString(output)
	} else {
		fmt.Fprintf(&sb, "    expected: |-\n      %v\n", expected)
		fmt.Fprintf(&sb, "    result: |-\n      %v\n", result)
	}
	if at != "" {
		fmt.Fprintf(&sb, "    %s\n", at)
	}
	if f.stackEnv != "0" && errorStack != "" {
		fmt.Fprintf(&sb, "    stack: |-\n%s\n", errorStack)
	}
	sb.WriteString("  ...\n\n")
	f.out.WriteString(sb.String())
	return ""
}

func (f *ProgressBarFormatter) Comment(message string) string {
	return fmt.Sprintf("# %s\n", message)
}

func (f *ProgressBarFormatter) End(count, passed, failed, skipped int) string {
	fmt.Fprintln(os.Stderr) // clear progress line
	var sb strings.Builder
	// flush buffered fail output
	sb.WriteString(f.out.String())
	sb.WriteString("\n")
	fmt.Fprintf(&sb, "1..%d\n", count)
	fmt.Fprintf(&sb, "# tests %d\n", count)
	fmt.Fprintf(&sb, "# pass %d\n", passed)
	if skipped > 0 {
		fmt.Fprintf(&sb, "# %s skip %d\n", skipEmoji, skipped)
	}
	sb.WriteString("\n")
	if failed > 0 {
		fmt.Fprintf(&sb, "# %s fail %d\n", failEmoji, failed)
	} else {
		fmt.Fprintf(&sb, "# %s ok\n", okMark)
	}
	sb.WriteString("\n\n")
	return sb.String()
}

// ansiEscape matches ANSI escape sequences.
var ansiEscape = regexp.MustCompile("\x1b\\[[0-9;]*m")

// visibleLen returns the visible length of a string ignoring ANSI escape codes.
func visibleLen(s string) int {
	return len([]rune(ansiEscape.ReplaceAllString(s, "")))
}

// truncateANSI truncates s to n visible runes, preserving ANSI codes.
func truncateANSI(s string, n int) string {
	visible := 0
	var out []byte
	i := 0
	runes := []byte(s)
	for i < len(runes) {
		if runes[i] == 0x1b && i+1 < len(runes) && runes[i+1] == '[' {
			// consume escape sequence
			j := i + 2
			for j < len(runes) && runes[j] != 'm' {
				j++
			}
			if j < len(runes) {
				j++
			}
			out = append(out, runes[i:j]...)
			i = j
			continue
		}
		if visible >= n {
			break
		}
		out = append(out, runes[i])
		visible++
		i++
	}
	return string(out)
}

// RenderBar renders a progress bar string.
func RenderBar(done, total int, color string) string {
	if total == 0 {
		return fmt.Sprintf("%s%s\033[0m", color, strings.Repeat(string(barEmpty), barWidth))
	}
	filled := done * barWidth / total
	if filled > barWidth {
		filled = barWidth
	}
	bar := strings.Repeat(string(barComplete), filled) + strings.Repeat(string(barEmpty), barWidth-filled)
	return fmt.Sprintf("%s%s\033[0m", color, bar)
}

// Truncate truncates a string to n runes, adding "..." if shortened.
func Truncate(s string, n int) string {
	if len([]rune(s)) <= n {
		return s
	}
	return string([]rune(s)[:n-3]) + "..."
}

