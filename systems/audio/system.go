// Package audio implements audio system
package audio

import (
	"log"
	"time"

	"github.com/hajimehoshi/oto"
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/systems/audio/proc"
)

// System initializes the audio system.
func System(g *gorge.Context) { FromContext(g) }

// System initializes the gorge system.
/*func System(g *gorge.Context) error {
	log.Println("Initializing system")
	audio := &Audio{
		sources: map[*gorge.AudioSource]*Processor{},
	}

	g.PutProp(&Context{audio})
	g.Handle(audio)

	return nil
}*/

// Audio struct state for tracked playing audio clips, effects stack etc
type Audio struct {
	gorge    *gorge.Context
	oto      *oto.Context
	listener rListenerEntity
	// Player per Source
	sources  map[*gorge.AudioSource]*Processor
	entities []rSourceEntity
}

// HandleEvent implements the eventhandler.
func (s *Audio) HandleEvent(ee event.Event) {
	e, ok := ee.(gorge.EventAddEntity)
	if !ok {
		return
	}

	if ae, ok := e.Entity.(rListenerEntity); ok {
		s.listener = ae
	}
	if ae, ok := e.Entity.(rSourceEntity); ok {
		if s.oto == nil {
			// Lazy start
			// Need options here
			o, err := oto.NewContext(44100, 2, 2, 2048)
			if err != nil {
				s.gorge.Error(err)
				return
			}
			s.oto = o
		}
		s.entities = append(s.entities, ae)
		comp := ae.AudioSourceComponent()

		player := s.oto.NewPlayer()

		proc := &Processor{
			audio:      s,
			player:     player,
			positional: proc.NewPositional(player),
			source:     ae.AudioSourceComponent(),
			entity:     ae,
		}
		log.Println("Adding audio source", proc)
		// Create a processor
		s.sources[comp] = proc
		proc.Run()
	}
}

// Processor processes an audio source
type Processor struct {
	audio      *Audio
	positional *proc.Positional
	player     *oto.Player

	entity rSourceEntity
	source *gorge.AudioSource
	cur    int
}

func (p *Processor) run() {
	for {
		time.Sleep(1000 / 60 * time.Millisecond)

		if p.source.Clip == nil {
			continue
		}
		clip, ok := p.source.Clip.Resource().(*gorge.AudioClipData)
		if !ok {
			panic("panic loading audio clip")
		}
		p.cur = 0
		updates := p.source.Updates
		for {
			// There should some kind of lock for this
			// or we handle changes via Update and update it here
			if !p.source.Playing || p.source.Clip == nil {
				break
			}
			if p.cur >= len(clip.Data) {
				if !p.source.Loop {
					p.source.Playing = false
					break
				}
			}
			if updates != p.source.Updates {
				break
			}
			buf := clip.Data[p.cur:]
			if len(buf) > 100 {
				buf = buf[:100]
			}

			pos := m32.Vec3{}
			if e, ok := p.entity.(matrixer); ok {
				pos = e.M4().Col(3).Vec3()
			}
			if l, ok := p.audio.listener.(matrixer); ok {
				pos = l.M4().Inv().
					MulV4(pos.Vec4(1)).Vec3()
			}
			p.positional.Position = pos
			n, err := p.positional.Write(buf)
			if err != nil {
				break
			}
			p.cur += n
		}
	}
}

// Run the audio processor in the background.
func (p *Processor) Run() {
	go p.run()
}

type matrixer interface {
	M4() m32.Mat4
}

type rSourceEntity interface {
	AudioSourceComponent() *gorge.AudioSource
}

type rListenerEntity interface {
	AudioListenerComponent() *gorge.AudioListener
}
