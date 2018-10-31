package x

import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/randr"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/edio/randrctl2/profile"
)

type XError struct {
	cause error
}

func (err *XError) Error() string {
	return err.cause.Error()
}

var (
	x           *xgb.Conn
	rootWindow  xproto.Window
	resources   *randr.GetScreenResourcesReply
	modeInfoIdx map[randr.Mode]randr.ModeInfo
)

func Connect(display string) error {
	var err error
	x, err = xgb.NewConnDisplay(display)
	if err != nil {
		return &XError{err}
	}
	err = randr.Init(x)
	if err != nil {
		return &XError{err}
	}

	// initialize package resources
	rootWindow = xproto.Setup(x).DefaultScreen(x).Root

	resources, err = randr.GetScreenResources(x, rootWindow).Reply()
	if err != nil {
		return &XError{err}
	}

	// index some resources for easier access
	modeInfoIdx = make(map[randr.Mode]randr.ModeInfo)
	for _, mode := range resources.Modes {
		modeInfoIdx[randr.Mode(mode.Id)] = mode
	}

	return nil
}

func Disconnect() {
	if x != nil {
		x.Close()

		x = nil
		rootWindow = 0
		resources = nil
		modeInfoIdx = nil
	}
}

type Geometry [2]int

type ModeFlags uint32

func (f ModeFlags) ToProfileModeFlags() []profile.ModeFlag {
	flags := make([]profile.ModeFlag, 0)
	if f&randr.ModeFlagHsyncPositive != 0 {
		flags = append(flags, profile.HsyncPositive)
	}
	if f&randr.ModeFlagHsyncNegative != 0 {
		flags = append(flags, profile.HsyncNegative)
	}
	if f&randr.ModeFlagVsyncPositive != 0 {
		flags = append(flags, profile.VsyncPositive)
	}
	if f&randr.ModeFlagVsyncNegative != 0 {
		flags = append(flags, profile.VsyncNegative)
	}
	if f&randr.ModeFlagInterlace != 0 {
		flags = append(flags, profile.Interlace)
	}
	if f&randr.ModeFlagDoubleScan != 0 {
		flags = append(flags, profile.DoubleScan)
	}
	if f&randr.ModeFlagCsync != 0 {
		flags = append(flags, profile.Csync)
	}
	if f&randr.ModeFlagCsyncPositive != 0 {
		flags = append(flags, profile.CsyncPositive)
	}
	if f&randr.ModeFlagCsyncNegative != 0 {
		flags = append(flags, profile.CsyncNegative)
	}
	if f&randr.ModeFlagHskewPresent != 0 {
		flags = append(flags, profile.HskewPresent)
	}
	if f&randr.ModeFlagBcast != 0 {
		flags = append(flags, profile.Bcast)
	}
	if f&randr.ModeFlagPixelMultiplex != 0 {
		flags = append(flags, profile.PixelMultiplex)
	}
	if f&randr.ModeFlagDoubleClock != 0 {
		flags = append(flags, profile.DoubleClock)
	}
	if f&randr.ModeFlagHalveClock != 0 {
		flags = append(flags, profile.HalveClock)
	}
	return flags
}

type Mode struct {
	Resolution Geometry
	Rate       float64
	Flags      ModeFlags
}

type OutputId uint32

type RotationFlags uint16

func (rf RotationFlags) ToProfileRotation() []profile.Rotation {
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

type Output struct {
	Id   OutputId
	Name string
	Crtc           int
	Edid           []byte
	SupportedModes []*Mode
	PreferredMode  *Mode
	Mode           *Mode
	Position       Geometry
	Panning        Geometry
	Scale          float64
	RotationFlags  RotationFlags
}

func (o *Output) IsActive() bool {
	return o.Mode != nil
}

func FindPrimary(connections []*Output) (int, *Output, error) {
	resp, err := randr.GetOutputPrimary(x, rootWindow).Reply()
	if err != nil {
		return -1, nil, &XError{err}
	}

	id := OutputId(resp.Output)

	for i, connection := range connections {
		if id == connection.Id {
			return i, connection, nil
		}
	}
	return -1, nil, nil
}

func GetOutputNames() ([]string, error) {
	names := make([]string, resources.NumOutputs)
	for oi, outputId := range resources.Outputs {
		outputInfo, err := randr.GetOutputInfo(x, outputId, 0).Reply()
		if err != nil {
			return nil, &XError{err}
		}

		names[oi] = string(outputInfo.Name)
	}
	return names, nil
}

func GetConnectedOutputs() ([]*Output, error) {
	outputs := make([]*Output, 0)
	for _, outputId := range resources.Outputs {
		outputInfo, err := randr.GetOutputInfo(x, outputId, 0).Reply()
		if err != nil {
			return nil, &XError{err}
		}

		if outputInfo.Connection != randr.ConnectionConnected {
			continue
		}

		output := Output{
			Name: string(outputInfo.Name),
		}

		// Edid
		properties, _ := randr.ListOutputProperties(x, outputId).Reply()
		for _, propAtom := range properties.Atoms {
			name, _ := xproto.GetAtomName(x, propAtom).Reply()
			if name.Name == "EDID" {
				prop, _ := randr.GetOutputProperty(x, outputId, propAtom, 0, 0, 100, false, false).Reply()
				output.Edid = prop.Data
			}
		}

		// Monitor.SupportedModes and PreferredMode
		supportedModes := make([]*Mode, outputInfo.NumModes)
		for i, modeId := range outputInfo.Modes {
			modeInfo := modeInfoIdx[randr.Mode(modeId)]
			rate := float64(modeInfo.DotClock) / (float64(modeInfo.Htotal) * float64(modeInfo.Vtotal))
			supportedModes[i] = &Mode{
				Resolution: Geometry{
					int(modeInfo.Width),
					int(modeInfo.Height),
				},
				Rate: rate,
			}
			if i < int(outputInfo.NumPreferred) {
				output.PreferredMode = supportedModes[i]
			}
		}
		output.SupportedModes = supportedModes

		if outputInfo.Crtc > 0 {
			// output is active

			// xrandr refers to crtc by its index in response. Let's do the same here
			for i, crtcId := range outputInfo.Crtcs {
				if crtcId == outputInfo.Crtc {
					output.Crtc = i
				}
			}

			crtcInfo, err := randr.GetCrtcInfo(x, outputInfo.Crtc, 0).Reply()
			if err != nil {
				return nil, &XError{err}
			}
			modeInfo := modeInfoIdx[crtcInfo.Mode]
			rate := float64(modeInfo.DotClock) / (float64(modeInfo.Htotal) * float64(modeInfo.Vtotal))

			output.Mode = &Mode{
				Resolution: Geometry{
					int(modeInfo.Width),
					int(modeInfo.Height),
				},
				Rate:  rate,
				Flags: ModeFlags(modeInfo.ModeFlags),
			}

			output.Panning = Geometry{
				int(crtcInfo.Width),
				int(crtcInfo.Height),
			}

			output.Position = Geometry{
				int(crtcInfo.X),
				int(crtcInfo.Y),
			}

			output.RotationFlags = RotationFlags(crtcInfo.Rotation)

			// TODO implement scaling
			output.Scale = 1
		}

		outputs = append(outputs, &output)
	}

	return outputs, nil
}
