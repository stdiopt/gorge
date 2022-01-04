package gltf

import (
	"fmt"
	"log"
	"unsafe"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/anim"
	"github.com/stdiopt/gorge/gorgeutil"
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/static"
	"github.com/stdiopt/gorge/systems/resource"
)

// GLTF gorge based gltf doc
// This loads all stuff
// Need to figure out something to lazy load, free memory and whatnots
type GLTF struct {
	Textures   []*gorge.Texture
	Materials  []*gorge.Material
	Scenes     []*GScene
	Nodes      []*GNode
	Meshes     []*GMesh
	Skins      []*GSkin
	Animations []*anim.Animation

	updateFn []func(dt float32)
}

type gltfCreator struct {
	gorge *gorge.Context
	doc   *Doc

	texRef  map[*Image]*gorge.TextureData
	primRef map[*MeshPrimitive]*gorge.MeshData

	Textures   []*gorge.Texture
	Materials  []*gorge.Material
	Scenes     []*GScene
	Nodes      []*GNode
	Meshes     []*GMesh
	Skins      []*GSkin
	Animations []*anim.Animation

	updateFn []func(dt float32)
}

// gltf Model into gorge based stuff
func create(g *gorge.Context, doc *Doc) *GLTF {
	c := gltfCreator{
		gorge:   g,
		doc:     doc,
		texRef:  map[*Image]*gorge.TextureData{},
		primRef: map[*MeshPrimitive]*gorge.MeshData{},
	}
	c.processTextures()
	c.processMaterials()
	c.processMeshes()
	c.processSkins()
	c.processNodes()
	c.processScenes()
	c.processAnimations()

	return &GLTF{
		Textures:   c.Textures,
		Materials:  c.Materials,
		Scenes:     c.Scenes,
		Nodes:      c.Nodes,
		Meshes:     c.Meshes,
		Skins:      c.Skins,
		Animations: c.Animations,
		updateFn:   c.updateFn,
	}
}

// ReleaseRawData releases the raw data from memory leaving the gpu ref only
func (r *GLTF) ReleaseRawData(g *gorge.Context) {
	for _, t := range r.Textures {
		t.ReleaseData(g)
		// g.ResourceReleaseData(t)
	}
	for _, m := range r.Meshes {
		for _, p := range m.primitives {
			p.ReleaseData(g)
		}
	}
}

// UpdateDelta to be manually called to trigger animations, morphs etc.
func (r *GLTF) UpdateDelta(dt float32) {
	for _, fn := range r.updateFn {
		fn(dt)
	}
}

func (c *gltfCreator) processTextures() {
	textures := []*gorge.Texture{}
	for _, t := range c.doc.Textures {
		textures = append(textures, c.getGTexture(t))
	}
	c.Textures = textures
}

func defMaterial() *gorge.Material {
	mat := gorge.NewShaderMaterial(static.Shaders.DefaultNew)
	mat.Define("USE_HDR", "USE_IBL")

	mat.Set("u_MipCount", 5)
	mat.Set("u_EmissiveFactor", m32.Vec3{0, 0, 0})
	mat.Set("u_AlphaCutoff", float32(0.5))
	mat.Set("u_Exposure", float32(1))

	// This is relative to MATERIAL_METALLICROUGHNESS probably shouldn't be here?
	// Maybe add some flag upthere
	mat.Set("u_BaseColorFactor", m32.Vec4{1, 1, 1, 1})
	mat.Set("u_MetallicFactor", float32(1))
	mat.Set("u_RoughnessFactor", float32(1))

	return mat
}

