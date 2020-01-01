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
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/x/obj"
)

type meshLoader struct {
	source string
	loader func() *gorge.MeshData
}

func (l *meshLoader) Data() *gorge.MeshData { return l.loader() }

// Mesh loads a mesh (obj only now)
/*func (a *Asset) Mesh(name string) *gorge.Mesh {
	return gorge.NewMesh(a.meshLoader(name))
}*/

type meshLoaderKey struct{ name string }

func (m *Manager) meshLoader(name string) *meshLoader {
	if m.loaders == nil {
		m.loaders = map[interface{}]interface{}{}
	}
	k := meshLoaderKey{name}
	if loader, ok := m.loaders[k]; ok {
		return loader.(*meshLoader)
	}

	loaderFunc := func() *gorge.MeshData {
		rd, err := m.Open(name)
		if err != nil {
			m.Error(err)
			return nil
		}
		data, err := obj.Decode(rd)
		return data.Loader().Data()
	}

	loader := &meshLoader{
		source: name,
		loader: loaderFunc,
	}
	m.loaders[k] = loader
	return loader
}
