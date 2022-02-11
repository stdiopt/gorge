package gorgeutil

import (
	"fmt"
	"math"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/math/gm"
)

// Light is a light entity which can be used directly on gorge.
type Light struct {
	Name string
	gorge.TransformComponent
	gorge.LightComponent
}

func (l *Light) String() string {
	return fmt.Sprintf("gorgeutil.Light(%v)", l.Type)
}

// SetName sets the light name for debugging purposes.
func (l *Light) SetName(n string) {
	l.Name = n
}

// NewLight returns a new default light.
func NewLight() *Light {
	return &Light{
		Name:               "",
		TransformComponent: gorge.TransformIdent(),
		LightComponent: gorge.LightComponent{
			Type:      gorge.LightPoint,
			Intensity: 100,
			Color:     [3]float32{1, 1, 1},
			Range:     1000,
		},
	}
}

func AddLight(a Contexter) *Light {
	l := NewLight()
	a.Add(l)
	return l
}

// NewPointLight returns a new PointLight.
func NewPointLight() *Light {
	// NewLight defaults to pointlight.
	return NewLight()
}

func AddPointLight(a Contexter) *Light {
	l := NewPointLight()
	a.Add(l)
	return l
}

// NewDirectionalLight returns a directional light.
func NewDirectionalLight() *Light {
	return &Light{
		TransformComponent: gorge.TransformIdent(),
		LightComponent: gorge.LightComponent{
			Type:      gorge.LightDirectional,
			Intensity: 1,
			Color:     [3]float32{1, 1, 1},
			Range:     1000,
		},
	}
}

func AddDirectionalLight(a Contexter) *Light {
	l := NewDirectionalLight()
	a.Add(l)
	return l
}

// NewSpotLight returns a spot light.
func NewSpotLight() *Light {
	return &Light{
		TransformComponent: gorge.TransformIdent(),
		LightComponent: gorge.LightComponent{
			Type:         gorge.LightPoint,
			Intensity:    100,
			Color:        [3]float32{1, 1, 1},
			Range:        100,
			InnerConeCos: gm.Cos(30 * (math.Pi / 180)),
			OuterConeCos: gm.Cos(40 * (math.Pi / 180)),
		},
	}
}

func AddSpotLight(a Contexter) *Light {
	l := NewSpotLight()
	a.Add(l)
	return l
}
