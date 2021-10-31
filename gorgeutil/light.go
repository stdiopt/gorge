package gorgeutil

import (
	"math"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/m32"
)

// Light is a light entity which can be used directly on gorge.
type Light struct {
	Name string
	gorge.TransformComponent
	gorge.LightComponent
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
			Type:        gorge.LightPoint,
			Intensity:   100,
			Color:       [3]float32{1, 1, 1},
			Range:       1000,
			CastShadows: true,
		},
	}
}

// NewPointLight returns a new PointLight.
func NewPointLight() *Light {
	// NewLight defaults to pointlight.
	return NewLight()
}

// NewDirectionalLight returns a directional light.
func NewDirectionalLight() *Light {
	return &Light{
		TransformComponent: gorge.TransformIdent(),
		LightComponent: gorge.LightComponent{
			Type:        gorge.LightDirectional,
			Intensity:   1,
			Color:       [3]float32{1, 1, 1},
			Range:       1000,
			CastShadows: true,
		},
	}
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
			InnerConeCos: m32.Cos(30 * (math.Pi / 180)),
			OuterConeCos: m32.Cos(40 * (math.Pi / 180)),
			CastShadows:  true,
		},
	}
}
