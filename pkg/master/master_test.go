package master

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPutFile(t *testing.T) {
	homeDir, _ := os.UserHomeDir()
	fileContent, _ := ioutil.ReadFile(filepath.Join(homeDir, "go", "src", "ds", "main.go"))
	PutFile("user_1/myfile.txt", bytes.NewReader(fileContent))
	assert.True(t, false)
}
