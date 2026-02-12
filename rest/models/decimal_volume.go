package models

import (
	"strconv"
)

type DecimalVolume string

func (dv DecimalVolume) Float64() (float64, error) {
	return strconv.ParseFloat(string(dv), 64)
}

func (dv DecimalVolume) String() string {
	return string(dv)
}
