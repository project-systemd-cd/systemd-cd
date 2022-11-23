package logrus

import (
	"systemd-cd/domain/model/logger"

	"github.com/sirupsen/logrus"
)

func New() logger.LoggerI {
	logrus := logrus.New()
	return &WrapLogrus{Logger: *logrus}
}

type WrapLogrus struct {
	logrus.Logger
}

func (w *WrapLogrus) SetLevel(level logger.Level) error {
	// ログレベルを検証
	_, err := logrus.Level(level).MarshalText()
	if err != nil {
		// 不正なログレベル
		w.Error(err)
	} else {
		// ログレベルを適用
		w.Logger.SetLevel(logrus.Level(level))
		w.Infof("Log level: %q", w.Level.String())
	}

	return err
}
