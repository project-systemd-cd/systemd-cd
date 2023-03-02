package unix

type CpOption struct {
	Recursive bool
	Parents   bool
	Force     bool
}

func Cp(o ExecuteOption, o1 CpOption, src string, target string) error {
	options := []string{}
	if o1.Recursive {
		options = append(options, "-R")
	}
	if o1.Parents {
		options = append(options, "--parents")
	}
	if o1.Force {
		options = append(options, "-f")
	}
	_, _, _, err := Execute(o, "cp", append(options, src, target)...)
	if err != nil {
		return err
	}

	return nil
}
