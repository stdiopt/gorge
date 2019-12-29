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

package asset

import (
	"fmt"
	"image"
	"io"
	"io/ioutil"

	// Image formats
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/stdiopt/gorge"
)

// Loader interface specifies the loader methods used to fetch assets
type Loader interface {
	// Allow Put too for memory resources
	Load(name string) (io.ReadCloser, error)
}

// System just the initial manager which contains the Loader
type System struct {
	Loader
	gorge    *gorge.Gorge
	internal map[string]interface{}
}

// NewSystem returns a new manager
func NewSystem(l Loader) gorge.SystemFunc {
	return func(g *gorge.Gorge) {
		s := &System{
			gorge:    g,
			Loader:   l,
			internal: map[string]interface{}{},
		}
		g.Persist(s)
	}
}

// FromECS gets the asset manager from gorge
func FromECS(g *gorge.Gorge) *System {
	var ret *System
	g.Query(func(s *System) { ret = s })
	return ret
}

// Put something in memory
func (s *System) Put(name string, v interface{}) {
	s.internal[name] = v
}

///////////////////////////////////////////////////////////////////////////////
// specialized methods
///////////////////////

// LoadImage get an asset as an image
func (s *System) LoadImage(name string) (image.Image, error) {
	if v, ok := s.internal[name].(image.Image); ok {
		return v, nil
	}

	rd, err := s.Loader.Load(name)
	if err != nil {
		return nil, fmt.Errorf("error getting image: %w", err)
	}

	img, _, err := image.Decode(rd)
	if err != nil {
		return nil, fmt.Errorf("error deconding image: %w", err)
	}
	return img, nil
}

// LoadBytes returns the asset as string
func (s *System) LoadBytes(name string) ([]byte, error) {
	rd, err := s.Loader.Load(name)
	if err != nil {
		return nil, err
	}
	defer rd.Close()

	data, err := ioutil.ReadAll(rd)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// LoadString returns asset as a string
func (s *System) LoadString(name string) (string, error) {
	data, err := s.LoadBytes(name)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s *System) Error(err error) {
	s.gorge.Persist(gorge.ErrorEvent(err))
}

// Opener holds a reference to the system and to the name for later opening
type Opener struct {
	name string
	s    *System
}

// Open the asset using the loader
func (a Opener) Open() (io.ReadCloser, error) {
	return a.s.Load(a.name)
}

// Asset returns a resource reference that can be reopened (reused)
func (s *System) Asset(name string) Opener {
	return Opener{name, s}
}
