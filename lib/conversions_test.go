package lib

import (
	"reflect"
	"testing"

	"github.com/edio/randrctl2/profile"
	"github.com/edio/randrctl2/x"
	"github.com/stretchr/testify/assert"
)

func Test_toProfileOutput(t *testing.T) {
	tests := []struct {
		name      string
		xOutput   *x.Output
		assertion func(t *testing.T, actual *profile.Output)
	}{
		{
			"should convert x.Output ot profile.Output",
			&x.Output{
				Crtc: 3,
				Mode: &x.Mode{
					Resolution: x.Geometry{1280, 720},
					Rate:       60,
					Flags:      4,
				},
				Position:      x.Geometry{1920, 1080},
				Panning:       x.Geometry{1366, 768},
				Scale:         1,
				RotationFlags: 2,
				// do not matter for this test
				Id:             x.OutputId(0),
				Name:           "",
				Edid:           []byte{},
				SupportedModes: []*x.Mode{},
				PreferredMode:  &x.Mode{},
			},
			func(t *testing.T, actual *profile.Output) {
				expected := profile.Output{
					Crtc: 3,
					Mode: profile.Mode{
						Resolution: "1280x720",
						RateHint:   60,
						FlagsHint:  toProfileModeFlags(x.ModeFlags(4)),
					},
					Panning:  "1366x768",
					Position: "1920x1080",
					Rotation: toProfileRotation(x.RotationFlags(2)),
					Scale:    1,
				}
				assert.Equal(t, expected, *actual)
			},
		},
		{
			"should round Mode.Rate to 2 decimals",
			&x.Output{
				Mode: &x.Mode{
					Rate: 59.9453,
					// do not matter for this test
					Resolution: x.Geometry{0, 0},
					Flags:      0,
				},
				// do not matter for this test
				Position:       x.Geometry{0, 0},
				Panning:        x.Geometry{0, 0},
				Scale:          0,
				RotationFlags:  0,
				Id:             x.OutputId(0),
				Name:           "",
				Crtc:           0,
				Edid:           []byte{},
				SupportedModes: []*x.Mode{},
				PreferredMode:  &x.Mode{},
			},
			func(t *testing.T, actual *profile.Output) {
				assert.Equal(t, 59.95, actual.Mode.RateHint)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := toProfileOutput(tt.xOutput)
			tt.assertion(t, actual)
		})
	}
}

func Test_toProfileRotation(t *testing.T) {
	tests := []struct {
		name string
		rf   x.RotationFlags
		want []profile.Rotation
	}{
		{"handle 0th LSB", 1 << 0, []profile.Rotation{profile.Rotate0}},
		{"handle 1st LSB", 1 << 1, []profile.Rotation{profile.Rotate90}},
		{"handle 2nd LSB", 1 << 2, []profile.Rotation{profile.Rotate180}},
		{"handle 3rt LSB", 1 << 3, []profile.Rotation{profile.Rotate270}},
		{"handle 4th LSB", 1 << 4, []profile.Rotation{profile.ReflectX}},
		{"handle 5th LSB", 1 << 5, []profile.Rotation{profile.ReflectY}},
		{"handle only 6 LSBs combined", 0xFFFF /* 16 bits */, []profile.Rotation{
			profile.Rotate0,
			profile.Rotate90,
			profile.Rotate180,
			profile.Rotate270,
			profile.ReflectX,
			profile.ReflectY,
		},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toProfileRotation(tt.rf); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toProfileRotation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_toProfileModeFlags(t *testing.T) {
	tests := []struct {
		name string
		mf   x.ModeFlags
		want []profile.ModeFlag
	}{
		{"handle 0th LSB", 1 << 0, []profile.ModeFlag{profile.HsyncPositive}},
		{"handle 1st LSB", 1 << 1, []profile.ModeFlag{profile.HsyncNegative}},
		{"handle 2nd LSB", 1 << 2, []profile.ModeFlag{profile.VsyncPositive}},
		{"handle 3rd LSB", 1 << 3, []profile.ModeFlag{profile.VsyncNegative}},
		{"handle 4th LSB", 1 << 4, []profile.ModeFlag{profile.Interlace}},
		{"handle 5th LSB", 1 << 5, []profile.ModeFlag{profile.DoubleScan}},
		{"handle 6th LSB", 1 << 6, []profile.ModeFlag{profile.Csync}},
		{"handle 7th LSB", 1 << 7, []profile.ModeFlag{profile.CsyncPositive}},
		{"handle 8th LSB", 1 << 8, []profile.ModeFlag{profile.CsyncNegative}},
		{"handle 9th LSB", 1 << 9, []profile.ModeFlag{profile.HskewPresent}},
		{"handle 10th LSB", 1 << 10, []profile.ModeFlag{profile.Bcast}},
		{"handle 11th LSB", 1 << 11, []profile.ModeFlag{profile.PixelMultiplex}},
		{"handle 12th LSB", 1 << 12, []profile.ModeFlag{profile.DoubleClock}},
		{"handle 13th LSB", 1 << 13, []profile.ModeFlag{profile.HalveClock}},
		{"handle only 14 LSBs combined", 0xFFFF /* 16 bits */, []profile.ModeFlag{
			profile.HsyncPositive,
			profile.HsyncNegative,
			profile.VsyncPositive,
			profile.VsyncNegative,
			profile.Interlace,
			profile.DoubleScan,
			profile.Csync,
			profile.CsyncPositive,
			profile.CsyncNegative,
			profile.HskewPresent,
			profile.Bcast,
			profile.PixelMultiplex,
			profile.DoubleClock,
			profile.HalveClock,
		},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toProfileModeFlags(tt.mf); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toProfileModeFlags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToProfile(t *testing.T) {
	type args struct {
		connected []*x.Output
		primary   *x.Output
	}
	tests := []struct {
		name      string
		args      args
		assertion func(t *testing.T, actual *profile.Profile)
	}{
		{
			"should create empty profile",
			args{
				[]*x.Output{},
				nil,
			},
			func(t *testing.T, actual *profile.Profile) {
				assert.Equal(t, 0, len(actual.Outputs))
				assert.Equal(t, 0, len(actual.Match))
				assert.Equal(t, 0, len(actual.Primary))
			},
		},
		{
			"should create profile without primary",
			args{
				[]*x.Output{
					{
						Id:   x.OutputId(1),
						Name: "Output1",
						Mode: &x.Mode{
							Resolution: x.Geometry{1280, 720},
							Rate:       60,
							Flags:      4,
						},
						Position:       x.Geometry{0, 0},
						Panning:        x.Geometry{1280, 720},
						Scale:          1,
						RotationFlags:  1,
						Crtc:           3,
						Edid:           []byte("edid"),
						PreferredMode:  &x.Mode{
							Resolution: x.Geometry{1920, 1080},
							Rate:       60,
							Flags:      4,
						},
						// does not matter
						SupportedModes: []*x.Mode{},
					},
				},
				nil,
			},
			func(t *testing.T, actual *profile.Profile) {
				assert.Equal(t, 1, len(actual.Outputs))
				output := *actual.Outputs["Output1"]
				assert.Equal(t, "1280x720", output.Mode.Resolution)

				assert.Equal(t, 1, len(actual.Match))
				rule := *actual.Match["Output1"]
				assert.Equal(t, profile.Rule {
					Edid: hash([]byte("edid")),
					Supports: "1280x720",
					Prefers: "1920x1080",
				}, rule)
				assert.Equal(t, 0, len(actual.Primary))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := ToProfile(tt.args.connected, tt.args.primary)
			tt.assertion(t, actual)
		})
	}
}
