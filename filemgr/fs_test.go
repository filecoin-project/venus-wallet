// stm: #unit
package filemgr

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
)

func TestNewFS(t *testing.T) {
	// stm: @VENUSWALLET_FILEMGR_FS_NEW_001
	fsPath, err := ioutil.TempDir("", "venus-repo-")
	defer os.RemoveAll(fsPath)
	if err != nil {
		t.Fatal(err)
	}
	targetAPI := "/ip4/0.0.0.0/tcp/1334/http"
	fs, err := NewFS(fsPath, &OverrideParams{
		API: targetAPI,
	})
	require.NoError(t, err)
	require.NotNil(t, fs)

	// stm: @VENUSWALLET_FILEMGR_FS_API_SECRET_001,
	secr, err := fs.APISecret()
	require.NoError(t, err)
	require.NotNil(t, secr)

	// stm: @VENUSWALLET_FILEMGR_FS_API_STRATEGY_TOKEN_001
	token, err := fs.APIStrategyToken("password")
	require.NoError(t, err)
	require.NotEqual(t, token, "")

	curAPI, err := fs.APIEndpoint()
	if err != nil {
		t.Fatal()
	}
	assert.Equal(t, curAPI, targetAPI)
}
