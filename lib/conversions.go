package lib

import (
	"crypto/md5"
	"fmt"
	"github.com/BurntSushi/xgb/randr"
	"github.com/edio/randrctl2/profile"
	"github.com/edio/randrctl2/x"
	"math"
)

func ToProfile(connected []*x.Output, primary *x.Output) *profile.Profile {
	outputs := make(map[string]*profile.Output, 0)
	rules := make(map[string]*profile.Rule, 0)

	for _, xOutput := range connected {
		rule := profile.Rule{
			Edid: hash(xOutput.Edid),
		}
		rules[xOutput.Name] = &rule

		if xOutput.PreferredMode != nil {
			rule.Prefers = toGeometryString(xOutput.PreferredMode.Resolution)
		}

		if xOutput.IsActive() {
			output := toProfileOutput(xOutput)
			outputs[xOutput.Name] = output

			rule.Supports = output.Mode.Resolution
		}
	}

	result := profile.Profile{
		Match:   rules,
		Outputs: outputs,
	}

	if primary != nil {
		result.Primary = primary.Name
	}

	return &result
}

func hash(data []byte) string {
	hasher := md5.New()
	hasher.Write(data)
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

func toProfileOutput(xOutput *x.Output) *profile.Output {
	rateRounded := math.Round(xOutput.Mode.Rate*100) / 100
	return &profile.Output{
		Crtc: xOutput.Crtc,
		Mode: profile.Mode{
			Resolution: toGeometryString(xOutput.Mode.Resolution),
			RateHint:   rateRounded,
			FlagsHint:  toProfileModeFlags(xOutput.Mode.Flags),
		},
		Panning:  toGeometryString(xOutput.Panning),
		Position: toGeometryString(xOutput.Position),
		Rotation: toProfileRotation(xOutput.RotationFlags),
		Scale:    xOutput.Scale,
	}
}

func toGeometryString(geometry x.Geometry) string {
	return fmt.Sprintf("%dx%d", geometry[0], geometry[1])
}

func toProfileModeFlags(mf x.ModeFlags) []profile.ModeFlag {
	flags := make([]profile.ModeFlag, 0)
	if mf&randr.ModeFlagHsyncPositive != 0 {
		flags = append(flags, profile.HsyncPositive)
	}
	if mf&randr.ModeFlagHsyncNegative != 0 {
		flags = append(flags, profile.HsyncNegative)
	}
	if mf&randr.ModeFlagVsyncPositive != 0 {
		flags = append(flags, profile.VsyncPositive)
	}
	if mf&randr.ModeFlagVsyncNegative != 0 {
		flags = append(flags, profile.VsyncNegative)
	}
	if mf&randr.ModeFlagInterlace != 0 {
		flags = append(flags, profile.Interlace)
	}
	if mf&randr.ModeFlagDoubleScan != 0 {
		flags = append(flags, profile.DoubleScan)
	}
	if mf&randr.ModeFlagCsync != 0 {
		flags = append(flags, profile.Csync)
	}
	if mf&randr.ModeFlagCsyncPositive != 0 {
		flags = append(flags, profile.CsyncPositive)
	}
	if mf&randr.ModeFlagCsyncNegative != 0 {
		flags = append(flags, profile.CsyncNegative)
	}
	if mf&randr.ModeFlagHskewPresent != 0 {
		flags = append(flags, profile.HskewPresent)
	}
	if mf&randr.ModeFlagBcast != 0 {
		flags = append(flags, profile.Bcast)
	}
	if mf&randr.ModeFlagPixelMultiplex != 0 {
		flags = append(flags, profile.PixelMultiplex)
	}
	if mf&randr.ModeFlagDoubleClock != 0 {
		flags = append(flags, profile.DoubleClock)
	}
	if mf&randr.ModeFlagHalveClock != 0 {
		flags = append(flags, profile.HalveClock)
	}
	return flags
}

func toProfileRotation(rf x.RotationFlags) []profile.Rotation {
	rotation := make([]profile.Rotation, 0)
	if rf&randr.RotationRotate0 != 0 {
		rotation = append(rotation, profile.Rotate0)
	}
	if rf&randr.RotationRotate90 != 0 {
		rotation = append(rotation, profile.Rotate90)
	}
	if rf&randr.RotationRotate180 != 0 {
		rotation = append(rotation, profile.Rotate180)
	}
	if rf&randr.RotationRotate270 != 0 {
		rotation = append(rotation, profile.Rotate270)
	}
	if rf&randr.RotationReflectX != 0 {
		rotation = append(rotation, profile.ReflectX)
	}
	if rf&randr.RotationReflectY != 0 {
		rotation = append(rotation, profile.ReflectY)
	}
	return rotation
}
