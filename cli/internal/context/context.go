package context

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/zoomoid/assignments/v1/internal/config"
	zap "go.uber.org/zap"
)

// AppContext is initialized by a cobra command and carried through the application
// to provide a shared space of root-level configuration
type AppContext struct {
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
	// Verbose toggles more explicit output down the line
	Verbose         bool
	rootInitialized bool
}

// Read uses the context's root to read a configmap into the context's struct field
func (c *AppContext) Read() error {
	// TODO: we don't strictly need this anymore because only those commands that *need* config also *load* it
	if !c.rootInitialized {
		if err := c.mustFindConfigFile(); err != nil {
			return err
		}
		c.rootInitialized = true
	}
	p := filepath.Join(c.Root, ".assignments.yaml")
	cfg, err := config.Read(p)
	if err != nil {
		return err
	}
	c.Configuration = cfg
	return nil
}

// Write writes the context's struct field to a file at the context's root
func (c *AppContext) Write() error {
	p := filepath.Join(c.Root, ".assignments.yaml")
	err := config.Write(c.Configuration, p)
	return err
}

func (c *AppContext) mustFindConfigFile() error {
	// if we cannot find a configuration file in here, traverse the file tree upwards
	// until either the root or we find a config file
	cfgPath, err := config.Find(c.Cwd)
	if err != nil {
		return errors.New("failed to find configmap in working directory or above. Is the directory initialized?")
	}
	c.Root = cfgPath
	return nil
}

// Clone copies all fields except the logger into a fresh context and returns a reference to it
func (c *AppContext) Clone() *AppContext {
	nc := &AppContext{
		Cwd:           c.Cwd,
		Root:          c.Root,
		Configuration: c.Configuration.Clone(),
	}
	return nc
}

// NewDevelopment creates a new AppContext with development logger from scratch
func NewDevelopment() (context *AppContext, err error) {
	if err != nil {
		return nil, err
	}
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return &AppContext{
		Configuration: nil,
		Cwd:           cwd,
		Root:          cwd,
	}, nil
}

// NewProduction creates a new AppContext with production logger from scratch
func NewProduction() (context *AppContext, err error) {
	if err != nil {
		return nil, err
	}
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return &AppContext{
		Configuration: nil,
		Cwd:           cwd,
		Root:          cwd,
	}, nil
}

// newLogger creates either a new production or a new development zap logger
func newLogger(production bool) (*zap.SugaredLogger, error) {
	var l *zap.Logger
	var err error
	if production {
		l, err = zap.NewProduction()
	} else {
		l, err = zap.NewDevelopment()
	}
	if err != nil {
		return nil, err
	}
	defer l.Sync()
	return l.Sugar(), nil
}
