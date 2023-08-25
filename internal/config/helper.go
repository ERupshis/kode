package config

import (
	"fmt"
	"strconv"
)

// SUPPORT FUNCTIONS.

func setEnvToParamIfNeed(param interface{}, val string) {
	if val == "" {
		return
	}

	switch param := param.(type) {
	case *int64:
		if envVal, err := atoi64(val); err == nil {
			*param = envVal
		} else {
			panic(err)
		}
	case *string:
		*param = val
	default:
		panic(fmt.Errorf("wrong input param type"))
	}
}

func atoi64(value string) (int64, error) {
	return strconv.ParseInt(value, 10, 64)
}
