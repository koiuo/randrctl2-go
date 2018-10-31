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
						Crtc: 0,
						Mode: Mode{
							Resolution: "1920x1080",
						},
						Panning:  "1920x1200",
						Position: "1920x0",
						Rotation: []Rotation{Rotate0},
						Scale:    1.4,
					},
				},
			},
			`
			outputs:
			  LVDS1:
			    crtc: 0
			    mode:
			      resolution: 1920x1080
			    panning: 1920x1200
			    position: 1920x0
			    rotation:
			    - rotate0
			    scale: 1.4
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
						Crtc: 0,
						Mode: Mode{
							Resolution: "1920x1080",
							RateHint: 60,
							FlagsHint: []ModeFlag{
								HsyncPositive,
								VsyncNegative,
							},
						},
						Panning:  "1920x1200",
						Position: "0x0",
						Rotation: []Rotation{Rotate0},
						Scale:    1.4,
					},
					"DP1": {
						Crtc: 1,
						Mode: Mode{
							Resolution: "3840x2160",
							RateHint: 60,
							FlagsHint: []ModeFlag{
								Interlace,
							},
						},
						Panning:  "3840x2160",
						Position: "1920x0",
						Rotation: []Rotation{Rotate270, ReflectY},
						Scale:    2,
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
			    crtc: 1
			    mode:
			      resolution: 3840x2160
			      ratehint: 60
			      flaghint:
			      - interlace
			    panning: 3840x2160
			    position: 1920x0
			    rotation:
			    - rotate270
			    - reflecty
			    scale: 2
			  LVDS1:
			    crtc: 0
			    mode:
			      resolution: 1920x1080
			      ratehint: 60
			      flaghint:
			      - hsync+
			      - vsync-
			    panning: 1920x1200
			    position: "0x0"
			    rotation:
			    - rotate0
			    scale: 1.4
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
