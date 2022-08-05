package runner

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/rs/zerolog/log"
	"github.com/zoomoid/assignments/v1/internal/config"
)

type substitutionContext struct {
	DOC              string
	DOCEXT           string
	DIR              string
	TMPDIR           string
	WORKSPACE_FOLDER string
	RELATIVE_DIR     string
	RELATIVE_DOC     string
	OUTDIR           string
}

func commandsFromRecipe(recipe *config.Recipe, cwd string, file string, quiet bool) ([]*exec.Cmd, error) {
	cmds := []*exec.Cmd{}

	ctx := makeSubstitutionContext(cwd, file)

	for i, tool := range *recipe {
		if tool.Command == "" {
			return nil, fmt.Errorf("failed to make build commands, missing program in recipe step %d", i)
		}
		program := tool.Command
		args := []string{}
		if len(tool.Args) > 0 {
			args = tool.Args
		}

		for i, arg := range args {
			args[i] = findAndSubstituteReservedSymbols(arg, ctx)
		}

		// args = append(args)
		cmd := exec.Command(program, args...)

		out := &bytes.Buffer{}

		if quiet {
			sink := bufio.NewWriter(out)
			cmd.Stdout = sink
		} else {
			cmd.Stdout = os.Stdout
		}
		cmd.Dir = cwd

		cmds = append(cmds, cmd)
	}

	return cmds, nil
}

func makeSubstitutionContext(cwd string, file string) *substitutionContext {

	dir := filepath.Dir(file)
	base := filepath.Base(file)
	ext := filepath.Ext(file)
	name := strings.ReplaceAll(base, ext, "")

	doc := filepath.Join(filepath.Clean(dir), name)

	relativeDir, err := filepath.Rel(cwd, dir)
	if err != nil {
		relativeDir = dir
	}

	relativeDoc, err := filepath.Rel(cwd, doc)
	if err != nil {
		relativeDoc = doc
	}

	return &substitutionContext{
		WORKSPACE_FOLDER: cwd,
		DOCEXT:           base,
		DIR:              filepath.Clean(dir),
		DOC:              doc,
		TMPDIR:           os.TempDir(),
		RELATIVE_DIR:     relativeDir,
		RELATIVE_DOC:     relativeDoc,
		OUTDIR:           filepath.Clean(dir),
	}
}

func findAndSubstituteReservedSymbols(arg string, context *substitutionContext) string {
	tpl, err := template.New("substitution").Parse(arg)
	if err != nil {
		log.Warn().Err(err)
		return arg
	}
	out := &bytes.Buffer{}
	err = tpl.Execute(out, *context)
	if err != nil {
		log.Warn().Err(err)
		return arg
	}
	return out.String()
}