func (c *gltfCreator) getMaterial(tfMat *Material) *gorge.Material {
	mat := defMaterial()

	alphaModeDef := "ALPHAMODE_OPAQUE"
	if v := tfMat.AlphaMode; v != nil {
		switch *v {
		case "MASK":
			alphaModeDef = "ALPHAMODE_MASK"
		case "BLEND":
			mat.Queue = 10
			alphaModeDef = "ALPHAMODE_BLEND"
		}
	}
	mat.Define(alphaModeDef)

	if v := tfMat.AlphaCutoff; v != nil {
		mat.SetFloat32("u_AlphaCutoff", *v)
	}
	if v := tfMat.EmissiveFactor; v != nil {
		mat.Set("u_EmissiveFactor", m32.Vec3(*v))
	}

	if tfMat.DoubleSided != nil {
		mat.DoubleSided = *tfMat.DoubleSided
	}

	// could it be also spec and gloss at the same time?
	if pbr := tfMat.PBRMetallicRoughness; pbr != nil {
		mat.Define("MATERIAL_METALLICROUGHNESS")

		if v := pbr.BaseColorFactor; v != nil {
			mat.Set("u_BaseColorFactor", m32.Vec4(*v))
		}
		if t := pbr.BaseColorTexture; t != nil {
			tex := c.Textures[t.Index]
			mat.Define("HAS_BASE_COLOR_MAP")
			mat.Set("u_BaseColorUVSet", t.TexCoord)
			mat.SetTexture("u_BaseColorSampler", tex)
		}

		if v := pbr.MetallicFactor; v != nil {
			mat.Set("u_MetallicFactor", *v)
		}
		if v := pbr.RoughnessFactor; v != nil {
			mat.Set("u_RoughnessFactor", *v)
		}
		if t := pbr.MetallicRoughnessTexture; t != nil {
			tex := c.Textures[t.Index]
			mat.Define("HAS_METALLIC_ROUGHNESS_MAP")
			mat.Set("u_MetallicRoughnessUVSet", t.TexCoord)
			mat.SetTexture("u_MetallicRoughnessSampler", tex)
		}
	}

	if t := tfMat.OcclusionTexture; t != nil { // aoMap
		tex := c.Textures[t.Index]
		mat.Define("HAS_OCCLUSION_MAP")
		mat.Set("u_OcclusionUVSet", t.TexCoord)
		mat.SetTexture("u_OcclusionSampler", tex)
		mat.SetFloat32("u_OcclusionStrength", 1)
		if v := t.Strength; v != nil {
			mat.Set("u_OcclusionStrength", *v)
		}
	}

	if t := tfMat.NormalTexture; t != nil { // normalMap
		tex := c.Textures[t.Index]
		mat.Define("HAS_NORMAL_MAP")
		mat.SetTexture("u_NormalSampler", tex)
		mat.Set("u_NormalUVSet", t.TexCoord)
		mat.SetFloat32("u_NormalScale", 1)
		if v := t.Scale; v != nil {
			mat.Set("u_NormalScale", v)
		}
	}

	if t := tfMat.EmissiveTexture; t != nil {
		tex := c.Textures[t.Index]
		mat.Define("HAS_EMISSIVE_MAP")
		mat.Set("u_EmissiveUVSet", t.TexCoord)
		mat.SetTexture("u_EmissiveSampler", tex)
	}

	if tfMat.Extensions == nil {
		return mat
	}

	if tfMat.Extensions.Unlit != nil {
		mat.Define("MATERIAL_UNLIT")
	}

	if m := tfMat.Extensions.Clearcoat; m != nil {
		mat.Define("MATERIAL_CLEARCOAT")
		mat.SetFloat32("u_ClearcoatFactor", 0)
		mat.SetFloat32("u_ClearcoatRoughnessFactor", 0)
		if v := m.ClearcoatFactor; v != nil {
			mat.Set("u_ClearcoatFactor", *v)
		}
		if t := m.ClearcoatTexture; t != nil {
			tex := c.Textures[t.Index]
			mat.Define("HAS_CLEARCOAT_TEXTURE_MAP")
			mat.Set("u_ClearcoatUVSet", t.TexCoord)
			mat.SetTexture("u_ClearcoatSampler", tex)
		}

		if v := m.ClearcoatRoughnessFactor; v != nil {
			mat.Set("u_ClearcoatRoughnessFactor", *v)
		}
		if t := m.ClearcoatRoughnessTexture; t != nil {
			tex := c.Textures[t.Index]
			mat.Define("HAS_CLEARCOAT_ROUGHNESS_MAP")
			mat.Set("u_ClearcoatRoughnessUVSet", t.TexCoord)
			mat.SetTexture("u_ClearcoatRoughnessSampler", tex)
		}

		if t := m.ClearcoatNormalTexture; t != nil {
			tex := c.Textures[t.Index]
			mat.Define("HAS_CLEARCOAT_NORMAL_MAP")
			mat.Set("u_ClearcoatNormalUVSet", t.TexCoord)
			mat.SetTexture("u_ClearcoatNormalSampler", tex)
			if v := t.Scale; v != nil {
				mat.Set("u_ClearcoatNormalScale", *v) // Not in use
			}
		}
	}
	return mat
}

