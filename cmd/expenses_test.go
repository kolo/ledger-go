package cmd

import (
	"testing"
)

func TestExpensesCmd(t *testing.T) {
	cmd := setupTestingCmd(newExpensesCmd())

	testCmdOutput(t, cmd, "testdata/fixtures/expenses.snap")
}

func TestWeeklyExpensesCmd(t *testing.T) {
	cmd := setupTestingCmd(newExpensesCmd())
	cmd.Flags().Set("weekly", "true")

	testCmdOutput(t, cmd, "testdata/fixtures/weekly_expenses.snap")
}
