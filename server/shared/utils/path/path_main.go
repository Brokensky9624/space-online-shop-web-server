package path

import (
	"os"
	"path/filepath"
	"runtime"
)

var (
	RootPath  string
	CfgDir    string
	LoggerDir string
)

func init() {
	initRootPath()
	CfgDir = JoinRootPath("cfg")
	LoggerDir = filepath.Join(CfgDir, "logger")
}

func initRootPath() {
	runtimePath := getRuntimePath()
	executePath := getExecutePath()

	rootPath := executePath
	if !isLoggerDirExist(executePath) && isLoggerDirExist(runtimePath) {
		rootPath = runtimePath
	}

	RootPath = rootPath
}

func getRuntimePath() string {
	_, fp, _, _ := runtime.Caller(0)
	for i := 0; i < 3; i++ {
		fp = filepath.Dir(fp)
	}
	return fp
}

func getExecutePath() string {
	ex, err := os.Executable()
	if err != nil {
		return ""
	}
	return filepath.Dir(ex)
}

func isLoggerDirExist(path string) bool {
	loggerPath := filepath.Join(path, "cfg", "logger")
	_, err := os.Stat(loggerPath)
	return !os.IsNotExist(err)
}

func JoinRootPath(pathList ...string) string {
	pathLen := len(pathList)
	if pathLen == 0 {
		return RootPath
	}
	p := make([]string, pathLen+1)
	p[0] = RootPath
	copy(p[1:], pathList)
	return filepath.Join(p...)
}
