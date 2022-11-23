package flag_with_env

import (
	"systemd-cd/domain/model/flag_with_env"
)

func Parse() {
	flag_with_env.Parse()
}

func Uint(paramName string, envName string, fallback uint, desc string) *uint {
	return flag_with_env.Uint(paramName, envName, fallback, desc)
}

func String(paramName string, envName string, fallback string, desc string) *string {
	return flag_with_env.String(paramName, envName, fallback, desc)
}

func Array(paramName string, desc string) *flag_with_env.ArrayParam {
	return flag_with_env.Array(paramName, desc)
}
