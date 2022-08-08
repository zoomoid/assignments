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
	// DOC is the absolute path of the document to build without its file
	// extension. This is required for most latex/pdflatex commands
	DOC string
	// DOCEXT is the same as DOC but with the file extension (i.e. ".tex")
	// included
	DOCEXT string
	// DIR is the absolute path to the directory of the file to build
	DIR string
	// TMPDIR is the OS's temporary directory, suitable for storing any
	// intermediate files that may be disposed of
	TMPDIR string
	// WORKSPACE_FOLDER is the folder that the current command runs in
	WORKSPACE_FOLDER string
	// RELATIVE_DIR is the relative path to the directory of the file to build
	RELATIVE_DIR string
	// RELATIVE_DOC is the relative path to the document to build
	RELATIVE_DOC string
	// OUTDIR is the directory to which to build, and in our scenario always the
	// same as the directory of the source file
	OUTDIR string
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
