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
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// EmbedLoader in memory asset loader
type EmbedLoader struct {
	Data map[string][]byte
}

// Load loads the asset
func (l EmbedLoader) Load(p string) (io.ReadCloser, error) {
	if l.Data == nil {
		return nil, fmt.Errorf("%s not found: %w", p, os.ErrNotExist)
	}
	d, ok := l.Data[p]
	if !ok {
		return nil, fmt.Errorf("%s not found: %w", p, os.ErrNotExist)
	}
	rd := bytes.NewReader(d)
	return ioutil.NopCloser(rd), nil
}
