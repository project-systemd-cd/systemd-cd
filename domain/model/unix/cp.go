package unix

import "systemd-cd/domain/model/logger"

type CpOption struct {
	Recursive bool
	Parents   bool
	Force     bool
}

func Cp(o ExecuteOption, o1 CpOption, src string, target string) error {
	logger.Logger().Trace(logger.Var2Text("Called", []logger.Var{{Value: o}, {Value: o1}, {Name: "src", Value: src}, {Name: "target", Value: target}}))

	options := []string{}
	if o1.Recursive {
		options = append(options, "-R")
	}
	if o1.Parents {
		options = append(options, "-P")
	}
	if o1.Force {
		options = append(options, "-f")
	}
	_, _, _, err := Execute(o, "cp", append(options, src, target)...)
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return err
	}

	logger.Logger().Trace("Finished")
	return nil
}
