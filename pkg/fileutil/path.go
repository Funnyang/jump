package fileutil

import "os"

// ExistPath 返回路径是否存在
func ExistPath(p string) bool {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return false
	}
	return true
}
