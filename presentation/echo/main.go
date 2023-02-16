package echo

import (
	"errors"
	"fmt"
	"systemd-cd/domain/logger"
	"systemd-cd/domain/runner"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	repository runner.IRepositoryInmemory
	jwtIssuer  *string
	jwtSecret  *string
	username   *string
	password   *string
)

type Args struct {
	Repository   runner.IRepositoryInmemory
	JwtIssuer    string
	JwtSecret    string
	Username     string
	Password     string
	AllowOrigins []string
}

func (args Args) validate() error {
	if args.Repository == nil {
		return errors.New("Args.Repository cannot be nil")
	}
	if args.JwtSecret == "" {
		return errors.New("Args.JwtSecret cannot be empty")
	}
	if args.Username == "" {
		return errors.New("Args.Username cannot be empty")
	}
	if args.Password == "" {
		return errors.New("Args.Password cannot be empty")
	}
	return nil
}

func Start(port uint, args Args) (err error) {
	logger.Logger().Info("-----------------------------------------------------------")
	logger.Logger().Info("START - http web api server")
	logger.Logger().Infof("< port = %d", port)
	logger.Logger().Info("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Info("-----------------------------------------------------------")
		if err == nil {
			logger.Logger().Info("END   - http web api server")
		} else {
			logger.Logger().Error("FAILED - http web api server")
			logger.Logger().Error(err)
		}
		logger.Logger().Info("-----------------------------------------------------------")
	}()

	if err = args.validate(); err != nil {
		return err
	}
	repository = args.Repository
	jwtIssuer = &args.JwtIssuer
	jwtSecret = &args.JwtSecret
	username = &args.Username
	password = &args.Password

	e := echo.New()
	e.Use(middleware.Gzip())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: args.AllowOrigins,
	}))
	e.HideBanner = true
	e.HidePort = true
	registerHandler(e, *jwtIssuer, *jwtSecret)

	err = e.Start(fmt.Sprintf(":%d", port))
	return err
}
