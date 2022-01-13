package resource

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/hajimehoshi/go-mp3"
	"github.com/stdiopt/gorge"
)

func init() {
	Register((*gorge.AudioClipData)(nil), ".mp3", audioClipDataLoader)
	Register((*gorge.AudioClip)(nil), ".mp3", audioClipLoader)
}

func audioClipLoader(res *Context, v any, name string, opts ...any) error {
	clip := v.(*gorge.AudioClip)

	var clipData gorge.AudioClipData
	if err := audioClipDataLoader(res, &clipData, name, opts...); err != nil {
		return err
	}

	clip.Resourcer = &clipData

	return nil
}

func audioClipDataLoader(res *Context, v any, name string, _ ...any) error {
	clipData := v.(*gorge.AudioClipData)

	rd, err := res.Open(name)
	if err != nil {
		return fmt.Errorf("error opening audio clip: %w", err)
	}

	ext := filepath.Ext(name)
	switch ext {
	case ".mp3":
		dec, err := mp3.NewDecoder(rd)
		if err != nil {
			return err
		}
		data, err := ioutil.ReadAll(dec)
		if err != nil {
			return err
		}
		// clipData.Format = ....
		clipData.Data = data
	default:
		return fmt.Errorf("unknown audioClip type: %s", ext)
	}
	return nil
}
