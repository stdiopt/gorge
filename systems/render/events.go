package render

import "time"

// EventStat track gpu resources for debugging
type EventStat struct {
	VAOs           int
	VBOs           int
	Textures       int
	Shaders        int
	Instances      int
	Buffers        int
	RenderDuration time.Duration
	DrawCalls      int
}
