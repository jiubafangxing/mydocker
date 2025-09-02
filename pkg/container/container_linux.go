//go:build linux

package container

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// NewParentProcess 创建一个新的父进程 (Linux版本)
func NewParentProcess(tty bool, command string) (*exec.Cmd, *os.File) {
	args := []string{"init", command}

	// 获取当前可执行文件的绝对路径
	exePath, err := os.Executable()
	if err != nil {
		logrus.Errorf("Failed to get executable path: %v", err)
		// 回退到使用 /proc/self/exe
		exePath = "/proc/self/exe"
	}
	logrus.Infof("start run proc/self.exe")
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

	//设置匿名管道用户主子进程进行通信

	r, w, err := os.Pipe()
	if nil != err {
		logrus.Errorf("NewParentProcess err to create pipe %s", err)
		return nil, nil
	}
	cmd.ExtraFiles = []*os.File{r}
	return cmd, w
}

// RunContainerInitProcess 运行容器初始化进程 (Linux版本)
func RunContainerInitProcess() error {
	command := readUserCommand()
	if nil == command || len(command) == 0 {
		return fmt.Errorf("not fetch command from master process")
	}
	logrus.Infof("init process started, command: %s", command)
	// 设置默认的mount flags
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")
	if err := syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), ""); err != nil {
		logrus.Errorf("Mount proc error: %v", err)
		return err
	}
	cmdAbsloutePath, err := exec.LookPath(command[0])
	if nil != err {
		return fmt.Errorf("command not exist")
	}

	if err := syscall.Exec(cmdAbsloutePath, command[0:], os.Environ()); err != nil {
		logrus.Errorf("Failed to exec command: %v", err)
		return err
	}
	return nil
}

func readUserCommand() []string {
	r := os.NewFile(uintptr(3), "pipe")
	contents, err := io.ReadAll(r)
	if nil != err {
		return nil
	}
	cmdStr := string(contents)
	return strings.Split(cmdStr, " ")
}
