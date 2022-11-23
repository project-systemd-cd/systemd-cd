package flag_with_env

import (
	"os"
	"strconv"
)

// uint型で環境変数取得
func getUintEnv(key string, fallback uint) uint {
	// 環境変数取得
	if value, ok := os.LookupEnv(key); ok {
		// uint型に変換
		var intValue, err = strconv.ParseUint(value, 10, 16)
		if err == nil {
			return uint(intValue)
		}
	}
	// uint型に変換失敗 または 環境変数が設定されていない場合、引数に指定したfallbackを返す
	return fallback
}

// string型で環境変数取得
func getEnv(key, fallback string) string {
	// 環境変数取得
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	// 環境変数が設定されていない場合は引数に指定したfallbackを返す
	return fallback
}
