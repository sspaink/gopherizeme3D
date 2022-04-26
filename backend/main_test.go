package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateProgram(t *testing.T) {
	err := GenerateProgram()
	require.NoError(t, err)
}
