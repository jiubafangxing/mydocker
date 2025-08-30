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
	
	// 获取当前可执行文件的绝对路径
	exePath, err := os.Executable()
	if err != nil {
		logrus.Errorf("Failed to get executable path: %v", err)
		// 回退到使用 /proc/self/exe
		exePath = "/proc/self/exe"
	}
	
	cmd := exec.Command(exePath, args...)
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
	logrus.Infof("init process started, command: %s", command)

	// 设置默认的mount flags
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV

	// 先卸载旧的proc文件系统（如果存在）
	syscall.Unmount("/proc", syscall.MNT_DETACH)

	// 重新挂载proc文件系统
	logrus.Info("Mounting /proc filesystem...")
	if err := syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), ""); err != nil {
		logrus.Errorf("Mount proc error: %v", err)
		return err
	}
	logrus.Info("Successfully mounted /proc")
	
	// 查找命令路径
	path, err := exec.LookPath(command)
	if err != nil {
		logrus.Errorf("Command not found: %s, error: %v", command, err)
		return err
	}
	logrus.Infof("Found command path: %s", path)
	
	// 执行用户命令，替换当前进程
	logrus.Infof("Executing command: %s with args: %v", path, args)
	if err := syscall.Exec(path, args[0:], os.Environ()); err != nil {
		logrus.Errorf("Failed to exec command: %v", err)
		return err
	}
	return nil
}
