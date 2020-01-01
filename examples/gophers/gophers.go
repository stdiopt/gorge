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
	glog "log"
	"math"
	"math/rand"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/gorgeutils"
	"github.com/stdiopt/gorge/input"
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/primitive"
	"github.com/stdiopt/gorge/x/text"
)

var (
	log = glog.New(glog.Writer(), "(gophers) ", 0)
)

const (
	nThings   = 15000
	timeScale = 20
	areaX     = 15
	areaY     = 10
)

// Cool aliases
type (
	vec2 = m32.Vec2
	vec3 = m32.Vec3
	vec4 = m32.Vec4
	mat4 = m32.Mat4
)

var (
	texList = []string{
		"gopher", "wasm", "wood", "grid",
	}
	textures map[string]*gorge.Texture
)

// Thing entity is a mover unit in screen controlled by the customSystem
type Thing struct {
	gorge.Transform
	gorge.Renderable
	custom
}

// Reset the thing
func (t *Thing) Reset(sz vec2) {
	t.Color = vec4{}
	t.speed = rand.Float32() * 0.2
	t.dir = rand.Float32() * math.Pi * 2
	t.life = 1
	t.lifeS = rand.Float32() * 0.01
	t.SetPosition(
		rand.Float32()*sz[0]*2-sz[0],
		0,
		rand.Float32()*sz[1]*2-sz[1],
	)
}

// Custom component
type custom struct {
	turner float32
	dir    float32
	speed  float32
	life   float32
	lifeS  float32
}

type gophersSystem struct {
	scene *gorge.Scene
	//gorge *gorge.Gorge
	//assets  *asset.System
	winSize vec2

	camTrans  vec3
	cameraRig *gorge.Transform
	camera    *gorgeutils.Camera
	light     *gorgeutils.Light

	ground       *primitive.MeshEntity
	things       []*Thing
	pointerLoc   *gorge.Transform
	pointerShape *primitive.MeshEntity
	pointerText  *text.Text
	gimbal       *gorgeutils.Gimbal
	cube         *primitive.MeshEntity
	dog          *primitive.MeshEntity
	dogText      *text.Text

	font      *gorge.Font
	minDist   float32
	totalTime float32
}

// System gopher system initializes a scene
func System(g *gorge.Gorge) {
	s := g.Scene(Gophers, Stat)
	log.Println("Loading scene")
	g.LoadScene(s)
	log.Println("Loaded")
	s.Start() // Start linkt stuff and what nots
}

// Gophers starts the gopher scne
func Gophers(s *gorge.Scene) {
	log.Println("Gophers starting")

	assets := s.Assets()
	font := assets.Font("fonts/open-sans.ttf", gorge.FontResolution(2048))

	dogMesh := assets.Mesh("obj/dog.obj")
	dogTex := assets.Texture2D("obj/dog.jpg")
	// Renderer create texture?
	textures = map[string]*gorge.Texture{
		"gopher": assets.Texture2D("gopher.png"),
		"wood":   assets.Texture2D("wood.png"),
		"grid":   assets.Texture2D("grid.png"),
		"wasm":   assets.Texture2D("wasm.png"),
	}

	gs := gophersSystem{
		scene: s,
		//gorge: gm,
		//assets:       assets,
		//input:        input,
		cameraRig:    gorge.NewTransform(),
		camera:       gorgeutils.NewCamera(),
		light:        gorgeutils.NewLight(),
		ground:       primitive.Cube(),
		pointerLoc:   gorge.NewTransform(),
		pointerShape: primitive.Plane(),
		pointerText:  text.New(font),
		gimbal:       gorgeutils.NewGimbal(),
		cube:         primitive.Cube(),
		dogText:      text.New(font),
		dog: &primitive.MeshEntity{
			Transform: *gorge.NewTransform(),
			Renderable: gorge.Renderable{
				Mesh:  dogMesh,
				Color: vec4{1, 1, 1, 1},
				Material: gorge.NewMaterial(nil).
					SetFloat32("metallic", 0).
					SetFloat32("ao", 10).
					SetTexture("albedoMap", dogTex),
			},
		},
		font:    font,
		minDist: 4,
	}

	txtMat := s.Assets().Material("shaders/unlit")
	txtMat.SetTexture("albedoMap", font.Texture)
	gs.dogText.SetMaterial(txtMat)

	//gs.dogText.Material.ShaderLoader = txtMat.ShaderLoader

	s.Persist(gs.camera)

	s.Handle(func(evt gorge.ResizeEvent) {
		gs.winSize = vec2(evt)
		gs.camera.Camera.AspectRatio = gs.winSize[0] / gs.winSize[1]
	})

	s.Handle(gs.pointerHandler())
	s.Handle(gs.Update).Describe("gophers-update")
	s.Handle(gs.Start)

	gs.Setup()
}

