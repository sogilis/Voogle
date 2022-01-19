package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/Sogilis/Voogle/services/api/config"
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
		t.Run(tt.name, func(t *testing.T) {
			// Given
			assert.Nil(t, os.Unsetenv("PORT"))

			// Application required basic auth
			// to allow this set var env
			os.Setenv("USER_AUTH", "user")
			os.Setenv("PWD_AUTH", "pwd")

			// if os.Setenv is call with "", it wreaks the env parse library
			if tt.valueToParse != "" {
				os.Setenv("PORT", tt.valueToParse)
			}

			// When
			config, err := NewConfig()

			// Then
			if tt.wantError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.expectedValue, config.Port)
			}
		})
	}
	assert.Nil(t, os.Unsetenv("PORT"))
}
func TestConfigBasicAuth(t *testing.T) {
	cases := []struct {
		name              string
		givenUser         string
		givenPwd          string
		userExpectedValue string
		pwdExpectedValue  string
		wantError         bool
	}{
		{name: "NoValue", givenUser: "", givenPwd: "", userExpectedValue: "", pwdExpectedValue: "", wantError: true},
		{name: "OnlyUser", givenUser: "test", givenPwd: "", userExpectedValue: "test", pwdExpectedValue: "", wantError: true},
		{name: "OnlyPwd", givenUser: "", givenPwd: "pwd", userExpectedValue: "", pwdExpectedValue: "pwd", wantError: true},
		{name: "Default", givenUser: "test", givenPwd: "pwd", userExpectedValue: "test", pwdExpectedValue: "pwd", wantError: false},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			assert.Nil(t, os.Unsetenv("PORT"))
			assert.Nil(t, os.Unsetenv("USER_AUTH"))
			assert.Nil(t, os.Unsetenv("PWD_AUTH"))

			// if os.Setenv is call with "", it wreaks the env parse library
			if tt.givenUser != "" {
				os.Setenv("USER_AUTH", tt.givenUser)
			}
			if tt.givenPwd != "" {
				os.Setenv("PWD_AUTH", tt.givenPwd)
			}

			// When
			config, err := NewConfig()

			// Then
			if tt.wantError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.userExpectedValue, config.UserAuth)
				assert.Equal(t, tt.pwdExpectedValue, config.PwdAuth)
			}
		})
	}
	assert.Nil(t, os.Unsetenv("PORT"))
	assert.Nil(t, os.Unsetenv("USER_AUTH"))
	assert.Nil(t, os.Unsetenv("PWD_AUTH"))
}
