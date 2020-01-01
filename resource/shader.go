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
	"strings"

	"github.com/stdiopt/gorge"
)

type shaderLoader struct {
	name   string
	loader func() *gorge.ShaderData
}

func (l *shaderLoader) Data() *gorge.ShaderData {
	return l.loader()
}

// Material gets a material with a specific shader
/*func (a *Asset) Material(names ...string) *gorge.Material {
	return &gorge.Material{
		ShaderLoader: a.shaderLoader(names...),
	}
}*/

type shaderLoaderKey struct {
	vertName string
	fragName string
}

func (m *Manager) shaderLoader(names ...string) *shaderLoader {
	var vertName, fragName string

	switch len(names) {
	case 0:
		panic("name is required")
	case 1:
		vertName = names[0] + ".vert"
		fragName = names[0] + ".frag"
	case 2:
		vertName = names[0]
		fragName = names[1]
	}

	if m.loaders == nil {
		m.loaders = map[interface{}]interface{}{}
	}
	k := shaderLoaderKey{vertName, fragName}

	if loader, ok := m.loaders[k]; ok {
		return loader.(*shaderLoader)
	}

	loaderFunc := func() *gorge.ShaderData {
		vertSrc, err := m.LoadString(vertName)
		if err != nil {
			m.Error(err)
			return nil
		}
		fragSrc, err := m.LoadString(fragName)
		if err != nil {
			m.Error(err)
			return nil
		}
		return &gorge.ShaderData{
			VertSrc: vertSrc,
			FragSrc: fragSrc,
		}
	}
	loader := &shaderLoader{
		name:   strings.Join(names, ","),
		loader: loaderFunc,
	}
	m.loaders[k] = loader
	return loader
}
