package resource

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/static"
)

const internalPath = "_gorge"

var reInclude = regexp.MustCompile(`#include\s+<(.*)>\s*`)

func init() {
	Register((*gorge.ShaderData)(nil), ".glsl", shaderDataLoader)
	Register((*gorge.Material)(nil), ".glsl", materialLoader)
}

func shaderDataLoader(res *Context, v any, name string, _ ...any) error {
	shaderData := v.(*gorge.ShaderData)

	source, err := shaderSource(res, name)
	if err != nil {
		return err
	}

	*shaderData = gorge.ShaderData{
		Name: name,
		Src:  source,
	}

	return nil
}

func materialLoader(res *Context, v any, name string, _ ...any) error {
	mat := v.(*gorge.Material)

	var shaderData gorge.ShaderData
	if err := shaderDataLoader(res, &shaderData, name); err != nil {
		return err
	}

	mat.Resourcer = &shaderData

	return nil
}

// shaderSource loads a shader into memory, preprocessing the includes
func shaderSource(res *Context, name string) ([]byte, error) {
	rd, err := res.Open(name)
	if err != nil {
		return nil, err
	}
	source, err := ioutil.ReadAll(rd)
	if err != nil {
		return nil, err
	}
	matches := reInclude.FindAllSubmatch(source, -1)
	for _, match := range matches {
		incStmt := match[0]
		incFile := string(match[1])
		var chunk []byte
		if strings.HasPrefix(incFile, internalPath) {
			d, err := static.Data(incFile[len(internalPath):])
			if err != nil {
				return nil, fmt.Errorf("[resource] %q internal shader not found", incFile)
			}
			chunk = d
		} else {
			fpath := filepath.Dir(name)
			d, err := shaderSource(res, filepath.Join(fpath, incFile))
			if err != nil {
				return nil, err
			}
			chunk = d
		}
		// chunk = append([]byte(fmt.Sprintf("#line 0 %q\n", incFile)), chunk...)
		// Add back original line ?? so we can debug better
		source = bytes.ReplaceAll(source, incStmt, chunk)
	}

	return source, nil
}
