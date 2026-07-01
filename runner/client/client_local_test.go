//go:build !windows

package client

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/runmedev/runme/v3/project"
	"github.com/runmedev/runme/v3/runner"
)

func init() {
	// The default env dump command re-executes the current binary, which in
	// tests is the test binary itself and would recurse indefinitely.
	runner.SetEnvDumpCommandForTesting()
}

func TestLocalRunner_RunTask_CellEnv(t *testing.T) {
	dir := t.TempDir()
	source := "```sh {\"name\": \"print-cell\", \"interactive\": false}\necho -n \"$RUNME_CELL_NAME/$RUNME_CELL_ID\"\n```\n"
	err := os.WriteFile(filepath.Join(dir, "README.md"), []byte(source), 0o600)
	require.NoError(t, err)

	proj, err := project.NewDirProject(dir)
	require.NoError(t, err)

	tasks, err := project.LoadTasks(context.Background(), proj)
	require.NoError(t, err)
	require.Len(t, tasks, 1)

	task := tasks[0]

	stdout := new(bytes.Buffer)

	localRunner, err := NewLocalRunner(
		WithInsecure(true),
		WithProject(proj),
		WithStdout(stdout),
		WithStderr(io.Discard),
	)
	require.NoError(t, err)

	require.NoError(t, localRunner.RunTask(context.Background(), task))

	expected := task.CodeBlock.Name() + "/" + task.CodeBlock.ID()
	assert.Equal(t, "print-cell", task.CodeBlock.Name())
	assert.Equal(t, expected, stdout.String())
}
