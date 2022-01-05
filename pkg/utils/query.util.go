package utils

import (
	"net/url"
	"strconv"
)

func GetQueryVariable(q *url.Values, name string, fallback interface{}) interface{} {
	val := q.Get(name)

	switch v := fallback.(type) {
	case int:
		v, err := strconv.Atoi(val)
		if err != nil {
			return fallback
		}
		return v
	default:
		if val == "" {
			return fallback
		}
		return val
	}
}
