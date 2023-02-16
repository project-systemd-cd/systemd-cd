package echo

import (
	"strconv"
	"time"
)

func parseTime(s string) (time.Time, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err == nil {
		return time.Unix(i, 0), nil
	}

	t, err := time.Parse(time.RFC3339, s)
	if err == nil {
		return t, nil
	}

	return time.Time{}, err
}
