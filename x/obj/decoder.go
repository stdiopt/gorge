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

// Package obj preliminary obj loader outside of the package as we probably
// will load more than 1 mesh, we don't have concept of "objects" with several meshes yet
//
// Entity
//   Transform
//   Mesh
//
//
package obj

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/m32"
)

type (
	vec3 = m32.Vec3
	vec2 = m32.Vec2
)

// we need asset loader

// Decoder Decodes an wavefront .obj file

type face []rawIndex

// This produces one or more meshes but lets try with one first

type rawIndex struct {
	//position,texture, normal
	indices [3]int
}

// Decode the obj file into some format
func Decode(rd io.Reader) (*gorge.Mesh, error) {
	mesh, err := readMesh(rd)
	if err != nil {
		return nil, err
	}

	return mesh, nil
}

type rawObj struct {
	vertices []vec3
	uvs      []vec2
	normals  []vec3
	faces    []face
}

func readMesh(rd io.Reader) (*gorge.Mesh, error) {
	s := bufio.NewScanner(rd)

	o := rawObj{
		vertices: []vec3{},
		uvs:      []vec2{},
		normals:  []vec3{},
		faces:    []face{},
	}
	line := 0
	oCount := 0
	//MainLoop:
	for s.Scan() {
		if err := s.Err(); err != nil {
			return nil, err
		}
		t := s.Text()
		line++

		t = strings.TrimSpace(t)
		if len(t) == 0 {
			continue
		}
		re := regexp.MustCompile(`\s+`)
		// Space split
		parts := re.Split(t, -1)
		strings.Split(t, " ")

		switch parts[0] {
		case "v":
			vert, err := getVec3(parts[1:])
			if err != nil {
				return nil, err
			}

			o.vertices = append(o.vertices, vert)
		case "vt":
			uv, err := getVec2(parts[1:])
			if err != nil {
				return nil, err
			}
			o.uvs = append(o.uvs, uv)
		case "vn":
			norm, err := getVec3(parts[1:])
			if err != nil {
				return nil, err
			}
			o.normals = append(o.normals, norm)
		case "f":
			fac := face{}
			for _, p := range parts[1:] {
				vpart := strings.Split(p, "/")
				ind := rawIndex{}
				for j, vp := range vpart {
					if len(vp) == 0 {
						continue
					}
					v, err := strconv.ParseInt(vp, 10, 64)
					if err != nil {
						return nil, err
					}
					ind.indices[j] = int(v)
				}
				fac = append(fac, ind)
			}
			o.faces = append(o.faces, fac)

			// Index per thing
		case "g": // Group?
		case "s": // Smoothing group
		case "o": // Object name
			oCount++
			if oCount > 2 {
				//break MainLoop
			}
		default:
		}
	}

	vertexRes := []gorge.VertexPTN{}
	vertexInd := []uint32{}

	// Good mapping
	vertexRef := map[string]uint32{}

	quad := 0
	// Convert from raw to Mesh
	for _, face := range o.faces {
		iface := []uint32{}
		//for i := len(face) - 1; i >= 0; i-- {
		//fi := face[i]
		for i, fi := range face {
			key := fmt.Sprintf("%d/%d/%d", fi.indices[0], fi.indices[1], fi.indices[2])

			if i >= 3 {
				quad++
				// We need to add first and penultimate indice
				v1, v2 := iface[0], iface[2]
				iface = append(iface, v1)
				iface = append(iface, v2)
			}
			// Check if we already have a thing
			if v, ok := vertexRef[key]; ok {
				iface = append(iface, v) // add index then
				continue
			}
			// If doesn't exists we get the vertex info and create a new vertex

			nv := gorge.VertexPTN{}
			if fi.indices[0] > 0 {
				nv.Pos = o.vertices[fi.indices[0]-1]
				nv.Pos[2] *= -1 // Invert Z
			}
			if fi.indices[1] > 0 {
				nv.Tex = o.uvs[fi.indices[1]-1]
				nv.Tex[1] *= -1
			}
			if fi.indices[2] > 0 {
				nv.Normal = o.normals[fi.indices[2]-1]
				//nv.Normal[0] *= -1
				//nv.Normal[1] *= -1
				nv.Normal[2] *= -1
			}

			// if the index > 3

			rind := uint32(len(vertexRes))
			vertexRes = append(vertexRes, nv)
			vertexRef[key] = rind
			iface = append(iface, rind)
		}
		vertexInd = append(vertexInd, iface...)
	}

	ptn := &gorge.MeshPTN{
		Vertices: vertexRes,
		Indices:  vertexInd,
	}
	return &gorge.Mesh{MeshLoader: ptn}, nil
}

func getVec3(parts []string) (m32.Vec3, error) {
	var ret m32.Vec3
	for i := 0; i < 3; i++ {
		s := parts[i]
		if err := parse(s, &ret[i]); err != nil {
			return m32.Vec3{}, err
		}
	}
	return ret, nil
}
func getVec2(parts []string) (m32.Vec2, error) {
	var ret m32.Vec2
	for i := 0; i < 2; i++ {
		s := parts[i]
		if err := parse(s, &ret[i]); err != nil {
			return m32.Vec2{}, err
		}
	}
	return ret, nil
}

func parse(s string, v interface{}) error {
	var err error
	switch v := v.(type) {
	case *int:
		var r int64
		r, err = strconv.ParseInt(s, 10, 64)
		*v = int(r)
	case *float32:
		var r float64
		r, err = strconv.ParseFloat(s, 64)
		*v = float32(r)
	}
	return err
}
