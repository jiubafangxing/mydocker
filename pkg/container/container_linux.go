//go:build linux

package container

import (
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"syscall"
)

// NewParentProcess 创建一个新的父进程 (Linux版本)
func NewParentProcess(tty bool, command string) *exec.Cmd {
	args := []string{"init", command}

	// 获取当前可执行文件的绝对路径
	exePath, err := os.Executable()
	if err != nil {
		logrus.Errorf("Failed to get executable path: %v", err)
		// 回退到使用 /proc/self/exe
		exePath = "/proc/self/exe"
	}
	log.Infof("start run proc/self.exe")
	cmd := exec.Command(exePath, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS | syscall.CLONE_NEWNET |
			syscall.CLONE_NEWIPC,
	}
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	return cmd
}

// RunContainerInitProcess 运行容器初始化进程 (Linux版本)
func RunContainerInitProcess(command string, args []string) error {
	logrus.Infof("init process started, command: %s", command)

	// 设置默认的mount flags
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	if err := syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), ""); err != nil {
		logrus.Errorf("Mount proc error: %v", err)
		return err
	}
	logrus.Info("Successfully mounted /proc")
	argv := []string{command}
	// 执行用户命令，替换当前进程
	//logrus.Infof("Executing command: %s with args: %v", path, args)
	if err := syscall.Exec(command, argv, os.Environ()); err != nil {
		logrus.Errorf("Failed to exec command: %v", err)
		return err
	}
	return nil
}
