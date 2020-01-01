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

package resource

import (
	"fmt"
	"image"
	"io"
	"io/ioutil"

	// Image formats
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	glog "log"

	"github.com/stdiopt/gorge"
)

var (
	log = glog.New(glog.Writer(), "(resource) ", 0)
)

// Loader interface specifies the loader methods used to fetch assets
type Loader interface {
	// Allow Put too for memory resources
	Open(name string) (io.ReadCloser, error)
}

// Manager just the initial manager which contains the Loader
type Manager struct {
	Loader
	gorge    *gorge.Gorge
	internal map[string]interface{}

	loaders map[interface{}]interface{}
}

// NewSystem returns a new manager
func NewSystem(l Loader) gorge.SystemFunc {
	return func(g *gorge.Gorge) {
		m := &Manager{
			gorge:    g,
			Loader:   l,
			internal: map[string]interface{}{},
		}
		g.Persist(m)

		g.Persist(gorge.AssetBundleProvider{
			Bundle: BundleFunc(m),
		})
	}
}

///////////////////////////////////////////////////////////////////////////////
// specialized methods
///////////////////////

// LoadImage get an asset as an image
func (m *Manager) LoadImage(name string) (image.Image, error) {
	rd, err := m.Loader.Open(name)
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
func (m *Manager) LoadBytes(name string) ([]byte, error) {
	rd, err := m.Loader.Open(name)
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
func (m *Manager) LoadString(name string) (string, error) {
	data, err := m.LoadBytes(name)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (m *Manager) Error(err error) {
	m.gorge.Persist(gorge.ErrorEvent{Err: err})
}

// Asset returns a resource reference that can be later open
func (m *Manager) Asset(name string) Opener {
	return Opener{name, m}
}

// Opener holds a reference to the system and to the name for later opening
type Opener struct {
	name    string
	manager *Manager
}

// Open the asset using the loader
func (o Opener) Open() (io.ReadCloser, error) {
	return o.manager.Open(o.name)
}
