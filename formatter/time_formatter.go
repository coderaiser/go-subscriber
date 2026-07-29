package formatter

import (
	"fmt"
	"os"
	"time"
)

// TimeFormatter is progress bar with wall clock timer.
type TimeFormatter struct {
	ProgressBarFormatter
	startTime time.Time
	clock     string
}

func NewTime(total int) *TimeFormatter {
	clock := os.Getenv("TAPE_TIME_CLOCK")
	if clock == "" {
		clock = "⏳"
	}
	return &TimeFormatter{
		ProgressBarFormatter: *NewProgressBar(total),
		clock:                clock,
	}
}

func (f *TimeFormatter) Start(total int) string {
	f.startTime = time.Now()
	return f.ProgressBarFormatter.Start(total)
}

func (f *TimeFormatter) TestEnd(count, total, failed int, name string) string {
	// override to add time to bar line
	elapsed := time.Since(f.startTime)
	m := int(elapsed.Minutes())
	s := int(elapsed.Seconds()) % 60
	timeStr := fmt.Sprintf("%s %02d:%02d", f.clock, m, s)

	failStr := okEmoji
	if failed > 0 {
		failStr = fmt.Sprintf("\033[31m%d\033[0m", failed)
	}
	bar := RenderBar(count, total, f.color)
	pct := 0
	if total > 0 {
		pct = count * 100 / total
	}
	truncName := Truncate(name, 30)
	line := fmt.Sprintf("%s %d%% | %s | %d/%d | %s | %s",
		bar, pct, failStr, count, total, timeStr, truncName)
	fmt.Fprintf(os.Stderr, "\r%s", line)
	return ""
}
