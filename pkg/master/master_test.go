package master

import (
	"bytes"
	"ds/pkg/master/metadata"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPutFile(t *testing.T) {
	homeDir, _ := os.UserHomeDir()
	fileContent, _ := ioutil.ReadFile(filepath.Join(homeDir, "go", "src", "ds", "main.go"))
	PutFile(metadata.UserID(1), "myfile.txt", bytes.NewReader(fileContent))
	assert.True(t, false)
}
