// Copyright 2020 Luis Figueiredo
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

package gorge

// TODO:
// Issues: we cannot import assets as it is because of cyclic dependencies
// but we could create some kind of placeholder to tell asset pkg to load
// the stuff

// SceneStartEvent starts a scene and pass scene to know which scene is tarting
type SceneStartEvent *Scene

// Scene creates a new Scene from gorge
func (g *Gorge) Scene(fns ...SceneFunc) *Scene {

	// Dependency checker
	var b *AssetBundleProvider
	g.Query(func(v AssetBundleProvider) { b = &v })
	// Creates a new scene
	s := &Scene{
		assets: b.Bundle(),
		gorge:  g,
	}

	// Calls scene initializrs
	for _, fn := range fns {
		fn(s)
	}

	return s
}

// LoadScene does:
// - Load scene assets
// - Add scene entities
func (g *Gorge) LoadScene(s *Scene) {
	// Ref assetBundle
	g.Trigger(LoadBundleEvent{s.assets.Assets()})
	ents := []Entity{}
	s.entities.Range(func(v interface{}) bool {
		ents = append(ents, v)
		return true
	})

	g.Trigger(EntitiesAddEvent(ents))
}

// StartScene loads and start a scene
func (g *Gorge) StartScene(s *Scene) {
	g.LoadScene(s)
	s.Start()
}

// SceneFunc scene initializer func
type SceneFunc func(*Scene)

// Scene contains stuff
type Scene struct {
	Messaging
	gorge    *Gorge
	assets   AssetBundler
	entities SetList
}

// Start ...
func (s *Scene) Start() {
	s.gorge.Link(&s.Messaging)
	//
	s.Trigger(SceneStartEvent(s))
}

// Assets returns the current scene asset bundler
func (s *Scene) Assets() AssetBundler {
	return s.assets
}

// AddEntity to scene
func (s *Scene) AddEntity(ents ...Entity) {

	// XXX: Edge case?
	// This will solve assets and put it on scene bundle so it can be unloaded
	type renderable interface {
		RenderableComponent() *Renderable
	}

	for _, e := range ents {
		if !s.entities.Add(e) {
			continue
		}
		// Solve asset
		v, ok := e.(renderable)
		if !ok {
			continue
		}
		rend := v.RenderableComponent()
		s.Assets().AddAsset(rend.Mesh)
		for _, t := range rend.Material.Textures {
			s.Assets().AddAsset(t)
		}
	}
	s.Trigger(EntitiesAddEvent(ents))
}

// RemoveEntity ..
func (s *Scene) RemoveEntity(ents ...Entity) {
	// Remove locally from scene
	for _, e := range ents {
		s.entities.Remove(e)
	}

	s.Trigger(EntitiesRemoveEvent(ents))
}