// TODO: fix this crap
func (s *gophersSystem) pointerHandler() func(evt input.PointerEvent) {
	var dragging = 0
	var lastP vec2
	var camRotVec vec2
	var camRot = vec2{-0.4, 0}

	var lastPinch float32
	var pinching bool

	return func(evt input.PointerEvent) {
		delta := vec2(evt.Pointers[0].Pos).Sub(lastP)
		lastP = vec2(evt.Pointers[0].Pos)
		if evt.Type == input.MouseWheel {
			dist := s.camera.WorldPosition().Len()
			multiplier := dist * 0.005
			s.camera.Translate(0, 0, -evt.Pointers[0].DeltaZ*multiplier)
			return
		}

		switch len(evt.Pointers) {
		case 1: // Only one pointer
			if evt.Type == input.MouseDown || evt.Type == input.PointerDown {
				p := s.screenToYPlane(vec2(evt.Pointers[0].Pos))
				dragging = 1
				cursor := s.pointerLoc.Position
				halfDist := float32(s.minDist * 0.5)
				min := vec2{cursor[0] - halfDist, cursor[2] - halfDist}
				max := vec2{cursor[0] + halfDist, cursor[2] + halfDist}
				if in2d(vec2{p[0], p[2]}, min, max) {
					dragging = 2
				}
			}
			if evt.Type == input.MouseUp || evt.Type == input.PointerEnd {
				dragging = 0
				pinching = false
			}
		}
		if dragging == 0 {
			return
		}
		// dragging state 2
		if dragging == 2 {
			nv := s.screenToYPlane(vec2(evt.Pointers[0].Pos))
			s.pointerLoc.SetPositionv(nv)
			p := s.pointerLoc.WorldPosition()
			s.gimbal.LookAt(p, m32.Up())
			return
		}

		if evt.Type == input.MouseMove || evt.Type == input.PointerMove {
			if len(evt.Pointers) == 1 {
				scale := float32(0.005)
				camRotVec = vec2{-delta[1], delta[0]}.Mul(scale)
				camRot = camRot.Add(camRotVec)

				s.cameraRig.SetRotation(m32.QuatEuler(camRot[0], camRot[1], 0))

			}
			if len(evt.Pointers) == 2 {
				v := vec2(evt.Pointers[0].Pos).Sub(vec2(evt.Pointers[1].Pos))
				curPinch := v.Len()
				if !pinching {
					lastPinch = curPinch
					pinching = true
				}
				deltaPinch := curPinch - lastPinch
				lastPinch = curPinch
				s.camera.Translate(0, 0, deltaPinch*0.1)
			}
		}
	}
}

