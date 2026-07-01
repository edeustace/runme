//go:build !windows

package command

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/runmedev/runme/v3/document"
	"github.com/runmedev/runme/v3/document/identity"
)

func TestNewProgramConfigFromCodeBlock_CellEnv(t *testing.T) {
	t.Parallel()

	idResolver := identity.NewResolver(identity.AllLifecycleIdentity)

	doc := document.New([]byte("```sh {\"name\": \"my-cell\"}\necho -n test\n```"), idResolver)
	node, err := doc.Root()
	require.NoError(t, err)

	blocks := document.CollectCodeBlocks(node)
	require.Len(t, blocks, 1)
	require.NotEmpty(t, blocks[0].ID())

	cfg, err := NewProgramConfigFromCodeBlock(blocks[0])
	require.NoError(t, err)

	assert.Contains(t, cfg.Env, "RUNME_CELL_NAME=my-cell")
	assert.Contains(t, cfg.Env, "RUNME_CELL_ID="+blocks[0].ID())
}

func TestCommand_FromCodeBlock_CellEnvAvailableInProgram(t *testing.T) {
	t.Parallel()

	idResolver := identity.NewResolver(identity.AllLifecycleIdentity)

	doc := document.New([]byte("```sh {\"name\": \"print-cell\"}\necho -n \"$RUNME_CELL_NAME/$RUNME_CELL_ID\"\n```"), idResolver)
	node, err := doc.Root()
	require.NoError(t, err)

	blocks := document.CollectCodeBlocks(node)
	require.Len(t, blocks, 1)

	cfg, err := NewProgramConfigFromCodeBlock(blocks[0])
	require.NoError(t, err)

	testExecuteCommand(
		t,
		cfg,
		bytes.NewReader(nil),
		"print-cell/"+blocks[0].ID(),
		"",
	)
}
