package echo

import (
	"net/http"

	"systemd-cd/domain/errors"
	"systemd-cd/domain/logger"

	"github.com/labstack/echo/v4"
)

type BodyUsersSignin struct {
	Username *string `json:"username"`
	Password *string `json:"password"`
}

func (b BodyUsersSignin) Validate() error {
	if b.Username == nil {
		return &errors.ErrValidation{
			Property:    "username",
			Given:       nil,
			Description: "cannot be empty",
		}
	} else if *b.Username == "" {
		return &errors.ErrValidation{
			Property:    "username",
			Given:       b.Username,
			Description: "cannot be empty",
		}
	}
	if b.Password == nil {
		return &errors.ErrValidation{
			Property:    "password",
			Given:       nil,
			Description: "cannot be empty",
		}
	} else if *b.Password == "" {
		return &errors.ErrValidation{
			Property:    "password",
			Given:       b.Password,
			Description: "cannot be empty",
		}
	}
	return nil
}

type ResUsersSigninPost struct {
	Token string `json:"token"`
}

func usersSigninPost(c echo.Context) (err error) {
	logger.Logger().Trace("-----------------------------------------------------------")
	logger.Logger().Trace("START - POST /usrs/sign_in")
	logger.Logger().Tracef("< RemoteAddr = %s", c.Request().RemoteAddr)
	logger.Logger().Trace("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Trace("-----------------------------------------------------------")
		if err == nil {
			logger.Logger().Tracef("> Status = %d", c.Response().Status)
			logger.Logger().Tracef("> ContentLength = %d", c.Response().Size)
			logger.Logger().Tracef("END    - POST /usrs/sign_in %d", c.Response().Status)
		} else {
			logger.Logger().Error("FAILED - POST /usrs/sign_in")
			logger.Logger().Error(err)
			err = c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}
		logger.Logger().Trace("-----------------------------------------------------------")
	}()

	body := new(BodyUsersSignin)
	if err = c.Bind(body); err != nil {
		err = c.JSONPretty(http.StatusBadRequest, map[string]string{"message": err.Error()}, "	")
		return err
	} else {
		err = body.Validate()
		if err != nil {
			err = c.JSONPretty(http.StatusBadRequest, map[string]string{"message": err.Error()}, "	")
			return err
		}
	}

	if *username == *body.Username && *password == *body.Password {
		var token string
		token, err = GenerateToken(GenerateTokenParam{Username: *body.Username})
		err = c.JSONPretty(http.StatusOK, ResUsersSigninPost{Token: token}, "	")
		return err
	}
	err = c.JSONPretty(http.StatusUnauthorized, map[string]string{"message": "failed to sign in"}, "	")
	return err
}