func (c *gltfCreator) processMaterials() {
	materials := []*gorge.Material{}
	for _, tfMat := range c.doc.Materials {
		materials = append(materials, c.getMaterial(tfMat))
	}
	c.Materials = materials
}

func (c *gltfCreator) processMeshes() {
	renderables := []*GMesh{}
	for _, m := range c.doc.Meshes {
		rr := []*gorge.RenderableComponent{}
		// meshTransform := gorge.TransformIdent()
		for _, prim := range m.Primitives {
			mesh := c.getGPrimitive(prim)
			var mat *gorge.Material
			if prim.Material != nil {
				mat = c.Materials[*prim.Material]
			} else {
				mat = defMaterial()
				mat.Define("ALPHAMODE_OPAQUE", "MATERIAL_METALLICROUGHNESS")
			}

			if len(m.Weights) > 0 {
				mesh.Define("USE_MORPHING", fmt.Sprintf("WEIGHT_COUNT %d", len(m.Weights)))
			}

			rr = append(rr, gorge.NewRenderableComponent(mesh, mat))
		}

		renderables = append(renderables, &GMesh{
			Name:       m.Name,
			Weights:    m.Weights,
			primitives: rr,
		})
	}
	c.Meshes = renderables
}

func (c *gltfCreator) processSkins() {
	skins := []*GSkin{}
	for _, s := range c.doc.Skins {
		var matrices []m32.Mat4
		if s.InverseBindMatrices != nil {
			matrices = bufMat4Slice(acBuf(c.doc.AccessorBuffer(*s.InverseBindMatrices)))
		}
		skins = append(skins, &GSkin{
			Matrices: matrices,
			Joints:   s.Joints,
		})
	}
	c.Skins = skins
}

