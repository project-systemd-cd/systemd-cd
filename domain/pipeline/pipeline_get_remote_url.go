package pipeline

import (
	"regexp"
	"strings"
)

func (p *pipeline) GetGitRemoteUrl() string {
	rep := regexp.MustCompile(`https?:\/\/[^\/:@]*:[^\/:@]*@`)
	if rep.MatchString(p.RepositoryLocal.RemoteUrl) {
		// Mask access token
		return rep.ReplaceAllString(p.RepositoryLocal.RemoteUrl, p.RepositoryLocal.RemoteUrl[:strings.Index(p.RepositoryLocal.RemoteUrl, "://")+3])
	}
	return p.RepositoryLocal.RemoteUrl
}
