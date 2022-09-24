/*
Copyright 2022 zoomoid.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package util

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

type FileVisitor struct {
	path string
}

func (v *FileVisitor) Visit(fn func(path string) error) error {
	return fn(v.path)
}

func ExpandPaths(cwd string, patterns []string, recursive bool) ([]FileVisitor, []error) {
	var errs []error

	vl := []FileVisitor{}
	for _, s := range patterns {
		matches, err := expandIfGlob(cwd, s)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		for _, m := range matches {
			_, err := os.Stat(m)
			if errors.Is(err, fs.ErrNotExist) {
				errs = append(errs, fmt.Errorf("the path %q does not exist", m))
			}
			if err != nil {
				errs = append(errs, fmt.Errorf("path %q cannot be accessed: %w", m, err))
				continue
			}

			visitors, err := expandPathsToVisitors(m, recursive)
			if err != nil {
				errs = append(errs, err)
			}
			vl = append(vl, visitors...)
		}
	}
	return vl, errs
}

func expandIfGlob(cwd string, pattern string) ([]string, error) {
	if _, err := os.Stat(pattern); errors.Is(err, fs.ErrNotExist) {
		p := filepath.Join(cwd, pattern)
		matches, err := filepath.Glob(p)
		if err == nil && len(matches) == 0 {
			// return nil, fmt.Errorf("the path %q does not exist", pattern)
			return []string{}, nil
		}
		if err == filepath.ErrBadPattern {
			return nil, fmt.Errorf("patterns %q is not valid: %w", pattern, err)
		}
		return matches, err
	}
	return []string{pattern}, nil
}

func expandPathsToVisitors(paths string, recursive bool) ([]FileVisitor, error) {
	var visitors []FileVisitor

	err := filepath.WalkDir(paths, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if path != paths && !recursive {
				// selected a directory that isn't the root and not running in recursive mode
				return filepath.SkipDir
			}
			return nil
		}

		v := FileVisitor{
			path: path,
		}

		visitors = append(visitors, v)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return visitors, nil
}
