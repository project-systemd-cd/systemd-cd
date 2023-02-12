package echo

import (
	"errors"
	"fmt"
	"systemd-cd/domain/logger"
	"systemd-cd/domain/runner"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var repository runner.IRepositoryInmemory

type Args struct {
	Repository runner.IRepositoryInmemory
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

	if args.Repository == nil {
		return errors.New("Args.Repository cannot be nil")
	}
	repository = args.Repository

	e := echo.New()
	e.Use(middleware.Gzip())
	e.HideBanner = true
	e.HidePort = true
	registerHandler(e)

	err = e.Start(fmt.Sprintf(":%d", port))
	return err
}
