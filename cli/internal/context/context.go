package context

import (
	config "github.com/zoomoid/assignments/v1/internal/config"
	zap "go.uber.org/zap"
)

// AppContext is initialized by a cobra command and carried through the application
// to provide a shared space of root-level configuration
type AppContext struct {
	// Logger is the shared application zap logger
	Logger *zap.SugaredLogger
	// Cwd is the working directory from which the command was invoked
	Cwd string
	// Root is the repository's root, either the Cwd, or any directory *above*
	// Cwd that contains the .assignments.yaml file, marking it the repository's root
	// This way, we can e.g. call `assignmentctl build -f .` in $ROOT/assignment-06
	// to build assignment-06 without having to return to root
	// This behaves similar to e.g. git when commands are used in subdirectories
	Root string
	// Contains the shared (read-in) Configuration struct used for e.g. build, generate,
	// bundle and more. NOTE that is might not be present all the time, as it is only read
	// from file when ctx.Read() is called. ctx.Read() should be followed by `defer ctx.Write()`
	// in case any mutations to the struct happen to persist those back into the file
	Configuration *config.Configuration
}

// Read uses the context's root to read a configmap into the context's struct field
func (c *AppContext) Read() error {
	cfg, err := config.ReadConfigMap(c.Root)
	if err != nil {
		return err
	}
	c.Configuration = cfg
	return nil
}

// Write writes the context's struct field to a file at the context's root
func (c *AppContext) Write() error {
	err := config.WriteConfigMap(c.Configuration, c.Root)
	return err
}

// Clone copies all fields except the logger into a fresh context and returns a reference to it
func (c *AppContext) Clone() *AppContext {
	nc := &AppContext{
		Cwd:           c.Cwd,
		Root:          c.Root,
		Configuration: c.Configuration.Clone(),
		Logger:        c.Logger, // don't actually clone the zap logger, this is fine to be aliased by all contexts
	}
	return nc
}
