package cfg

import (
	"os"
	"path/filepath"
	"runtime"
)

var (
	RootPath  string
	CfgDir    string
	LogDir    string
	LoggerDir string
)

func init() {
	initRootPath()
	CfgDir = filepath.Join(RootPath, "cfg")
	LoggerDir = filepath.Join(CfgDir, "logger")
	LogDir = filepath.Join(RootPath, "log")
}

func JoinRootPath(pathList ...string) string {
	l := len(pathList)
	tmpPathList := make([]string, l+1)
	tmpPathList[0] = RootPath
	for i, path := range pathList {
		tmpPathList[i+1] = path
	}
	return filepath.Join(tmpPathList...)
}

func initRootPath() {
	runtimeRootPath := getRootPathFromeRuneTimeCaller()
	executableRootPath := getRootPathFromExecutable()
	rootPath := executableRootPath
	if !isLoggerDirExists(executableRootPath) && isLoggerDirExists(runtimeRootPath) {
		rootPath = runtimeRootPath
	}
	RootPath = rootPath
}

func getRootPathFromeRuneTimeCaller() string {
	// get current file Path
	// ex: {RootPath}\cfg\path.go
	_, p, _, _ := runtime.Caller(0)
	// find RootPath
	for i := 0; i < 2; i++ {
		p = filepath.Dir(p)
	}
	return p
}

func getRootPathFromExecutable() string {
	fp, err := os.Executable()
	if err != nil {
		return ""
	}
	return filepath.Dir(fp)
}

func isLoggerDirExists(rootPath string) bool {
	fullPath := filepath.Join(rootPath, "cfg", "logger")
	_, err := os.Stat(fullPath)
	return !os.IsNotExist(err)
}
