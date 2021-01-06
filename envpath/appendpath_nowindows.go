// +build !windows

package envpath

import "fmt"

// AppendToEnvPath appends the parameter to the PATH environment
func AppendToEnvPath(newPath string) {
	fmt.Println("Not supported in the current OS.")
}
