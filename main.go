package main

import (
	"systemd-cd/application/flag_with_env"
	"systemd-cd/application/logrus"
	"systemd-cd/domain/model/git"
	"systemd-cd/domain/model/logger"
	"systemd-cd/infrastructure/externalapi/git_command"
)

// コマンドライン引数 / 環境変数
var (
	logLevel = flag_with_env.Uint("log-level", "LOG_LEVEL", 3, "Log level (0: Panic, 1: Fatal, 2: Error, 3: Warn, 4; Info, 5: Debug, 6: Trace)")
)

func main() {
	// コマンドライン引数 / 環境変数 の取得
	flag_with_env.Parse()

	// ロガーのインスタンス化, ログレベルを指定
	l := logrus.New()
	l.SetLevel(logger.Level(*logLevel))

	g := git.New(git_command.New())
	gitUser := "tingtt"
	gitAccessToken := ""
	repo, err := g.NewLocalRepository("./sample", git.RepositoryRemote{
		RemoteUrl: "https://github.com/tingtt/systemd-cd.git",
		User:      &gitUser,
		Token:     &gitAccessToken,
	}, "main")
	if err != nil {
		l.Error(err)
		return
	}
	exists, err := repo.DiffExists(true)
	if err != nil {
		l.Error(err)
		return
	}
	l.Debugf("diff exists: %v\n", exists)
}