func (s *gophersSystem) Setup() {
	s.createGophers()
	// Setup camera
	s.cameraRig.Rotate(0.4, 0, 0)
	s.camera.Camera.
		SetPerspective(math.Pi/4, s.winSize[0]/s.winSize[1], 0.1, 1000).
		SetAmbient(0.4, 0.4, 0.4)
	s.camera.
		SetParent(s.cameraRig).
		SetEuler(0, 0, 0).
		SetPosition(0, 0, -17)
	// Camera stuff

	// Set Ground
	s.ground.
		SetPosition(0, -1.05, 0).
		SetScale(areaX+.2, 1, areaY+0.2)
	s.ground.Material.
		SetTexture("albedoMap", textures["wood"]).
		Set("roughness", float32(0.1)).
		Set("metallic", float32(0.2)).
		Set("ao", float32(5))

	// Setup big gopher (pointer)
	s.pointerShape.
		SetParent(s.pointerLoc).
		SetScale(s.minDist)
	s.pointerShape.Material.
		SetTexture("albedoMap", textures["gopher"])

	s.pointerText.Material.Depth = false
	s.pointerText.SetColor(vec4{0, 0, 0, 1})
	s.pointerText.SetParent(s.pointerLoc).
		SetEuler(math.Pi/2, 0, 0).
		SetScale(0.4)

	s.gimbal.SetPosition(0, 1, -4)

	s.cube.
		SetPosition(0, areaY*0.5, areaY+1).
		SetScale(areaX*0.5, areaY*0.5, 1)
	s.cube.Material.
		SetFloat32("roughness", 0.2).
		SetFloat32("ao", 1)

	dogLoc := gorge.NewTransform().SetPosition(-areaX+2, 0, areaY-2)
	s.dog.SetParent(dogLoc).
		SetEuler(math.Pi/2, 0, 0).
		SetScale(0.1)
	s.dogText.SetText("random dog")
	s.dogText.SetColor(vec4{1, 1, 1, 1})

	s.dogText.SetParent(dogLoc).
		SetPosition(1, 3, -1).
		SetScale(0.6)

	s.light.
		SetParent(s.pointerLoc).
		SetPosition(0, 4, 0)
	s.light.Color = vec3{1, 1, 1}
	lightGimbal := gorgeutils.NewGimbal()
	lightGimbal.SetParent(s.light)

	g := s.scene
	g.AddEntity(s.ground)
	g.AddEntity(thingsToEntities(s.things)...)
	g.AddEntity(s.pointerShape, s.pointerText)
	g.AddEntity(s.gimbal.Entities...)
	g.AddEntity(s.cube)
	g.AddEntity(s.dog, s.dogText)
	g.AddEntity(s.camera)
	g.AddEntity(s.light)
	g.AddEntity(lightGimbal.Entities...)
}
func (s *gophersSystem) Start(evt gorge.StartEvent) {
	log.Println("START EVENT....")
	for _, t := range s.things {
		t.Reset(s.winSize)
	}
}

