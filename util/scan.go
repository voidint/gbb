package util

import (
	"bufio"
	"os"
	"strings"
)

// Scanln 从标准输入中读取一行
func Scanln() (line string, err error) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line = scanner.Text()
		break
	}
	return strings.TrimSpace(line), scanner.Err()
}