func (c *gltfCreator) processNodes() {
	nodes := []*GNode{}
	for _, n := range c.doc.Nodes {
		node := &GNode{
			TransformComponent: gorge.NewTransformComponent(),
		}
		if n.Skin != nil {
			node.skin = c.Skins[*n.Skin]
		}

		if n.Mesh != nil {
			node.mesh = c.Meshes[*n.Mesh]
			// Create Primitives here
			for _, r := range node.mesh.primitives {
				// Clone mesh too
				primMesh := r.Clone()
				primMesh.Define("HAS_SINGLE_INSTANCE")
				p := gorgeutil.NewRenderable(primMesh, r.Material)
				p.SetParent(node)
				node.entities = append(node.entities, p)
				if node.skin == nil {
					continue
				}
				// Use skinning
				primMesh.Define(
					"USE_SKINNING",
					fmt.Sprintf("JOINT_COUNT %d", len(node.skin.Matrices)),
				)
				fn := func(_ float32) {
					for i, ni := range node.skin.Joints {
						m := node.Mat4().Inv()
						m = m.Mul(nodes[ni].Mat4())
						m = m.Mul(node.skin.Matrices[i])
						// This should be set somewhere automatically within mesh
						primMesh.Set(fmt.Sprintf("u_jointMatrix[%d]", i), m)
						primMesh.Set(fmt.Sprintf("u_jointNormalMatrix[%d]", i), m.Inv().Transpose())
						primMesh.Update()
					}
				}
				c.updateFn = append(c.updateFn, fn)

			}
		}

		if n.Matrix != nil {
			// https://answers.unity.com/questions/402280/how-to-decompose-a-trs-matrix.html
			node.Transform().SetMat4Decompose(m32.Mat4(*n.Matrix))
		}

		if n.Rotation != nil {
			r := *n.Rotation
			node.Transform().Rotation = m32.Quat(r)
		}

		if n.Translation != nil {
			node.Transform().Position = m32.Vec3(*n.Translation)
		}

		if n.Scale != nil {
			node.Transform().Scale = m32.Vec3(*n.Scale)
		}

		nodes = append(nodes, node)

	}

	// Solve children stuff
	for i, n := range nodes {
		for _, ci := range c.doc.Nodes[i].Children {
			child := nodes[ci]
			child.SetParent(n)
			n.children = append(n.children, child)
		}
	}
	c.Nodes = nodes
}

func (c *gltfCreator) processScenes() {
	scenes := []*GScene{}
	for _, s := range c.doc.Scenes {
		scene := &GScene{
			TransformComponent: gorge.NewTransformComponent(),
		}
		for _, ni := range s.Nodes {
			c.Nodes[ni].SetParent(scene)
			scene.Nodes = append(scene.Nodes, c.Nodes[ni])
		}

		scenes = append(scenes, scene)
	}
	c.Scenes = scenes
}

func (c *gltfCreator) processAnimations() {
	animations := []*anim.Animation{}
	for _, a := range c.doc.Animations {
		animations = append(animations, c.getGAnimation(a))
	}
	c.Animations = animations
}

func (c *gltfCreator) getGAnimation(a *Animation) *anim.Animation {
	gAnim := &anim.Animation{}
	gAnim.SetLoop(anim.LoopAlways)
	for _, ch := range a.Channels {

		s := a.Samplers[ch.Sampler]
		keys := bufF32Slice(acBuf(c.doc.AccessorBuffer(s.Input)))

		off := 0
		ds := 1
		if s.Interpolation == "CUBICSPLINE" {
			off = 1 // we fetch the key from center
			log.Println("ANIMATION: Cubic spline")
			ds = 3
		}

		var track *anim.Channel
		var val []interface{}
		// We have to manually add node as we don't have it in gorge stuff
		// Track translation for now
		targetNode := c.Nodes[ch.Target.Node]
		switch ch.Target.Path {
		case "translation":
			data := bufVec3Slice(acBuf(c.doc.AccessorBuffer(s.Output)))
			// XXX: debug condition Remove when done
			if s.Interpolation == "CUBICSPLINE" {
				for i, k := range keys {
					log.Printf("Key info: %v -> %v,%v,%v",
						k,
						data[i*ds],
						data[i*ds+1],
						data[i*ds+2],
					)
				}
			}
			for i := range keys {
				val = append(val, data[i*ds+off])
			}
			track = gAnim.Channel(anim.Vec3(&targetNode.Position))
		case "rotation":
			data := bufVec4Slice(acBuf(c.doc.AccessorBuffer(s.Output)))

			track = gAnim.Channel(anim.Quat(&targetNode.Rotation))
			for i := range keys {
				v := data[i*ds+off]
				// q := m32.Q{W: v[3], V: mgl32.Vec3{v[0], v[1], v[2]}}
				// val = append(val, m32.Q(q))
				val = append(val, m32.Quat(v))
			}
		case "scale":
			data := bufVec3Slice(acBuf(c.doc.AccessorBuffer(s.Output)))

			track = gAnim.Channel(anim.Vec3(&targetNode.Scale))
			for i := range keys {
				val = append(val, data[i*ds+off])
			}
		case "weights":
			data := bufF32Slice(acBuf(c.doc.AccessorBuffer(s.Output)))

			wlen := len(targetNode.mesh.Weights)
			weightProps := make([]string, wlen)
			for i := 0; i < wlen; i++ {
				weightProps[i] = fmt.Sprintf("u_morphWeights[%d]", i)
			}

			track = gAnim.Channel(anim.InterpolatorFunc(func(a, b interface{}, dt float32) {
				va, vb := a.([]float32), b.([]float32)
				for i := range va {
					v := m32.Lerp(va[i], vb[i], dt) // Might be different according to interpolator
					// Set in every entity?
					for _, e := range targetNode.entities {
						p := e.(*gorgeutil.Renderable)
						p.Mesh.Set(weightProps[i], v)
					}
				}
			}))

			for i := range keys {
				kd := make([]float32, wlen)
				off := i*ds*wlen + off
				end := off + wlen
				copy(kd, data[off:end])
				val = append(val, kd)
			}
		}
		for i, k := range keys {
			kk := track.SetKey(k, val[i])
			switch s.Interpolation {
			case "LINEAR":
			case "STEP":
				kk.SetEase(anim.Step)
			case "CUBICSPLINE":
				log.Println("Cubic spline: not supported yet")
			}
		}

	}
	// Just mark as started, do not actually start animating
	gAnim.Start()
	return gAnim
}

