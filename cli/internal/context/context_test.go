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

package context

import (
	"fmt"
	"os"
	"testing"
)

func TestNewProduction(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := NewProduction()
	if err != nil {
		t.Fatal(err)
	}

	if ctx.Cwd != cwd {
		t.Fatal(fmt.Errorf("working directory does not match, expected %s, found %s", cwd, ctx.Cwd))
	}
	if ctx.Root != cwd {
		t.Fatal(fmt.Errorf("root does not match, expected %s, found %s", cwd, ctx.Root))
	}
}

func TestNewDevelopment(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}

	if ctx.Cwd != cwd {
		t.Fatal(fmt.Errorf("working directory does not match, expected %s, found %s", cwd, ctx.Cwd))
	}
	if ctx.Root != cwd {
		t.Fatal(fmt.Errorf("root does not match, expected %s, found %s", cwd, ctx.Root))
	}
}
