// +build !windows

package ostools

import "fmt"

// RunMeElevated runs the current process as elevated
func RunMeElevated() {
	fmt.Println("Not supported in this OS.")
}
