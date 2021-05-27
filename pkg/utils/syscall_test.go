package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestSetHijackStdPipeline(t *testing.T) {
	// init
	stderrFile := "/tmp/test_stderr"
	os.Remove(stderrFile)
	// call, test std error only
	SetHijackStdPipeline(stderrFile, false, true)
	time.Sleep(time.Second) // wait goroutine run
	fmt.Fprintf(os.Stderr, "test stderr")
	// verify
	if !verifyFile(stderrFile, "test stderr") {
		t.Error("stderr hijack failed")
	}
	ResetHjiackStdPipeline()
	fmt.Fprintf(os.Stderr, "repaired\n")
}

func verifyFile(p string, data string) bool {
	b, err := ioutil.ReadFile(p)
	if err != nil {
		return false
	}
	return string(b) == data
}
