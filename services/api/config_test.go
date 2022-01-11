package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigPort(t *testing.T) {
	cases := []struct {
		name          string
		valueToParse  string
		expectedValue uint32
		wantError     bool
	}{
		{name: "NoValue", valueToParse: "", expectedValue: 4444, wantError: false},
		{name: "Default value equal to default", valueToParse: "4444", expectedValue: 4444, wantError: false},
		{name: "Numerical value", valueToParse: "5555", expectedValue: 5555, wantError: false},
		{name: "Incorrect value", valueToParse: "44A4", expectedValue: 0, wantError: true},
	}

	for _, tt := range cases {
		fmt.Println("Testing:", tt.name)

		//given
		assert.Nil(t, os.Unsetenv("PORT"))
		// if os.Setenv is call with "", it wreaks the env parse library
		if tt.valueToParse != "" {
			os.Setenv("PORT", tt.valueToParse)
		}

		// when
		config, err := NewConfig()

		// then
		if tt.wantError {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, tt.expectedValue, config.Port)
		}
	}

}