// Primitive is a gorge Mesh
func (c *gltfCreator) getGPrimitive(prim *MeshPrimitive) *gorge.Mesh {
	if ref, ok := c.primRef[prim]; ok {
		return gorge.NewMesh(ref)
	}

	type elem struct {
		data []float32 // this is the vertices data type anyway
		sz   int       // size in floats
	}
	elems := []elem{}
	var format gorge.VertexFormat

	attrs := map[string]gorge.VertexFormatAttrib{
		"POSITION":   gorge.VertexAttrib(3, "a_Position", "HAS_POSITION"),
		"NORMAL":     gorge.VertexAttrib(3, "a_Normal", "HAS_NORMALS"),
		"TEXCOORD_0": gorge.VertexAttrib(2, "a_UV1", "HAS_UV_SET1"),
		"TEXCOORD_1": gorge.VertexAttrib(2, "a_UV2", "HAS_UV_SET2"),
		"TANGENT":    gorge.VertexAttrib(4, "a_Tangent", "HAS_TANGENTS"),
		"COLOR_0":    gorge.VertexAttrib(4, "a_Color", "HAS_VERTEX_COLOR_VEC4"),
		"JOINTS_0":   gorge.VertexAttrib(4, "a_Joint1", "HAS_JOINT_SET1"),
		"JOINTS_1":   gorge.VertexAttrib(4, "a_Joint2", "HAS_JOINT_SET2"),
		"WEIGHTS_0":  gorge.VertexAttrib(4, "a_Weight1", "HAS_WEIGHT_SET1"),
		"WEIGHTS_1":  gorge.VertexAttrib(4, "a_Weight2", "HAS_WEIGHT_SET2"),

		"TARGET_POSITION": gorge.VertexAttrib(3, "a_Target_Position", "HAS_TARGET_POSITION"),
		"TARGET_NORMAL":   gorge.VertexAttrib(3, "a_Target_Normal", "HAS_TARGET_NORMAL"),
		"TARGET_TANGENT":  gorge.VertexAttrib(3, "a_Target_Tangent", "HAS_TARGET_TANGENT"),
	}
	for a, ai := range prim.Attributes {
		attrib, ok := attrs[a]
		if !ok {
			log.Printf("attr: %v not supported yet", a)
			continue
		}
		buf, sz, ty := c.doc.AccessorBuffer(ai)

		if a == "COLOR_0" && sz == 3*4 {
			log.Println("Attrib as color3")
			attrib.Define = "HAS_VERTEX_COLOR_VEC3"
		}

		var data []float32
		switch ty {
		case ComponentUByte:
			for _, v := range buf {
				data = append(data, float32(v))
			}
		// This is an extra case for joints that I've seen so far
		case ComponentUShort:
			d := bufUI16Slice(buf)
			for _, v := range d {
				data = append(data, float32(v))
			}
		case ComponentFloat:
			data = bufF32Slice(buf)
		default:
			panic(fmt.Sprint("Dont't know what to do:", ty))
		}

		elems = append(elems, elem{data, attrib.Size})
		format = append(format, attrib)
	}

	for i, t := range prim.Targets {
		for a, ai := range t {
			if (a == "NORMAL" || a == "TAGENT") && i > 3 {
				log.Printf("TARGET_%v%v not supported", a, i)
				continue
			}
			attrib, ok := attrs["TARGET_"+a]
			if !ok {
				log.Printf("target attr: %v not supported yet", a)
				continue
			}
			attrib.Attrib += fmt.Sprint(i)
			attrib.Define += fmt.Sprint(i)

			buf, _, _ := c.doc.AccessorBuffer(ai)
			data := bufF32Slice(buf)
			elems = append(elems, elem{data, attrib.Size})
			format = append(format, attrib)
		}
	}

	nVerts := len(elems[0].data) / elems[0].sz

	var verts []float32
	var indices interface{}

	add := func(verts []float32, idx ...int) []float32 {
		for _, i := range idx {
			// Add interleaved verts for gorge
			for _, e := range elems {
				off := i * e.sz
				end := off + e.sz
				verts = append(verts, e.data[off:end]...)
			}
		}
		return verts
	}

	for i := 0; i < nVerts; i++ {
		verts = add(verts, i)
	}

	if prim.Indices != nil {
		indices = acBufIndices(c.doc.AccessorBuffer(*prim.Indices))
	}

	meshData := gorge.MeshData{
		Format:      format,
		FrontFacing: gorge.FrontFacingCCW,
		Vertices:    verts,
		Indices:     indices,
	}
	// var ref gorge.Resourcer // nolint
	// c.gorge.RunInMain(func() {
	// ref = c.gorge.ResourceRef(&meshData)
	// })
	c.primRef[prim] = &meshData

	return gorge.NewMesh(&meshData)
}

