package util

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"syscall"
)

func CheckPlatform() {
	// mac
	if runtime.GOOS == "darwin" {
		if !IsRunningAsRoot() {
			fmt.Println("Trying to Relauch with Sudo")
			RelaunchAsRoot()
			return
		}
	}

}

func IsRunningAsRoot() bool {
	return syscall.Geteuid() == 0
}

func RelaunchAsRoot() {
	executable, err := os.Executable()
	if err != nil {
		fmt.Println("Error getting executable:", err)
		return
	}

	cmd := exec.Command("sudo", executable)

	// same terminal
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := syscall.Exec("/usr/bin/sudo", []string{"sudo", executable}, os.Environ()); err != nil {
		fmt.Println("Error executing command:", err)
	}
}
