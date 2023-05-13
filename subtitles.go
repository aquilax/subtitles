package subtitles

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/aquilax/with"
)

type SubRipSubtitle struct {
	id    int
	start time.Time
	end   time.Time
	text  string
}

type SubtitleCallback func(s SubRipSubtitle) (bool, error)

func ParseSubRip(r io.Reader, cb SubtitleCallback) error {
	scanner := bufio.NewScanner(r)

	lineNum := 0
	var err error
	var id int
	var line, start, end, text string
	var startTime, endTime time.Time
	for scanner.Scan() {
		line = scanner.Text()
		lineNum++
		if line == "" && id != 0 && start != "" && end != "" && text != "" {
			if startTime, err = parseSubRipTime(start); err != nil {
				return fmt.Errorf("failed to parse start time: %v", err)
			}
			if endTime, err = parseSubRipTime(end); err != nil {
				return fmt.Errorf("failed to parse end time: %v", err)
			}
			if stop, err := cb(SubRipSubtitle{
				id:    id,
				start: startTime,
				end:   endTime,
				text:  text,
			}); stop || err != nil {
				return err
			}
			id = 0
			start = ""
			end = ""
			text = ""
			continue
		}
		if id == 0 {
			id, err = strconv.Atoi(line)
			if err != nil {
				return fmt.Errorf("failed to parse subtitle number: %v", err)
			}
			continue
		}

		if start == "" || end == "" {
			timings := strings.Split(line, " --> ")
			if len(timings) != 2 {
				return fmt.Errorf("invalid subtitle timing: %s", line)
			}
			start = timings[0]
			end = timings[1]
			continue
		}

		if text == "" {
			text = line
		} else {
			text += "\n" + line
		}
	}

	if id != 0 && start != "" && end != "" && text != "" {
		if startTime, err = parseSubRipTime(start); err != nil {
			return fmt.Errorf("failed to parse start time: %v", err)
		}
		if endTime, err = parseSubRipTime(end); err != nil {
			return fmt.Errorf("failed to parse end time: %v", err)
		}
		if stop, err := cb(SubRipSubtitle{
			id:    id,
			start: startTime,
			end:   endTime,
			text:  text,
		}); stop || err != nil {
			return err
		}
	}

	if err = scanner.Err(); err != nil {
		return err
	}
	return nil
}

func parseSubRipTime(t string) (time.Time, error) {
	var hour, minute, second, millisecond int
	_, err := fmt.Sscanf(t, "%d:%d:%d,%d", &hour, &minute, &second, &millisecond)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse timestamp: %v", err)
	}

	return time.Date(0, 0, 0, hour, minute, second, millisecond*int(time.Millisecond), time.UTC), nil
}

func formatSubRipTime(t time.Time) string {
	return t.Format("15:04:05,000")
}

func Encode(w io.Writer, s SubRipSubtitle) error {
	return with.Errors(
		with.GetSecond(func() (any, error) { return fmt.Fprintln(w, s.id) }),
		with.GetSecond(func() (any, error) {
			return fmt.Fprintf(w, "%s --> %s\n", formatSubRipTime(s.start), formatSubRipTime(s.end))
		}),
		with.GetSecond(func() (any, error) { return fmt.Fprintln(w, s.text) }),
		with.GetSecond(func() (any, error) { return fmt.Fprintln(w, "") }),
	)
}