func (c *gltfCreator) getGTexture(tex *Texture) *gorge.Texture {
	im := c.doc.Images[tex.Source]

	ref, ok := c.texRef[im]
	if !ok {
		// We will remove Images?
		texData, err := loadImage(
			resource.FromContext(c.gorge),
			c.doc,
			im,
		)
		if err != nil {
			c.gorge.Error(err)
			return nil
		}
		ref = texData

		// c.gorge.RunInMain(func() {
		// ref = c.gorge.ResourceRef(texData)
		// })
		c.texRef[im] = ref
	}
	if ref == nil {
		panic("Ref is nil")
	}
	gtex := gorge.NewTexture(ref)

	if c.doc.Samplers != nil && c.doc.Samplers[tex.Sampler] != nil {
		s := c.doc.Samplers[tex.Sampler]
		if v := s.MinFilter; v != nil {
			gtex.FilterMode = gorgeSamplerFilter(*v)
		}
		if v := s.WrapS; v != nil {
			gtex.Wrap[0] = gorgeSamplerWrap(*v)
		}
		if v := s.WrapT; v != nil {
			gtex.Wrap[1] = gorgeSamplerWrap(*v)
		}
	}
	return gtex
}

// GSkin skin information.
type GSkin struct {
	Matrices []m32.Mat4
	Joints   []int
}

