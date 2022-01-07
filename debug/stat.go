package debug

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"runtime"
	"time"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/gorgeutil"
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/primitive"
	"github.com/stdiopt/gorge/static"
	"github.com/stdiopt/gorge/systems/input"
	"github.com/stdiopt/gorge/systems/render"
	"github.com/stdiopt/gorge/systems/render/gl"
	"github.com/stdiopt/gorge/systems/resource"
	"github.com/stdiopt/gorge/text"
)

const padding = 1.0

var gi = 0

type statSystem struct {
	gorge          *gorge.Context
	camera         *gorgeutil.Camera
	txt            *text.Entity
	rendererStat   render.EventStat
	updateDuration time.Duration
	glRenderer     string
	glVersion      string
}

// StatText logs stats on default logger.
func StatText(g *gorge.Context) {
	ic := input.FromContext(g)
	s := &statSystem{
		gorge:      g,
		glVersion:  gl.GetString(gl.VERSION),
		glRenderer: gl.GetString(gl.RENDERER),
	}
	gi++

	timeInitial := float32(2)
	timeInterval := float32(0)

	var preMark time.Time
	gorge.HandleFunc(g, func(gorge.EventPreUpdate) {
		preMark = time.Now()
	})
	gorge.HandleFunc(g, func(e gorge.EventUpdate) {
		if ic.KeyUp(input.KeyF10) {
			s.txt.SetText(s.Update())
		}
		if ic.KeyUp(input.KeyF9) {
			runtime.GC()
			s.txt.SetText(s.Update())
		}

		timeInterval -= float32(e)
		if timeInterval > 0 {
			return
		}
		timeInterval = timeInitial

		fmt.Println("\033[01;37m-----------------------------------")
		fmt.Println(s.Update())
		fmt.Println("\033[0m")
	})
	gorge.HandleFunc(g, func(gorge.EventRender) {
		s.updateDuration = time.Since(preMark)
	})
	gorge.HandleFunc(g, func(e render.EventStat) {
		s.rendererStat = e
	})
}

// Stat inits stat system
func Stat(g *gorge.Context) error {
	rc := resource.FromContext(g)
	ic := input.FromContext(g)
	s := &statSystem{
		gorge:      g,
		glVersion:  gl.GetString(gl.VERSION),
		glRenderer: gl.GetString(gl.RENDERER),
	}
	gi++

	font := &text.Font{}
	if err := rc.Load(font, "_gorge/fonts/font.ttf", text.FontResolution(512)); err != nil {
		return err
	}

	s.camera = gorgeutil.NewOrthoCamera(20, -100, 100)
	s.camera.SetCullMask(1 << 16)
	s.camera.SetClearFlag(gorge.ClearDepthOnly)
	s.camera.SetOrder(10)
	g.Add(s.camera)

	planeMat := gorge.NewShaderMaterial(static.Shaders.Unlit)
	texMat := gorge.NewShaderMaterial(static.Shaders.Unlit)
	texMat.SetTexture("albedoMap", font.Texture)

	s.txt = text.New(font)
	s.txt.SetMaterial(texMat)
	s.txt.SetText(s.Update())
	s.txt.Color = m32.Vec4{1, 1, 1, 1}
	s.txt.SetScale(.02)
	s.txt.SetPosition(-0.75, 0, 5)
	// s.txt.Position[2] = 5
	s.txt.CullMask = 1 << 16
	s.txt.CastShadow = gorge.CastShadowDisabled

	plane := gorgeutil.NewPlane(primitive.PlaneDirY)
	plane.SetMaterial(planeMat)
	plane.Color = m32.Vec4{0, 0, 0, 0.8}
	plane.Blend = gorge.BlendOneOneMinusSrcAlpha
	plane.SetParent(s.txt)
	plane.SetEuler(math.Pi/2, 0, 0)
	plane.CullMask = 1 << 16
	plane.CastShadow = gorge.CastShadowDisabled

	g.Add(plane, s.txt)

	timeInitial := float32(1)
	timeInterval := float32(0)

	var preMark time.Time
	gorge.HandleFunc(g, func(gorge.EventPreUpdate) {
		preMark = time.Now()
	})
	gorge.HandleFunc(g, func(e gorge.EventUpdate) {
		s.recalc()

		if ic.KeyUp(input.KeyF10) {
			log.Println("Update")
			s.txt.SetText(s.Update())
		}
		if ic.KeyUp(input.KeyF9) {
			log.Println("GC")
			runtime.GC()
			s.txt.SetText(s.Update())
		}

		timeInterval -= float32(e)
		if timeInterval > 0 {
			return
		}
		timeInterval = timeInitial

		s.txt.SetText(s.Update())

		delta := s.txt.Max.Sub(s.txt.Min)
		half := delta.Mul(0.5)
		center := s.txt.Min.Add(half)
		// Might be updated with new string
		plane.SetPosition(center[0], center[1], -0.01)

		plane.SetScale(half[0]+padding, 0, half[1]+padding)
	})
	gorge.HandleFunc(g, func(gorge.EventRender) {
		s.updateDuration = time.Since(preMark)
	})
	gorge.HandleFunc(g, func(e render.EventStat) {
		s.rendererStat = e
	})
	return nil
}

func (s *statSystem) recalc() {
	ws := s.gorge.ScreenSize()
	cam := s.camera
	aspect := ws[0] / ws[1]
	h := cam.OrthoSize / 2
	pad := padding * .03 * h
	s.txt.SetPosition(-h*aspect+pad, h-pad, 0)
	s.txt.SetScale(.05 * h * (1 / ws[1] * 800))
}

// Update Read stats into a formated string
func (s *statSystem) Update() string {
	memStat := runtime.MemStats{}
	buf := &bytes.Buffer{}
	ws := s.gorge.ScreenSize()
	fmt.Fprintf(buf, "GORGE (F9) runtime.GC (F10) statUpdate now\n")
	fmt.Fprintf(buf, "Width: %.2f, Height: %.2f\n", ws[0], ws[1])
	fmt.Fprintf(buf, "GL Wrapper: %v\n", gl.Global())
	fmt.Fprintf(buf, "GL Version: %v\n", s.glVersion)
	fmt.Fprintf(buf, "GL Renderer: %v\n", s.glRenderer)

	fmt.Fprintf(buf, "gpu vbos: %v, tex: %v, shaders: %v, ins: %v, buffers: %v\n",
		s.rendererStat.VBOs,
		s.rendererStat.Textures,
		s.rendererStat.Shaders,
		s.rendererStat.Instances,
		s.rendererStat.Buffers,
	)

	runtime.ReadMemStats(&memStat)
	fmt.Fprintf(buf,
		"GC: %v Pause: %v, CurMem: %.2fk\n",
		memStat.NumGC,
		time.Duration(memStat.PauseNs[(memStat.NumGC+255)%256]),
		float64(memStat.Alloc)/1024,
	)
	fmt.Fprintf(buf, "(PostUpdate): %v fps: %.2f\n",
		s.updateDuration.Round(time.Millisecond/100),
		float64(time.Second)/float64(s.updateDuration),
	)
	dt := s.rendererStat.RenderDuration
	fmt.Fprintf(buf, "(Renderer): %v fps: %.2f\n",
		dt.Round(time.Millisecond/100),
		float64(time.Second)/float64(dt),
	)
	fmt.Fprintf(buf, "Drawcalls:       %10v\n", s.rendererStat.DrawCalls)
	fmt.Fprintf(buf, "Transforms:      %10v, Saved: %v", gorge.TransformBuilds, gorge.TransformBuildSave)
	return buf.String()
}
