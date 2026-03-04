package scanner

import (
	"os"
	"regexp"
	"strconv"
	"strings"
)

type LyricLine struct {
	Time float64 `json:"time"` // seconds
	Text string  `json:"text"`
}

// ParseLRC parses LRC format lyrics into structured lines.
func ParseLRC(content string) ([]LyricLine, bool) {
	timeTag := regexp.MustCompile(`\[(\d{1,2}):(\d{2})(?:\.(\d{1,3}))?\]`)
	lines := strings.Split(content, "\n")
	var result []LyricLine
	isSynced := false

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		matches := timeTag.FindAllStringSubmatchIndex(line, -1)
		if len(matches) == 0 {
			continue
		}

		isSynced = true

		var times []float64
		lastEnd := 0
		for _, m := range matches {
			minStr := line[m[2]:m[3]]
			secStr := line[m[4]:m[5]]

			min, _ := strconv.ParseFloat(minStr, 64)
			sec, _ := strconv.ParseFloat(secStr, 64)

			var ms float64
			if m[6] >= 0 && m[7] >= 0 {
				msStr := line[m[6]:m[7]]
				msVal, _ := strconv.ParseFloat(msStr, 64)
				switch len(msStr) {
				case 1:
					ms = msVal * 100
				case 2:
					ms = msVal * 10
				case 3:
					ms = msVal
				}
			}

			t := min*60 + sec + ms/1000
			times = append(times, t)
			lastEnd = m[1]
		}

		text := strings.TrimSpace(line[lastEnd:])

		for _, t := range times {
			result = append(result, LyricLine{Time: t, Text: text})
		}
	}

	// Sort by time
	for i := 1; i < len(result); i++ {
		for j := i; j > 0 && result[j].Time < result[j-1].Time; j-- {
			result[j], result[j-1] = result[j-1], result[j]
		}
	}

	if !isSynced {
		return nil, false
	}

	return result, true
}

// ReadLRCFile reads a .lrc file and returns the content.
func ReadLRCFile(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
