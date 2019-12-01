package cmd

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBalanceCmd(t *testing.T) {
	cmd := setupTestingCmd(newBalanceCmd())

	flags := cmd.Flags()
	flags.Set("ledger-dir", "testdata/ledger")

	testCmdOutput(t, cmd, "testdata/fixtures/balance.snap")
}

func setupTestingCmd(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().Set("ledger-dir", "testdata/ledger")
	return cmd
}

func testCmdOutput(t *testing.T, cmd *cobra.Command, snapshot string) {
	buf := new(bytes.Buffer)
	cmd.SetOutput(buf)

	_, err := cmd.ExecuteC()
	require.NoError(t, err)

	assert.Equal(t, loadSnapshot(t, snapshot), buf.String())
}

func loadSnapshot(t *testing.T, filename string) string {
	b, err := ioutil.ReadFile(filename)
	require.NoErrorf(t, err, "couldn't read snapshot file %s", filename)

	return string(b)
}
