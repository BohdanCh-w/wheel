package fs_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/bohdanch-w/wheel/fs"
	whlogger "github.com/bohdanch-w/wheel/logger"
)

func TestHotSwapFileReader(t *testing.T) {
	filename := t.TempDir() + "/testfile.txt"

	err := os.WriteFile(filename, []byte("content"), 0o600)
	require.NoError(t, err)

	file, err := fs.NewHotSwapFileReader(whlogger.NewNullLogger(), filename)
	require.NoError(t, err)

	defer file.Close()

	{
		content, changed := file.Content()
		require.True(t, changed)
		require.Equal(t, "content", string(content))
	}

	{
		content, changed := file.Content()
		require.False(t, changed)
		require.Equal(t, "content", string(content))
	}

	err = os.WriteFile(filename, []byte("new content"), 0o600)
	require.NoError(t, err)

	time.Sleep(200 * time.Millisecond)

	{
		content, changed := file.Content()
		require.True(t, changed)
		require.Equal(t, "new content", string(content))
	}
}
