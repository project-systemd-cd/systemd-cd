package main

import (
	"fmt"
	"os"
	"systemd-cd/application/flag_with_env"
	"systemd-cd/application/systemd"
	"systemd-cd/domain/model/logger"
	"systemd-cd/domain/model/logrus"
	"systemd-cd/infrastructure/externalapi/systemctl"
	"time"

	logruss "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

// CLI args / ENV variables
var (
	logLevel                  = flag_with_env.String("log.level", "LOG_LEVEL", "info", "Only log messages with the given severity or above. One of: [panic, fatal, error, warn, info, debug, trace]")
	varDir                    = flag_with_env.String("var-dir", "VAR_DIR", "/var/lib/systemd-cd/", "")
	srcDestDir                = flag_with_env.String("src-dest-dir", "SRC_DEST_DIR", "/usr/local/systemd-cd/src/", "")
	binaryDestDir             = flag_with_env.String("binary-dest-dir", "BINARY_DEST_DIR", "/usr/local/systemd-cd/bin/", "")
	etcDestDir                = flag_with_env.String("etc-dest-dir", "ETC_DEST_DIR", "/usr/local/systemd-cd/etc/", "")
	optDestDir                = flag_with_env.String("opt-dest-dir", "OPT_DEST_DIR", "/usr/local/systemd-cd/opt/", "")
	systemdUnitFileDestDir    = flag_with_env.String("systemd-unit-file-dest-dir", "SYSTEMD_UNIT_FILE_DEST_DIR", "/usr/local/lib/systemd/system/", "")
	systemdUnitEnvFileDestDir = flag_with_env.String("systemd-unit-env-file-dest-dir", "SYSTEMD_UNIT_ENV_FILE_DEST_DIR", "/usr/local/systemd-cd/etc/default/", "")
	backupDestDir             = flag_with_env.String("backup-dest-dir", "BACKUP_DEST_DIR", "/var/backups/systemd-cd/", "")
)

func convertLogLevel(str string) (ok bool, lv logger.Level) {
	switch str {
	case "panic":
		return true, logger.PanicLevel
	case "fatal":
		return true, logger.FatalLevel
	case "error":
		return true, logger.ErrorLevel
	case "warn":
		return true, logger.WarnLevel
	case "info":
		return true, logger.InfoLevel
	case "debug":
		return true, logger.DebugLevel
	case "trace":
		return true, logger.TraceLevel
	default:
		return false, logger.InfoLevel
	}
}

func main() {
	logger.Init(logrus.New(logrus.Param{
		RepeatCaller: func() *bool { var b = true; return &b }(),
		Formatter: &logruss.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339Nano,
		},
	}))

	// parse flags
	pflag.Parse()

	// `--log.level`
	ok, lv := convertLogLevel(*logLevel)
	if !ok {
		logger.Logger().Fatalf("`--log.level` must be specified as \"panic\", \"fatal\", \"error\", \"warn\", \"info\", \"debug\" or \"trace\"")
		os.Exit(1)
	}
	logger.Logger().SetLevel(lv)

	i, err := systemd.New(systemctl.New(), *systemdUnitFileDestDir)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		os.Exit(1)
	}
	envFile := *systemdUnitEnvFileDestDir + "system-cd-go"
	us, err := i.NewService(
		"systemd-cd-go",
		systemd.UnitFileService{
			Unit: systemd.UnitDirective{
				Description:   "gitops agent for systemd-based linux",
				Documentation: "https://github.com/tingtt/systmed-cd",
				After:         []string{"syslog.target", "network.target"},
				Requires:      nil,
				Wants:         nil,
				Conflicts:     []string{"sendmail.servic", "exim.service"},
			},
			Service: systemd.ServiceDirective{
				Type:            &systemd.UnitTypeSimple,
				EnvironmentFile: &envFile,
				ExecStart:       "watch tail /var/log/syslog",
				ExecStop:        nil,
				ExecReload:      nil,
				Restart:         nil,
				RemainAfterExit: nil,
			},
			Install: systemd.InstallDirective{
				Alias:           nil,
				RequiredBy:      nil,
				WantedBy:        []string{"multi-user.target"},
				Also:            nil,
				DefaultInstance: nil,
			},
		},
		map[string]string{},
	)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		os.Exit(1)
	}
	s, err := us.GetStatus()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("status: %v\n", s)
}