func (s *gophersSystem) Update(evt gorge.UpdateEvent) {
	s.totalTime += float32(evt)
	dt := float32(evt) * timeScale
	count := 0
	for _, t := range s.things {

		t.life -= t.lifeS
		if t.life <= 0 {
			t.Reset(vec2{areaX, areaY})
		}

		speed := t.speed
		dir := t.dir
		opacity := m32.Sin((1 - t.life) * math.Pi)

		nearest := float32(1000)
		target := s.pointerLoc.WorldPosition()

		p := t.Position
		// Dist from Point
		dx, dy := target[0]-p[0], target[2]-p[2]
		dist := m32.Hypot(dx, dy)
		if dist >= nearest {
			continue
		}
		dir = m32.Atan2(dy, dx)
		nearest = dist
		switch {
		case nearest < 0.1:
			count++
			t.life = m32.Max(0.5, t.life)
			t.dir = dir
			speed = 0
			//t.Reset(vec2{areaX, areaY})
		case nearest < s.minDist:
			count++
			t.Color = vec4{0.9, 0.9, 1, opacity}
			t.life = m32.Max(0.3, t.life)
			t.dir = dir
			speed = 0.3
		default:
			t.Color = vec4{0.8, 0.8, 0.8, opacity}
			t.turner = m32.Limit(t.turner+(float32(rand.NormFloat64())*0.2), -0.2, 0.2)
			t.dir += t.turner * dt
		}
		//dog Area
		if t.Position[0] < -areaX+4 && t.Position[2] > areaY-4 {
			t.Position[0] = 0
			t.Position[2] = 0
		}

		// Move gophers
		sin, cos := m32.Sincos(t.dir)

		t.SetEuler(0, t.dir, 0)

		position := t.Position
		np := position.Add(vec3{cos * speed * dt, 0, sin * speed * dt})
		t.SetPosition(
			m32.Limit(np[0], -areaX, areaX),
			np[1],
			m32.Limit(np[2], -areaY, areaY),
		)
	}
	s.minDist = m32.Min(2+float32(count)/nThings*8, 10)
	s.pointerShape.SetScale(s.minDist)

	s.pointerText.SetTextf("Gophers: %v", count)
	s.pointerText.SetPosition(-s.pointerText.Max[0]/4, 0, -0.8*s.minDist)

	// XXX: Testing things
	s.cube.Material.Set("metallic", m32.Sin(s.totalTime))

	pickTex := int(s.totalTime*0.3) % len(texList)
	s.cube.Material.SetTexture("albedoMap", textures[texList[pickTex]])

	input := input.FromECS(s.scene)
	if input == nil {
		panic("gophers requires input")
	}
	const mmax = 4
	const stp = .1
	if input.GetKey("a") {
		s.camTrans[0] = m32.Max(s.camTrans[0]-stp, -mmax)
	}
	if input.GetKey("d") {
		s.camTrans[0] = m32.Min(s.camTrans[0]+stp, mmax)
	}
	if input.GetKey("w") {
		s.camTrans[2] = m32.Min(s.camTrans[2]+stp, mmax)
	}
	if input.GetKey("s") {
		s.camTrans[2] = m32.Max(s.camTrans[2]-stp, -mmax)
	}
	if input.GetKey("c") {
		s.camTrans = vec3{}
		s.cameraRig.Position = vec3{}
	}
	s.cameraRig.Translatev(s.camTrans.Mul(dt))
	s.camTrans = s.camTrans.Mul(.9)
}

func (s *gophersSystem) screenToYPlane(p vec2) vec3 {
	m := s.camera.Camera.Projection()
	m = m.Mul4(s.camera.Inv())
	//PVInv := s.camera.Mat4().Inv()
	PVInv := m.Inv()
	ndc := vec4{2*p[0]/s.winSize[0] - 1, 1 - 2*p[1]/s.winSize[1], -1, 1}
	dir := PVInv.Mul4x1(ndc).Vec3().Normalize()

	cp := PVInv.Col(3) // Camera position
	return m32.IntersectPlane(
		dir,
		vec3{cp[0] / cp[3], cp[1] / cp[3], cp[2] / cp[3]},
		vec3{0, 1, 0}, // plane
		vec3{0, 0, 0},
	)
}

func (s *gophersSystem) createGophers() {
	log.Println("Adding NThings:", nThings)

	mat := gorge.NewMaterial(nil)
	mat.Depth = false
	mat.Set("metallic", float32(0.5)).
		Set("roughness", float32(0.8)).
		SetTexture("albedoMap", textures["gopher"])

	mesh := primitive.PlaneMesh()

	ret := []*Thing{}
	// Creating entities
	for i := 0; i < nThings; i++ {
		t := Thing{
			*gorge.NewTransform(),
			gorge.Renderable{
				Name:     "gopher",
				Mesh:     mesh,
				Material: mat,
			},
			custom{},
		}
		t.SetEuler(0, 0, 0).
			SetScale(0.4)
		t.Reset(s.winSize)
		ret = append(ret, &t)
	}
	s.things = ret
}
func thingsToEntities(things []*Thing) []gorge.Entity {
	ret := make([]gorge.Entity, len(things))
	for i := range things {
		ret[i] = things[i]
	}
	return ret
}

func in2d(p vec2, min vec2, max vec2) bool {
	if p[0] < min[0] || p[0] > max[0] {
		return false
	}
	if p[1] < min[1] || p[1] > max[1] {
		return false
	}
	return true
}
