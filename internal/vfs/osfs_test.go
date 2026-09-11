// Copyright (C) 2019 Nicola Murino
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, version 3.

package vfs

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolvePathAllowsRootEscapeOnlyInLegacyMode(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink behavior differs on Windows")
	}

	baseDir := t.TempDir()
	homeDir := filepath.Join(baseDir, "home")
	outsideDir := filepath.Join(baseDir, "outside")
	require.NoError(t, os.Mkdir(homeDir, 0o755))
	require.NoError(t, os.Mkdir(outsideDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(outsideDir, "file.txt"), []byte("outside"), 0o644))
	require.NoError(t, os.Symlink(outsideDir, filepath.Join(homeDir, "external")))

	fs := NewOsFs("", homeDir, "", nil)
	_, err := fs.ResolvePath("/external/file.txt")
	require.Error(t, err)

	SetAllowRootEscape(true)
	defer SetAllowRootEscape(false)
	resolved, err := fs.ResolvePath("/external/file.txt")
	require.NoError(t, err)

	file, _, _, err := fs.Open(resolved, 0)
	require.NoError(t, err)
	defer file.Close()
	content, err := io.ReadAll(file)
	require.NoError(t, err)
	assert.Equal(t, "outside", string(content))
}
