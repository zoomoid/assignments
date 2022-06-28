package context

import (
	"fmt"
	"os"
	"testing"
)

func TestNewLogger(t *testing.T) {
	t.Run("development", func(t *testing.T) {
		_, err := newLogger(false)
		if err != nil {
			t.Fatal(err)
		}
	})
	t.Run("production", func(t *testing.T) {
		_, err := newLogger(true)
		if err != nil {
			t.Fatal(err)
		}
	})
}

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
