package resource

import (
	"fmt"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/x/obj"
)

func init() {
	Register((*gorge.MeshData)(nil), ".obj", meshDataLoader)
	Register((*gorge.Mesh)(nil), ".obj", meshLoader)
}

func meshDataLoader(res *Context, v any, name string, _ ...any) error {
	meshData := v.(*gorge.MeshData)

	rd, err := res.Open(name)
	if err != nil {
		return fmt.Errorf("error opening mesh: %w", err)
	}

	d, err := obj.Decode(rd)
	if err != nil {
		return err
	}
	*meshData = *d

	return nil
}

func meshLoader(res *Context, v any, name string, opts ...any) error {
	mesh := v.(*gorge.Mesh)

	var meshData gorge.MeshData
	if err := meshDataLoader(res, &meshData, name, opts...); err != nil {
		return err
	}

	mesh.Resourcer = &meshData

	return nil
}
