package main

import (
	"crypto/md5"
	"fmt"
	"github.com/edio/randrctl/profile"
	"github.com/edio/randrctl/x"
	"math"

	"gopkg.in/yaml.v2"
	"os"
)

func main() {
	err := x.Connect(":0")
	if err != nil {
		os.Exit(1)
	}
	defer x.Disconnect()

	connected, _ := x.GetConnectedOutputs()
	outputs := make(map[string]*profile.Output, 0)
	rules := make(map[string]*profile.Rule, 0)

	for _, c := range connected {
		hasher := md5.New()
		hasher.Write(c.Edid)
		rules[c.Name] = &profile.Rule{
			Edid:     fmt.Sprintf("%x", hasher.Sum(nil)),
		}

		if c.IsActive() {
			rateRounded := math.Round(c.Mode.Rate*100) / 100

			output := &profile.Output{
				Crtc: c.Crtc,
				Mode: profile.Mode{
					Resolution: fmt.Sprintf("%dx%d", c.Mode.Resolution[0], c.Mode.Resolution[1]),
					RateHint:   &rateRounded,
					FlagsHint:  c.Mode.Flags.ToProfileModeFlags(),
				},
				Panning:  fmt.Sprintf("%dx%d", c.Panning[0], c.Panning[1]),
				Position: fmt.Sprintf("%dx%d", c.Position[0], c.Panning[1]),
				Rotation: c.RotationFlags.ToProfileRotation(),
				Scale:    c.Scale,
			}
			outputs[c.Name] = output

			rules[c.Name].Supports = output.Mode.Resolution
			rules[c.Name].Prefers = output.Mode.Resolution
		}
	}

	primaryName := ""

	_, primary, _ := x.FindPrimary(connected)

	if primary != nil {
		primaryName = primary.Name
	}

	profile := profile.Profile{
		Match:   rules,
		Outputs: outputs,
		Primary: primaryName,
	}

	enc := yaml.NewEncoder(os.Stdout)
	enc.Encode(profile)
	//cmd.Execute()
}
