package unix

import "systemd-cd/domain/model/logger"

type CpOption struct {
	Recursive bool
	Parents   bool
	Force     bool
}

func Cp(o ExecuteOption, o1 CpOption, src string, target string) error {
	logger.Logger().Tracef("Called:\n\toption: %+v\n\toption:%+v\n\tsrc: %v\n\ttarget: %v", o, o1, src, target)

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
		logger.Logger().Errorf("Error:\n\terr: %v", err)
		return err
	}

	logger.Logger().Trace("Finished")
	return nil
}
