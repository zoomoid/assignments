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

type empty struct{}

type Set map[string]empty

func NewSet(items ...string) Set {
	ss := Set{}
	ss.Insert(items...)
	return ss
}

func (s Set) Insert(items ...string) Set {
	for _, item := range items {
		s[item] = empty{}
	}
	return s
}

func (s Set) Delete(items ...string) Set {
	for _, item := range items {
		delete(s, item)
	}
	return s
}

func (s Set) Has(item string) bool {
	_, contained := s[item]
	return contained
}

func (s Set) Len() int {
	return len(s)
}
