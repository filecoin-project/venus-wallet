package filemgr

import (
	"io/ioutil"
	"os"
	"testing"

	"gotest.tools/assert"
)

func TestNewFS(t *testing.T) {
	fsPath, err := ioutil.TempDir("", "venus-repo-")
	defer os.RemoveAll(fsPath)
	if err != nil {
		t.Fatal(err)
	}
	targetAPI := "/ip4/0.0.0.0/tcp/1334/http"
	fs, err := NewFS(fsPath, &OverrideParams{
		API: targetAPI,
	})
	if err != nil {
		t.Fatal(err)
	}
	curAPI, err := fs.APIEndpoint()
	if err != nil {
		t.Fatal()
	}
	assert.Equal(t, curAPI, targetAPI)
}
