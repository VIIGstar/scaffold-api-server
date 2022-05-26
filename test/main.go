package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"scaffold-api-server/internal/services/log"
	"scaffold-api-server/pkg/config"
	"strings"
	"sync"
)

const (
	NoTestFileResult = "[no test files]"
)

var (
	_wg = sync.WaitGroup{}
)

func LoopDirsFiles(path string) {
	lg := log.NewLogger(config.LogConfig{})
	_wg.Add(1)
	defer _wg.Done()
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}

	for _, file := range files {
		if file.IsDir() && strings.Index(file.Name(), ".") != 0 {
			//Recursively go further in the tree
			go LoopDirsFiles(filepath.Join(path, file.Name()))
		}
	}

	if len(files) > 0 {
		cmd := exec.Command("bash", "-c", fmt.Sprintf("go test -v %v", path))
		rs, err := cmd.Output()
		if err == nil && !strings.Contains(string(rs), NoTestFileResult) {
			rs, _ = exec.Command("bash", "-c", fmt.Sprintf("go test %v -coverprofile cover.out", path)).Output()
			lg.Info(string(rs))
		}
	}
}

func main() {
	src, _ := os.Getwd()
	LoopDirsFiles(src)
	_wg.Wait()
}
