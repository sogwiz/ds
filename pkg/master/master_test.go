package master

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPutFile(t *testing.T) {
	fileContent, _ := ioutil.ReadFile("/Users/agilbert/go/src/ds/main.go")
	PutFile(UserID(1), "myfile.txt", bytes.NewReader(fileContent))
	assert.True(t, false)
}
