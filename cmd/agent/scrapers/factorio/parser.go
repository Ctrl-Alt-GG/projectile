package factorio

import (
	"fmt"
	"strconv"
	"strings"
)

// This one was partially written by ChatGPT

func parseRCONPlayersList(output string) ([]string, error) {
	lines := strings.Split(output, "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("empty /players output")
	}

	header := strings.TrimSpace(lines[0])
	if !strings.HasPrefix(header, "Players (") || !strings.HasSuffix(header, "):") {
		return nil, fmt.Errorf("unexpected /players header: %q", header)
	}

	countStr := strings.TrimSuffix(
		strings.TrimPrefix(header, "Players ("),
		"):",
	)
	total, err := strconv.Atoi(countStr)
	if err != nil {
		return nil, fmt.Errorf("invalid player count in header: %q %w", countStr, err)
	}

	var online []string

	numPly := 0
	for _, line := range lines[1:] {
		if !strings.HasPrefix(line, "  ") {
			continue
		}

		name := strings.TrimSpace(line)

		if name != "" {
			numPly++
			if strings.HasSuffix(name, " (online)") {
				name = strings.TrimSuffix(name, " (online)")
				name = strings.TrimSpace(name)
				if name != "" {
					online = append(online, name)
				}
			}
		}
	}

	if numPly != total {
		return nil, fmt.Errorf("unexpected number of players returned %d != %d", total, numPly)
	}

	return online, nil
}
