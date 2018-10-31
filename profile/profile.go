package profile

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
)

type Rotation string

const (
	Rotate0   Rotation = "rotate0"
	Rotate90  Rotation = "rotate90"
	Rotate180 Rotation = "rotate180"
	Rotate270 Rotation = "rotate270"
	ReflectX  Rotation = "reflectx"
	ReflectY  Rotation = "reflecty"
)

type ModeFlag string

const (
	HsyncPositive  ModeFlag = "hsync+"
	HsyncNegative  ModeFlag = "hsync-"
	VsyncPositive  ModeFlag = "vsync+"
	VsyncNegative  ModeFlag = "vsync-"
	Interlace      ModeFlag = "interlace"
	DoubleScan     ModeFlag = "doublescan"
	Csync          ModeFlag = "csync"
	CsyncPositive  ModeFlag = "csync+"
	CsyncNegative  ModeFlag = "csync-"
	HskewPresent   ModeFlag = "hskew"
	Bcast          ModeFlag = "bcast"
	PixelMultiplex ModeFlag = "pixelmultiplex"
	DoubleClock    ModeFlag = "doubleclock"
	HalveClock     ModeFlag = "halveclock"
)

type Profile struct {
	Name    string             `yaml:"-"`
	Match   map[string]*Rule   `yaml:"match,omitempty"`
	Outputs map[string]*Output `yaml:"outputs,omitempty"`
	Primary string             `yaml:"primary,omitempty"`
}

type Rule struct {
	Edid     string `yaml:"edid,omitempty"`
	Prefers  string `yaml:"prefers,omitempty"`
	Supports string `yaml:"supports,omitempty"`
}

type Mode struct {
	Resolution string     `yaml:"resolution"`
	RateHint   *float64   `yaml:"ratehint,omitempty"`
	FlagsHint  []ModeFlag `yaml:"flaghint,omitempty"`
}

type Output struct {
	Crtc     int        `yaml:"crtc"`
	Mode     Mode       `yaml:"mode"`
	Panning  string     `yaml:"panning"`
	Position string     `yaml:"position"`
	Rotation []Rotation `yaml:"rotation"`
	Scale    float64    `yaml:"scale"`
}

func Write(writer io.Writer, profile *Profile) error {
	if len(profile.Outputs) == 0 {
		return fmt.Errorf("outputs are empty: %v", profile)
	}
	enc := yaml.NewEncoder(writer)
	defer enc.Close()
	return enc.Encode(*profile)
}

func Read(reader io.Reader) (*Profile, error) {
	dec := yaml.NewDecoder(reader)
	p := Profile{}
	err := dec.Decode(&p)
	return &p, err
}

//func FromConnections(connections []*x.Output) (*Profile, error) {
//	outputs := make([]*Output, len(connections))
//	for i, connection := range connections {
//		println(connection)
//		outputs[i] = &Output{
//			//Mode: connection.Mode
//		}
//	}
//	return nil, nil
//}

func serializeGeometry() {

}
