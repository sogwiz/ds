package master

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestPutFile(t *testing.T) {
	fileContent, _ := ioutil.ReadFile("/Users/agilbert/go/src/ds/main.go")
	PutFile(UserID(1), "myfile.txt", bytes.NewReader(fileContent))
	t.Fail()
}
