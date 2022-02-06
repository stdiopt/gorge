package gorge

import (
	"fmt"

	"github.com/stdiopt/gorge/math/gm"
)

// LightType type for lights
type LightType int

// Light types
const (
	LightPoint = LightType(iota)
	LightDirectional
	LightSpot
)

func (l LightType) String() string {
	switch l {
	case LightPoint:
		return "LightPoint"
	case LightDirectional:
		return "LightDirectional"
	case LightSpot:
		return "LightSpot"
	default:
		return fmt.Sprintf("LightUnknown(%d)", l)
	}
}

// LightComponent component type of light and what nots
// type of light as in
// position and direction determined by transform Z direction
type LightComponent struct {
	Type      LightType // default point
	Intensity float32
	Color     gm.Vec3
	Range     float32

	InnerConeCos float32
	OuterConeCos float32

	DisableShadow bool
}

// NewLightComponent returns a New light component with some defaults (pointsLight)
func NewLightComponent() *LightComponent {
	return &LightComponent{
		Type:      LightPoint,
		Intensity: 100,
		Color:     gm.Vec3{1, 1, 1},
		Range:     100,
	}
}

// Light method to satisfy component
func (l *LightComponent) Light() *LightComponent { return l }

// SetType sets the light type Directional, Spot, Point.
func (l *LightComponent) SetType(t LightType) {
	l.Type = t
}

// SetColor sets light Color
func (l *LightComponent) SetColor(r, g, b float32) {
	l.Color = gm.Vec3{r, g, b}
}

// SetIntensity gets light intensity
func (l *LightComponent) SetIntensity(v float32) {
	l.Intensity = v
}

// SetRange sets Point or Spot light range
func (l *LightComponent) SetRange(v float32) {
	l.Range = v
}

// SetCastShadows convinient accessor that sets the CastShadows field and
// returns self.
func (l *LightComponent) SetDisableShadow(b bool) {
	l.DisableShadow = b
}

// SetInnerConeCos sets the InnerConeCos for spot lights.
func (l *LightComponent) SetInnerConeCos(v float32) {
	l.InnerConeCos = v
}

// SetOuterConeCos sets the OuterConeCos for spot lights.
func (l *LightComponent) SetOuterConeCos(v float32) {
	l.OuterConeCos = v
}
