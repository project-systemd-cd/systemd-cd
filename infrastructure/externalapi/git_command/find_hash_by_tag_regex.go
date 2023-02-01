package git_command

import (
	"regexp"
	"strings"
	"systemd-cd/domain/errors"
	"systemd-cd/domain/git"
	"systemd-cd/domain/unix"
)

func (*GitCommand) FindHashByTagRegex(workingDir git.Path, regex string) (hash string, err error) {
	r, err := open(workingDir)
	if err != nil {
		return "", err
	}

	var wd string = string(workingDir)
	_, b, _, err := unix.Execute(unix.ExecuteOption{WorkingDirectory: &wd}, "git", "tag")
	if err != nil {
		return "", err
	}
	found := false
	tags := strings.Split(b.String(), "\n")
	for i, j := 0, len(tags)-1; i < j; i, j = i+1, j-1 {
		tags[i], tags[j] = tags[j], tags[i]
	}
	for _, v := range tags {
		if v == "" {
			continue
		}
		matched, err := regexp.MatchString(regex, v)
		if err != nil {
			continue
		}
		if matched {
			r2, err := r.Tag(v)
			if err != nil {
				return "", err
			}
			hash = r2.Hash().String()
			found = true
			break
		}
	}

	if !found {
		err := &errors.ErrNotFound{Object: "Tag", IdName: "name", Id: regex}
		return "", err
	}

	return hash, nil
}
