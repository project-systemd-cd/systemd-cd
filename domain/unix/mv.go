package unix

type MvOption struct {
	Force bool
}

func Mv(o ExecuteOption, o1 MvOption, src string, target string) error {
	options := []string{}
	if o1.Force {
		options = append(options, "-f")
	}
	_, _, _, err := Execute(o, "mv", append(options, src, target)...)
	if err != nil {
		return err
	}

	return nil
}
