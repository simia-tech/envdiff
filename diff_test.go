package envdiff_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/simia-tech/envdiff"
)

func TestDiff(t *testing.T) {
	referenceFile := "# Test comment\nTEST_ONE=\"value\"\n"
	processOutputOne := "TEST_ONE=\"another value\"\nTEST_TWO=\"value\""
	processOutputTwo := "TEST_ONE=\"another value\"\nTEST_THREE=\"value\""

	output := bytes.Buffer{}
	require.NoError(t, envdiff.Diff(&output,
		strings.NewReader(referenceFile),
		strings.NewReader(processOutputOne),
		strings.NewReader(processOutputTwo)))

	assert.Equal(t, "  : # Test comment\n12: TEST_ONE=\"value\"\n", output.String())
}
