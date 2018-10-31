package profile

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"regexp"
	"strings"
	"testing"
)

func unindent(str string) string {
	nonwhitespace := regexp.MustCompile("^\n\\s+")
	indent := nonwhitespace.FindString(str)
	if len(indent) > 0 {
		return strings.Replace(str, indent, "\n", -1)[1:]
	} else {
		return str
	}
}

func TestWrite(t *testing.T) {
	tests := []struct {
		name       string
		profile    Profile
		wantWriter string
		assertErr  func(assert.TestingT, error, ...interface{}) bool
	}{
		{
			"should write minimal profile",
			Profile{
				Outputs: map[string]*Output{
					"LVDS1": {
						Mode:    "1920x1080",
						Pos:     "1920x0",
						Rotate:  "normal",
						Panning: "1920x1200",
						Scale:   "1.4x1.4",
						Rate:    "75",
						Crtc:    0,
					},
				},
			},
			`
			outputs:
			  LVDS1:
			    mode: 1920x1080
			    pos: 1920x0
			    rotate: normal
			    panning: 1920x1200
			    scale: 1.4x1.4
			    rate: "75"
			    crtc: 0
			`,
			assert.NoError,
		},
		{
			"should write full profile sorting keys and ignoring Name field",
			Profile{
				Name: "should be transient",
				Match: map[string]*Rule{
					"LVDS1": {
						Edid:     "70b13ad1e146a7e9a63a3e1f733996bb",
						Prefers:  "1920x1080",
						Supports: "1920x1080",
					},
					"DP1": {
						Edid:     "73e0b78b21eccb78174dc4325ab459e6",
						Prefers:  "3840x2160",
						Supports: "3840x2160",
					},
				},
				Outputs: map[string]*Output{
					"LVDS1": {
						Mode:    "1920x1080",
						Pos:     "0x0",
						Rotate:  "normal",
						Panning: "1920x1200",
						Scale:   "1.4x1.4",
						Rate:    "75",
						Crtc:    0,
					},
					"DP1": {
						Mode:    "3840x2160",
						Pos:     "1920x0",
						Rotate:  "right",
						Panning: "3840x2160",
						Scale:   "2x2",
						Rate:    "75",
						Crtc:    1,
					},
				},
				Primary: "DP1",
			},
			`
			match:
			  DP1:
			    edid: 73e0b78b21eccb78174dc4325ab459e6
			    prefers: 3840x2160
			    supports: 3840x2160
			  LVDS1:
			    edid: 70b13ad1e146a7e9a63a3e1f733996bb
			    prefers: 1920x1080
			    supports: 1920x1080
			outputs:
			  DP1:
			    mode: 3840x2160
			    pos: 1920x0
			    rotate: right
			    panning: 3840x2160
			    scale: 2x2
			    rate: "75"
			    crtc: 1
			  LVDS1:
			    mode: 1920x1080
			    pos: "0x0"
			    rotate: normal
			    panning: 1920x1200
			    scale: 1.4x1.4
			    rate: "75"
			    crtc: 0
			primary: DP1
			`,
			assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &bytes.Buffer{}
			err := Write(writer, &tt.profile)
			tt.assertErr(t, err)
			assert.Equal(t, unindent(tt.wantWriter), writer.String())
		})
	}
}
