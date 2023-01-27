package unix

import (
	"strings"
	"systemd-cd/domain/model/logger"
)

type LsOption struct {
	ReverceOrder         bool
	SortByDescendingTime bool
}

func Ls(o ExecuteOption, o1 LsOption, target string) ([]string, error) {
	logger.Logger().Tracef("Called:\n\toption: %+v\n\toption: %+v\n\ttarget: %v", o, o1, target)

	options := []string{}
	if o1.ReverceOrder {
		options = append(options, "-r")
	}
	if o1.SortByDescendingTime {
		options = append(options, "-t")
	}
	_, stdout, _, err := Execute(o, "ls", append(options, target)...)
	if err != nil {
		logger.Logger().Errorf("Error:\n\terr: %v", err)
		return nil, err
	}

	res := strings.Split(stdout.String(), "\n")
	logger.Logger().Tracef("Finished:\n\tstdout: %v", res)
	return res, nil
}
