package unix

import "systemd-cd/domain/logger"

type MvOption struct {
	Force bool
}

func Mv(o ExecuteOption, o1 MvOption, src string, target string) error {
	logger.Logger().Trace(logger.Var2Text("Called", []logger.Var{{Value: o}, {Value: o1}, {Name: "src", Value: src}, {Name: "target", Value: target}}))

	options := []string{}
	if o1.Force {
		options = append(options, "-f")
	}
	_, _, _, err := Execute(o, "mv", append(options, src, target)...)
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return err
	}

	logger.Logger().Trace("Finished")
	return nil
}
