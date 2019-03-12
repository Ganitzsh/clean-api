package api

import (
	"net/http"
	"strconv"
)

func readLimOff(r *http.Request) (lim int, off int) {
	if r == nil {
		return 0, 0
	}
	val, err := strconv.Atoi(r.URL.Query().Get("lim"))
	if err == nil {
		lim = val
	}
	val, err = strconv.Atoi(r.URL.Query().Get("off"))
	if err == nil {
		off = val
	}
	return lim, off
}