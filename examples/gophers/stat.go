// Copyright 2019 Luis Figueiredo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gophers

import (
	"bytes"
	"fmt"
	"math"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/gl"
	"github.com/stdiopt/gorge/gorgeutils"
	"github.com/stdiopt/gorge/input"
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/primitive"
	"github.com/stdiopt/gorge/x/text"
)

// Stat do the thing
func Stat(s *gorge.Scene) {
	// Add necessary stuff here
	//assets := asset.FromECS(g)

	//root := &gorge.Transform{}

	font := s.Assets().Font("fonts/meslo.ttf")

	txt := text.New(font)
	txtMat := s.Assets().Material("shaders/unlit")
	txtMat.SetTexture("albedoMap", font.Texture)

	txt.SetMaterial(txtMat)
	// TODO: find a better way
	//txt.Material.ShaderLoader = txtMat.ShaderLoader

	txt.SetScale(.02).SetPosition(-0.75, 0, 1)

	plane := primitive.Plane()
	plane.Color = vec4{0, 0, 0, 0.7}
	plane.SetParent(txt).SetEuler(-math.Pi/2, 0, 0)

	var camera *gorgeutils.Camera
	var winSize vec2
	s.AddEntity(plane, txt)

	recalc := func() {
		if camera == nil {
			return
		}
		m := camera.Projection()
		pvInv := m.Inv()

		p := vec2{20, 20}
		ndc := m32.Vec4{
			2*p[0]/winSize[0] - 1,
			1 - 2*p[1]/winSize[1],
			-1, 1}
		dir := pvInv.Mul4x1(ndc).Vec3().Normalize()
		r := m32.IntersectPlane(
			dir,
			m32.Vec3{0, 0, 0}, // zero point here
			m32.Vec3{0, 0, 1},
			m32.Vec3{0, 0, 1},
		)
		txt.SetPositionv(r)
		txt.SetScale(.02 * (1 / winSize[1] * 800))
	}
	// Goal make text on a right position
	s.Handle(func(c *gorgeutils.Camera) {
		camera = c
		txt.SetParent(camera)
	})
	s.Handle(func(e gorge.ResizeEvent) {
		winSize = vec2(e)
	})

	// Profiling
	padding := float32(1.0)

	timeInitial := float32(2)
	timeInterval := timeInitial
	s.Handle(func(e gorge.UpdateEvent) {
		recalc()
		timeInterval -= float32(e)
		if timeInterval > 0 {
			return
		}
		timeInterval = timeInitial

		statStr := statUpdate(s)
		txt.SetText(statStr)

		delta := txt.Max.Sub(txt.Min)
		half := delta.Mul(0.5)
		center := txt.Min.Add(half)
		// Might be updated with new string
		plane.SetPosition(center[0], center[1], 0.01)

		plane.SetScale(half[0]+padding, 0, half[1]+padding)
	})
	s.Handle(func(evt input.KeyEvent) {
		if evt.Type == input.KeyUp {
			switch evt.Key {
			case "F10":
				statUpdate(s)
			case "F9":
				runtime.GC()
			}
		}
	}).Describe("stat key")

	//panel := primitive.Plane("", m32.Vec4{0, 0, 0, 0.2})
}

// Read stats into a formated string
func statUpdate(s *gorge.Scene) string {
	var winSize m32.Vec2
	s.Handle(func(evt gorge.ResizeEvent) {
		winSize = m32.Vec2(evt)
	})

	var gw *gl.Wrapper
	s.Query(func(evt *gl.Wrapper) { gw = evt })

	memStat := runtime.MemStats{}
	buf := &bytes.Buffer{}

	fmt.Fprintf(buf, "Width: %.2f, Height: %.2f\n", winSize[0], winSize[1])
	fmt.Fprintf(buf, "Renderer: %v\n", gw.GetString(gl.RENDERER))
	fmt.Fprintf(buf, "Gophers Count: %d\n", nThings)
	runtime.ReadMemStats(&memStat)
	fmt.Fprintf(buf,
		"GC: %v Pause: %v, CurMem: %.2fk\n",
		memStat.NumGC,
		time.Duration(memStat.PauseNs[(memStat.NumGC+255)%256]),
		float64(memStat.Alloc)/1024,
	)

	hgroups := []*gorge.HandlerGroup{}
	s.Range(func(k, v interface{}) bool {
		hgroups = append(hgroups, v.(*gorge.HandlerGroup))
		return true
	})
	// Since its a map we maintain an order
	/*for _, e := range m.Messaging.Entries {
		entries = append(entries, e)
	}*/
	sort.SliceStable(hgroups, func(i, j int) bool {
		return hgroups[i].Type.Name() < hgroups[j].Type.Name()
	})

	for _, hg := range hgroups {
		dt := hg.CallEnd.Sub(hg.CallStart)
		/*fmt.Fprintf(buf, "Handler: %s delta: %v fps: %.2f\n", e.Type.Name(), dt.Round(time.Millisecond/100), float64(time.Second)/float64(dt))
		for _, h := range e.Handlers {
			fmt.Fprintf(buf, "  delta %v: %s\n", h.CallEnd.Sub(h.CallStart), h.Desc)
		}*/

		if hg.Type.Name() == "RenderEvent" {
			fmt.Fprintf(buf, "(Renderer) PostUpdate: %v fps: %.2f\n", dt.Round(time.Millisecond/100), float64(time.Second)/float64(dt))
		}
		for _, h := range hg.Handlers {
			if strings.HasPrefix(h.Desc, "gopher") {
				dt := h.CallEnd.Sub(h.CallStart)
				fmt.Fprintf(buf, "(CPU) GophersUpdate: %v fps: %.2f\n", dt.Round(time.Millisecond/100), float64(time.Second)/float64(dt))
			}
		}
	}
	fmt.Fprintf(buf, "Transforms:      %10v", gorge.TransformBuilds)
	return buf.String()
}
