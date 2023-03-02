package unix

type RmOption struct {
	Recursive bool
}

func Rm(o ExecuteOption, o1 RmOption, target ...string) error {
	args := []string{}
	if o1.Recursive {
		args = append(args, "-r")
	}
	for _, t := range target {
		args = append(args, t)
	}
	_, _, _, err := Execute(o, "rm", args...)
	if err != nil {
		return err
	}

	return nil
}
