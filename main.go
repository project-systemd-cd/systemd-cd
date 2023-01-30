package main

import (
	"fmt"
	"os"
	"systemd-cd/domain/git"
	"systemd-cd/domain/logger"
	"systemd-cd/domain/logrus"
	"systemd-cd/domain/pipeline"
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

	s, err := systemd.New(systemctl.New(), *systemdUnitFileDestDir)
	if err != nil {
		logger.Logger().Fatalf("Failed:\n\terr: %v", err)
		os.Exit(1)
	}

	g := git.NewService(git_command.New())

	repo, err := toml.NewRepositoryPipeline(*varDir)
	if err != nil {
		logger.Logger().Fatalf("Failed:\n\terr: %v", err)
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
		logger.Logger().Fatalf("Failed:\n\terr: %v", err)
		os.Exit(1)
	}

	p1, err := p.NewPipeline(pipeline.ServiceManifestLocal{
		GitRemoteUrl:    "https://github.com/tingtt/prometheus_sh_exporter.git",
		GitTargetBranch: "main",
		GitTagRegex:     func() *string { s := "v*"; return &s }(),
		GitManifestFile: nil,
		Name:            "prometheus_sh_exporter",
		TestCommands:    nil,
		BuildCommands:   func() *[]string { s := []string{"/usr/bin/go build"}; return &s }(),
		Opt:             &[]string{},
		Binaries:        func() *[]string { s := []string{"prometheus_sh_exporter"}; return &s }(),
		SystemdOptions: []pipeline.SystemdOption{{
			Name:           "prometheus_sh_exporter",
			Description:    func() *string { s := "The shell exporter allows probing with shell scripts."; return &s }(),
			ExecuteCommand: "prometheus_sh_exporter",
			Args:           "",
			EnvVars:        []pipeline.EnvVar{},
			Etc: []pipeline.PathOption{{
				Target: "sh.yml",
				Option: "-config.file",
			}},
			Port: func() *uint16 { p := uint16(9923); return &p }(),
		}, {
			Name:           "prometheus_sh_exporter2",
			Description:    func() *string { s := "The shell exporter allows probing with shell scripts."; return &s }(),
			ExecuteCommand: "prometheus_sh_exporter",
			Args:           "--port 9924",
			EnvVars:        []pipeline.EnvVar{},
			Etc: []pipeline.PathOption{{
				Target: "sh.yml",
				Option: "-config.file",
			}},
			Port: func() *uint16 { p := uint16(9924); return &p }(),
		}},
	})
	if err != nil {
		logger.Logger().Fatalf("Failed:\n\terr: %v", err)
		os.Exit(1)
	}

	fmt.Printf("p1.GetStatus(): %v\n", p1.GetStatus())
	fmt.Printf("p1.GetCommitRef(): %v\n", p1.GetCommitRef())
}
