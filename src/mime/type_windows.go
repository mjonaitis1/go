// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mime

import (
	"internal/syscall/windows/registry"
)

func init() {
	osInitMime = initMimeWindows
}

func initMimeWindows() {
	names := make([]string, 0)
	err := registry.CLASSES_ROOT.ReadSubKeyNames(func(s string) error {
		names = append(names, s)
		return nil
	})
	if err != nil {
		return
	}

	for _, name := range names {
		if len(name) < 2 || name[0] != '.' { // looking for extensions only
			continue
		}
		k, err := registry.OpenKey(registry.CLASSES_ROOT, name, registry.READ)
		if err != nil {
			continue
		}
		v, _, err := k.GetStringValue("Content Type")
		k.Close()
		if err != nil {
			continue
		}
		setExtensionType(name, v)
	}
}

func initMimeForTests() map[string]string {
	return map[string]string{
		".PnG": "image/png",
	}
}
