package cmdutil_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
)

func TestSaveUsesOwnerOnlyPermissions(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Cleanup(homedir.Reset)
	homedir.Reset()
	t.Cleanup(func() { viper.SetConfigFile("") })
	viper.SetConfigFile("")

	config := cmdutil.NewFactory(cmdutil.Version{}).Config()
	config.SetString(cmdutil.CONF_TOKEN, "a token")
	require.NoError(t, config.Save())

	dir := filepath.Join(home, ".config", "clockify-cli")
	filename := filepath.Join(dir, ".clockify-cli.yaml")

	dirInfo, err := os.Stat(dir)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o700), dirInfo.Mode().Perm())

	fileInfo, err := os.Stat(filename)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o600), fileInfo.Mode().Perm())

	// a file left world-readable by an older version gets tightened on save
	require.NoError(t, os.Chmod(filename, 0o644))
	require.NoError(t, config.Save())

	fileInfo, err = os.Stat(filename)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o600), fileInfo.Mode().Perm())
}
