// +build !windows

package adb

import (
	"fmt"
)

// InstallAdb installs adb when not existing
func InstallAdb() bool {
	fmt.Println("Not supported on this OS.")
	return false
}
