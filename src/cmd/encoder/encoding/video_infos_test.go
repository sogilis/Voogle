package encoding

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_videoHaveSound(t *testing.T) {
	t.SkipNow()
	cases := []struct {
		Name          string
		GivenFilepath string
		ExpectSound   bool
		ExpectError   bool
	}{
		{Name: "With Sound", GivenFilepath: "../../../../samples/1280x720_2mb.mp4", ExpectSound: true, ExpectError: false},
		{Name: "Without Sound", GivenFilepath: "../../../../samples/aerial.mp4", ExpectSound: false, ExpectError: false},
	}

	for _, tt := range cases {
		t.Run(tt.Name, func(t *testing.T) {
			sound, err := checkContainsSound(tt.GivenFilepath)
			if tt.ExpectError {
				assert.NotNil(t, err)
				return
			}
			assert.Nil(t, err)
			assert.Equal(t, tt.ExpectSound, sound)
		})
	}
}
