package cmd

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestParseArgs(t *testing.T) {
    cmd, args, flags, err := parseArgs([]string{"find", "demo", "-P"}, map[string]int{"-P": 0, "--plain": 0})
    assert.Nil(t, err)
    assert.Equal(t, "find", cmd)
    assert.Equal(t, []string{"demo"}, args)
    assert.Equal(t, map[string][]string{"-P": {}}, flags)

    cmd, args, flags, err = parseArgs([]string{"find", "demo", "-P", "123", "456"}, map[string]int{"-P": 1, "--plain": 0})
    assert.Nil(t, err)
    assert.Equal(t, "find", cmd)
    assert.Equal(t, []string{"demo", "456"}, args)
    assert.Equal(t, map[string][]string{"-P": {"123"}}, flags)

    cmd, args, flags, err = parseArgs([]string{"find", "demo", "--plain", "123", "456"}, map[string]int{"-P": 1, "--plain": 2})
    assert.Nil(t, err)
    assert.Equal(t, "find", cmd)
    assert.Equal(t, []string{"demo"}, args)
    assert.Equal(t, map[string][]string{"--plain": {"123", "456"}}, flags)
}
