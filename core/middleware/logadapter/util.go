package logadapter

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

func sourceDir() {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	dir = filepath.Dir(dir)
	dir = filepath.Dir(dir)

	s := filepath.Dir(dir)
	if filepath.Base(s) != "logadapter" && filepath.Base(s) != "helper" {
		s = dir
	}
	baseSourceDir = filepath.ToSlash(s) + "/"
}

func getCaller() string {
	for i := 2; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		if ok && (!strings.HasPrefix(file, baseSourceDir) || strings.HasSuffix(file, "_test.go")) {
			return fmt.Sprintf("%s:%d", file, line)
		}
	}

	return ""
}
