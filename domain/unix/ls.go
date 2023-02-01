package unix

import (
	"strings"
)

type LsOption struct {
	ReverceOrder         bool
	SortByDescendingTime bool
	DirTrailiingSlash    bool
}

func Ls(o ExecuteOption, o1 LsOption, target string) ([]string, error) {
	options := []string{}
	if o1.ReverceOrder {
		options = append(options, "-r")
	}
	if o1.SortByDescendingTime {
		options = append(options, "-t")
	}
	if o1.DirTrailiingSlash {
		options = append(options, "-p")
	}
	_, stdout, _, err := Execute(o, "ls", append(options, target)...)
	if err != nil {
		return nil, err
	}

	res := strings.Split(stdout.String(), "\n")
	return res, nil
}
