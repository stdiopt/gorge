package gltf

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/systems/resource"
)

func init() {
	resource.Register(&GLTF{}, ".gltf", gltfLoader)
	resource.Register(&GLTF{}, ".glb", glbLoader)
}

func gltfLoader(res *resource.Context, v interface{}, name string, _ ...interface{}) error {
	gOut := v.(*GLTF)

	rd, err := res.Open(name)
	if err != nil {
		return err
	}
	defer rd.Close() // nolint: errcheck

	var root Doc

	if err := json.NewDecoder(rd).Decode(&root); err != nil {
		return err
	}

	basePath := filepath.Dir(name)
	// Load this right away?
	for _, b := range root.Buffers {
		err := func() error {
			var bufRd io.Reader
			switch {
			case strings.HasPrefix(b.URI, "data:"):
				b64data := b.URI[strings.Index(b.URI, ",")+1:]
				bufRd = base64.NewDecoder(base64.StdEncoding, strings.NewReader(b64data))
			case b.URI != "":
				rd, err := res.Open(filepath.Join(basePath, b.URI))
				if err != nil {
					return err
				}
				bufRd = rd
			}
			data, err := ioutil.ReadAll(bufRd)
			if err != nil {
				return err
			}
			b.RawData = data
			return nil
		}()
		if err != nil {
			return err
		}
	}
	root.BasePath = basePath

	*gOut = *create(res.Gorge(), &root)
	return nil
}

func glbLoader(res *resource.Context, v interface{}, name string, _ ...interface{}) error {
	gOut := v.(*GLTF)

	rd, err := res.Open(name)
	if err != nil {
		return err
	}
	defer rd.Close() // nolint: errcheck

	buf, err := ioutil.ReadAll(rd)
	if err != nil {
		return err
	}

	var root Doc
	nChunk := buf[12:]                             // Skip magic
	chunkLen := binary.LittleEndian.Uint32(nChunk) // first part of chunk is size
	jsonChunk := nChunk[8:][:chunkLen]             // map jsonChunk

	nChunk = nChunk[8+chunkLen:]                 // skip jsonChunk
	bufLen := binary.LittleEndian.Uint32(nChunk) // read buffer chunkSize
	bufChunk := nChunk[8:][:bufLen]              // map buffer Chunk

	jsonReader := bytes.NewReader(jsonChunk)
	if err := json.NewDecoder(jsonReader).Decode(&root); err != nil {
		return err
	}

	// binary exclusive
	root.Buffers[0].RawData = bufChunk
	root.BasePath = filepath.Dir(name)

	*gOut = *create(res.Gorge(), &root)
	return nil
}

// Do not load directly? we should lazy load these? based on image stuff
// Although it could cause a double load?
// Should be on the other side
func loadImage(res *resource.Context, root *Doc, m *Image) (*gorge.TextureData, error) {
	switch {
	case strings.HasPrefix(m.URI, "data:"):
		b64data := m.URI[strings.Index(m.URI, ",")+1:]
		return resource.ReadTexture(
			base64.NewDecoder(base64.StdEncoding, strings.NewReader(b64data)),
		)
	case m.URI != "":
		var texData gorge.TextureData
		if err := res.Load(&texData, filepath.Join(root.BasePath, m.URI)); err != nil {
			return nil, err
		}
		return &texData, nil
	case m.BufferView != nil:
		bv := root.BufferViews[*m.BufferView]
		// buffer should be 0 here
		buf := root.Buffers[bv.Buffer].RawData
		viewBuf := buf[bv.ByteOffset:][:bv.ByteLength]
		return resource.ReadTexture(bytes.NewReader(viewBuf))
	default:
		return nil, errors.New("Image must have either URI or BufferView")
	}
}
