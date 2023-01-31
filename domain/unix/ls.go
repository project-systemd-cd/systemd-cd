package unix

import (
	"strings"
	"systemd-cd/domain/logger"
)

type LsOption struct {
	ReverceOrder         bool
	SortByDescendingTime bool
	DirTrailiingSlash    bool
}

func Ls(o ExecuteOption, o1 LsOption, target string) ([]string, error) {
	logger.Logger().Trace(logger.Var2Text("Called", []logger.Var{{Value: o}, {Value: o1}, {Name: "target", Value: target}}))

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
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return nil, err
	}

	res := strings.Split(stdout.String(), "\n")
	logger.Logger().Trace(logger.Var2Text("Finished", []logger.Var{{Name: "[]string", Value: res}}))
	return res, nil
}
