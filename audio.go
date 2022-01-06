package gorge

// AudioResourcer interface to return an audio resource.
type AudioResourcer interface {
	Resource() AudioResource
}

// AudioResource audio resource.
type AudioResource interface {
	isAudio()
}

// AudioFormat audio format for clipData.
type AudioFormat int

// AudioClip is the resource controller for audio (similar to material, texture, mesh).
type AudioClip struct {
	Resourcer AudioResourcer
}

// NewAudioClip creates a New audio clip based on resource.
func NewAudioClip(ref AudioResourcer) *AudioClip {
	return &AudioClip{ref}
}

// Resource returns the current resource from resourcer.
func (a *AudioClip) Resource() AudioResource {
	if a.Resourcer == nil {
		return nil
	}
	return a.Resourcer.Resource()
}

// AudioSource component.
type AudioSource struct {
	Playing bool
	Loop    bool
	Clip    *AudioClip
	Updates int
}

// AudioSourceComponent implements the component
func (a *AudioSource) AudioSourceComponent() *AudioSource { return a }

// Play sets the play state to playing
func (a *AudioSource) Play(c *AudioClip) {
	a.Updates++
	a.Clip = c
	a.Playing = true
}

// AudioListener is where the audio will be listened usually set on cameras
type AudioListener struct{}

// AudioListenerComponent implements the component
func (a *AudioListener) AudioListenerComponent() *AudioListener { return a }

// AudioClipData base audio data.
type AudioClipData struct {
	Format  AudioFormat
	Data    []byte
	Updates int
}

// Resource implements the AudioResourcer interface.
func (d *AudioClipData) Resource() AudioResource { return d }

func (AudioClipData) isAudio() {}
