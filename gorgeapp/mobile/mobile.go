//go:build android || mobile
// +build android mobile

// Package mobile provides platform initialization for golang/x/mobile
package mobile

import (
	"image"
	"log"
	"runtime"
	"time"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/systems/input"
	"github.com/stdiopt/gorge/systems/render/gl"
	"github.com/stdiopt/gorge/systems/resource"

	"github.com/stdiopt/gorge/gorgeapp/mobile/app"

	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
)

func init() {
	runtime.LockOSThread()
}

// update this to be similar to other platforms with a mobileSystem.

// Start stars a mobile app
func Run(opt Options, systems ...interface{}) error {
	resourceFS := resource.MobileFS{}

	var g *gorge.Gorge
	var glw *gl.Wrapper
	var bounds image.Rectangle
	done := make(chan struct{})

	mark := time.Now()
	app.Handle(func(a app.App, ee interface{}) bool {
		e := a.Filter(ee)
		if e, ok := e.(lifecycle.Event); ok {
			switch e.Crosses(lifecycle.StageVisible) {
			case lifecycle.CrossOn:
				log.Println("Starting mobile")
				glw = &gl.Wrapper{}
				ggArgs := []interface{}{
					func(g *gorge.Context, res *resource.Context) {
						res.AddFS("/", resourceFS)
						g.PutProp(glw)
						g.SetScreenSize(m32.Vec2{float32(bounds.Dx()), float32(bounds.Dy())})
					},
				}
				ggArgs = append(ggArgs, systems...)

				g = gorge.New(ggArgs...)
				g.Start()
			case lifecycle.CrossOff:
				g.Close()
				return false
				close(done)
			}
		}
		if g == nil {
			return true
		}
		// Again but with stuff
		switch e := e.(type) {
		case paint.Event:
			if glw == nil {
				return true
			}

			now := time.Now()
			sub := float32(now.Sub(mark)) / 1000000000
			g.Update(sub)
			go func() {
				a.Publish()
				a.Send(paint.Event{})
			}() // keep animating
			mark = now
		case mouse.Event:
			switch e.Button {
			case mouse.ButtonWheelUp:
				in
				g.Trigger(input.EventPointer{
					Type: input.MouseWheel,
					Pointers: map[int]input.PointerData{
						0: {
							DeltaZ: -1,
							Pos:    m32.Vec2{e.X, e.Y},
						},
					},
				})
			case mouse.ButtonWheelDown:
				g.Trigger(input.EventPointer{
					Type: input.MouseWheel,
					Pointers: map[int]input.PointerData{
						0: {
							DeltaZ: 1,
							Pos:    m32.Vec2{e.X, e.Y},
						},
					},
				})
			}
		case key.Event:
			log.Println("KeyboardEvent:", e)
		case touch.Event:
			var typ input.PointerType = input.MouseMove
			switch e.Type {
			case touch.TypeBegin:
				typ = input.MouseDown
			case touch.TypeEnd:
				typ = input.MouseUp
			}

			g.Trigger(input.EventPointer{
				Type: typ,
				Pointers: map[int]input.PointerData{
					0: {Pos: m32.Vec2{e.X, e.Y}},
				},
			})
		case size.Event:
			bounds = e.Bounds()
			g.SetScreenSize(m32.Vec2{float32(bounds.Dx()), float32(bounds.Dy())})
		}
	})
	<-done
	log.Println("Exiting app")
	return nil
}
