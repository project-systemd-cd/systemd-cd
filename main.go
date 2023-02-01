package main

import (
	"os"
	"systemd-cd/domain/git"
	"systemd-cd/domain/logger"
	"systemd-cd/domain/logrus"
	"systemd-cd/domain/pipeline"
	"systemd-cd/domain/runner"
	"systemd-cd/domain/systemd"
	"systemd-cd/infrastructure/datasource/toml"
	"systemd-cd/infrastructure/externalapi/git_command"
	"systemd-cd/infrastructure/externalapi/systemctl"
	"time"

	logruss "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

// flags
var (
	logLevel                  = pflag.String("log.level", "info", "Only log messages with the given severity or above. One of: [panic, fatal, error, warn, info, debug, trace]")
	logReportCaller           = pflag.Bool("log.report-caller", false, "Enable log report caller")
	logTimestamp              = pflag.Bool("log.timestamp", false, "Enable log timestamp.")
	varDir                    = pflag.String("dir.var", "/var/lib/systemd-cd/", "Path to variable files")
	srcDestDir                = pflag.String("dir.src", "/usr/local/systemd-cd/src/", "Path to service source files")
	binaryDestDir             = pflag.String("dir.binary", "/usr/local/systemd-cd/bin/", "Path to service binary files")
	etcDestDir                = pflag.String("dir.etc", "/usr/local/systemd-cd/etc/", "Path to service etc files")
	optDestDir                = pflag.String("dir.opt", "/usr/local/systemd-cd/opt/", "Path to service opt files")
	systemdUnitFileDestDir    = pflag.String("dir.systemd-unit-file", "/usr/local/lib/systemd/system/", "Path to systemd unit files.")
	systemdUnitEnvFileDestDir = pflag.String("dir.systemd-unit-env-file", "/usr/local/systemd-cd/etc/default/", "Path to systemd env files")
	backupDestDir             = pflag.String("dir.backup", "/var/backups/systemd-cd/", "Path to service backup files")

	manifestPaths    = pflag.StringSliceP("file.manifest", "f", nil, "Manifeset file path")
	pipelineInterval = pflag.Uint32("pipeline.interval", 180, "Interval of repository polling (second)")
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
	// parse flags
	pflag.Parse()

	logger.Init(logrus.New(logrus.Param{
		ReportCaller: logReportCaller,
		Formatter: &logruss.TextFormatter{
			FullTimestamp:   *logTimestamp,
			TimestampFormat: time.RFC3339Nano,
		},
	}))

	// `--log.level`
	ok, lv := convertLogLevel(*logLevel)
	if !ok {
		logger.Logger().Fatalf("`--log.level` must be specified as \"panic\", \"fatal\", \"error\", \"warn\", \"info\", \"debug\" or \"trace\"")
	}
	logger.Logger().SetLevel(lv)

	logger.Logger().Trace("-----------------------------------------------------------")
	logger.Logger().Trace("Analize cli options")
	logger.Logger().Trace("-----------------------------------------------------------")
	logger.Logger().Debugf("Log level is %s", *logLevel)
	logger.Logger().Tracef("< --log.level = %v", *logLevel)
	logger.Logger().Tracef("< --log.report-caller = %v", *logReportCaller)
	logger.Logger().Tracef("< --log.timestamp = %v", *logTimestamp)
	logger.Logger().Tracef("< --dir.var = %v", *varDir)
	logger.Logger().Tracef("< --dir.src = %v", *srcDestDir)
	logger.Logger().Tracef("< --dir.binary = %v", *binaryDestDir)
	logger.Logger().Tracef("< --dir.etc = %v", *etcDestDir)
	logger.Logger().Tracef("< --dir.opt = %v", *optDestDir)
	logger.Logger().Tracef("< --dir.systemd-unit-file = %v", *systemdUnitFileDestDir)
	logger.Logger().Tracef("< --dir.systemd-unit-env-file = %v", *systemdUnitEnvFileDestDir)
	logger.Logger().Tracef("< --dir.backup = %v", *backupDestDir)
	logger.Logger().Trace("-----------------------------------------------------------")

	s, err := systemd.New(systemctl.New(), *systemdUnitFileDestDir)
	if err != nil {
		logger.Logger().Fatal(err)
		os.Exit(1)
	}

	g := git.NewService(git_command.New())

	repo, err := toml.NewRepositoryPipeline(*varDir)
	if err != nil {
		logger.Logger().Fatal(err)
		os.Exit(1)
	}

	p, err := pipeline.NewService(
		repo, g, s,
		pipeline.Directories{
			Src:                *srcDestDir,
			Binary:             *binaryDestDir,
			Etc:                *etcDestDir,
			Opt:                *optDestDir,
			SystemdUnitFile:    *systemdUnitFileDestDir,
			SystemdUnitEnvFile: *systemdUnitEnvFileDestDir,
			Backup:             *backupDestDir,
		},
	)
	if err != nil {
		logger.Logger().Fatal(err)
		os.Exit(1)
	}

	runner, err := runner.NewService(p, runner.Option{PollingInterval: time.Duration(*pipelineInterval) * time.Second})
	if err != nil {
		logger.Logger().Fatal(err)
		os.Exit(1)
	}

	manifests := []pipeline.ServiceManifestLocal{}

	for _, path := range *manifestPaths {
		sml, err := loadManifest(path)
		if err != nil {
			logger.Logger().Fatal(err)
			os.Exit(1)
		}
		manifests = append(manifests, sml)
	}

	err = runner.Start(&manifests)
	if err != nil {
		logger.Logger().Fatal(err)
		os.Exit(1)
	}
}
