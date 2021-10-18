package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	seccomp "github.com/seccomp/libseccomp-golang"
)

func main() {
	if err := xmain(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func xmain() error {
	if len(os.Args) < 2 {
		fmt.Println("Workaround program for running glibc >= 2.34 distros on Docker <= 20.10.9")
		fmt.Printf("Usage: %s COMMAND [ARGS...]\n", os.Args[0])
		major, minor, micro := seccomp.GetLibraryVersion()
		fmt.Printf("seccomp version: %d.%d.%d\n", major, minor, micro)
		return nil
	}
	argv := os.Args[1:]
	arg0, err := exec.LookPath(argv[0])
	if err != nil {
		return err
	}
	filter, err := seccomp.NewFilter(seccomp.ActAllow)
	if err != nil {
		return fmt.Errorf("failed to create a seccomp filter: %w", err)
	}
	clone3, err := seccomp.GetSyscallFromName("clone3")
	if err != nil {
		return fmt.Errorf("failed to get syscall \"clone3\": %w", err)
	}
	act := seccomp.ActErrno
	act = act.SetReturnCode(int16(syscall.ENOSYS))
	if err := filter.AddRule(clone3, act); err != nil {
		return fmt.Errorf("failed to add action %v for \"clone3\": %w", act, err)
	}
	if err := filter.Load(); err != nil {
		return fmt.Errorf("failed to load the seccomp filter: %w", err)
	}
	return syscall.Exec(arg0, argv, os.Environ())
}
