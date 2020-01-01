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

// Asset common struct
type asset struct {
	release bool
	//xgpu    interface{}
}

func (a *asset) Release() {
	a.release = true
}

// XXX: Experimental store stuff on asset
// If performance gets affected we can use these
// to store gl ID's, it will affect texture fetch performance
// as we are rendering vbo's directly, texture is being fetched by
// gorge.Texture on renderer

/*func (a *asset) SetGPU(x interface{}) {
	a.xgpu = x
}
func (a *asset) GetGPU() interface{} {
	return a.xgpu
}*/
