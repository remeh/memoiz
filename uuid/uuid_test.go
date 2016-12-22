package uuid

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	require := require.New(t)

	str := "f022d162-2aa1-4ff9-a49b-b5ad2b55f860"
	uid, err := Parse(str)
	require.NoError(err, "parsing %s shouldn't fail", str)
	require.Equal(uid.String(), str, "read uuid should be equal to the str uid")

	str = "f162-2aa1-4ff9-a49b-b5ad2b55f860"
	uid, err = Parse(str)
	require.Error(err, "parsing %s should fail", str)
}

func TestScan(t *testing.T) {
	require := require.New(t)

	zero, err := Parse("00000000-0000-0000-0000-000000000000")
	require.NoError(err, "parsing this value shouldn't fail.")
	require.Equal(zero.String(), "00000000-0000-0000-0000-000000000000")

	str := "f0aaa162-2aa1-4ff9-a49b-b5ad2b55f860"
	err = zero.Scan(str)
	require.NoError(err, "scanning %s in the UUID shouldn't fail", str)
	require.Equal(zero.String(), str, "the scanned UUID should be equal to %s", str)

	err = zero.Scan("cccc")
	require.Error(err, "scanning %s in the UUID should fail", str)
}
