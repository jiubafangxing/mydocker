//go:build linux

package container

import (
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"syscall"

	"golang.org/x/sys/unix"
)

// NewParentProcess 创建一个新的父进程 (Linux版本)
func NewParentProcess(tty bool, command string) *exec.Cmd {
	args := []string{"init", command}
	cmd := exec.Command("/proc/self/exe", args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: unix.CLONE_NEWUTS | unix.CLONE_NEWPID |
			unix.CLONE_NEWNS | unix.CLONE_NEWNET |
			unix.CLONE_NEWIPC,
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
	logrus.Infof("command %s", command)

	// 设置默认的mount flags
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV

	// 重新挂载proc文件系统
	if err := syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), ""); err != nil {
		logrus.Errorf("Mount proc error %v", err)
		return err
	}
	
	// 查找命令路径
	path, err := exec.LookPath(command)
	if err != nil {
		logrus.Errorf("Exec loop path error %v", err)
		return err
	}
	logrus.Infof("Find path %s", path)
	
	// 执行用户命令
	if err := syscall.Exec(path, args[0:], os.Environ()); err != nil {
		logrus.Error(err.Error())
		return err
	}
	return nil
}