// GNode represents a gorge container.
type GNode struct {
	*gorge.TransformComponent

	mesh *GMesh
	// mesh *GMesh
	skin *GSkin
	// Not that we need this?
	entities []gorge.Entity
	children []*GNode
}

// GetEntities implements the entity container and returns the underlying
// primitive entities.
func (n *GNode) GetEntities() []gorge.Entity {
	entities := append([]gorge.Entity{}, n.entities...)
	for _, c := range n.children {
		// This will be solved in the gorge.Add
		entities = append(entities, c)
	}
	return entities
}

// GScene entity container with nodes.
type GScene struct {
	*gorge.TransformComponent
	Nodes []*GNode
}

// GetEntities implements the entity container and returns the existing gltf
// scene entities.
func (s *GScene) GetEntities() []gorge.Entity {
	entities := []gorge.Entity{}
	for _, n := range s.Nodes {
		entities = append(entities, n)
	}
	return entities
}

// GMesh is not a gorge mesh but rather gltf gorged mesh with multiple
// renderables
type GMesh struct {
	Name       string
	Weights    []float32
	primitives []*gorge.RenderableComponent
}

// Special case
func acBufIndices(buf []byte, _ int, ty ComponentType) interface{} {
	switch ty {
	case ComponentUByte:
		return append([]byte{}, buf...) // copy
	case ComponentUShort:
		return append([]uint16{}, bufUI16Slice(buf)...)
	case ComponentUInt:
		return append([]uint32{}, bufUI32Slice(buf)...)
	}
	panic(fmt.Sprintf("unknown indice type: %v", ty))
}

func acBuf(buf []byte, _ int, _ ComponentType) []byte {
	return buf
}

func bufUI16Slice(buf []byte) []uint16 {
	bufLen := len(buf) / 2
	return (*(*[^uint32(0)]uint16)(unsafe.Pointer(&buf[0])))[:bufLen:bufLen]
}

func bufUI32Slice(buf []byte) []uint32 {
	bufLen := len(buf) / 4
	return (*(*[^uint32(0)]uint32)(unsafe.Pointer(&buf[0])))[:bufLen:bufLen]
}

func bufF32Slice(buf []byte) []float32 {
	bufLen := len(buf) / 4
	return (*(*[^uint32(0)]float32)(unsafe.Pointer(&buf[0])))[:bufLen:bufLen]
}

func bufVec3Slice(buf []byte) []m32.Vec3 {
	bufLen := len(buf) / (4 * 3) // 4 bytes per float32
	return (*(*[^uint32(0)]m32.Vec3)(unsafe.Pointer(&buf[0])))[:bufLen:bufLen]
}

func bufVec4Slice(buf []byte) []m32.Vec4 {
	bufLen := len(buf) / (4 * 4) // 4 bytes per float32
	return (*(*[^uint32(0)]m32.Vec4)(unsafe.Pointer(&buf[0])))[:bufLen:bufLen]
}

func bufMat4Slice(buf []byte) []m32.Mat4 {
	bufLen := len(buf) / (4 * 16)
	return (*(*[^uint32(0)]m32.Mat4)(unsafe.Pointer(&buf[0])))[:bufLen:bufLen]
}

func gorgeSamplerFilter(f SamplerFilter) gorge.TextureFilter {
	switch f {
	case SamplerNearest:
		return gorge.TextureFilterPoint
	case SamplerLinear:
		return gorge.TextureFilterLinear
		// Missing others but gorge doesn't support that properly
	default:
		return gorge.TextureFilterLinear
	}
}

func gorgeSamplerWrap(w SamplerWrap) gorge.TextureWrap {
	switch w {
	case SamplerClamp:
		return gorge.TextureWrapClamp
	case SamplerMirroredRepeat:
		return gorge.TextureWrapMirror
	case SamplerRepeat:
		return gorge.TextureWrapRepeat
	default:
		log.Printf("warning sampler wrapper didn't match any: %v", w)
		return 0
	}
}
