package kerlong

import (
	"os"
	"path/filepath"
)

//GetExePath() 获取当前可执行文件所在目录
func GetExePath() (string, error) {
	ePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(ePath), nil
}
