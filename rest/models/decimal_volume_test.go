package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecimalVolume_Float64(t *testing.T) {
	tests := []struct {
		name    string
		dv      DecimalVolume
		want    float64
		wantErr assert.ErrorAssertionFunc
	}{
		{"dv 2.5", DecimalVolume("2.5"), 2.5, assert.NoError},
		{"dv invalid", DecimalVolume("2.b"), 0, assert.Error},
		{"dv high precision", DecimalVolume("2.323233223223223223322"), 2.323233223223223, assert.NoError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.dv.Float64()
			if !tt.wantErr(t, err, "Float64()") {
				return
			}
			assert.Equalf(t, tt.want, got, "Float64()")
		})
	}
}

func TestDecimalVolume_String(t *testing.T) {
	tests := []struct {
		name string
		dv   DecimalVolume
		want string
	}{
		{"dv 2.5", DecimalVolume("2.5"), "2.5"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.dv.String(), "String()")
		})
	}
}
