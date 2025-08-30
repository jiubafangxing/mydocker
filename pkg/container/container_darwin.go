//go:build darwin

package container

import (
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

// NewParentProcess 创建一个新的父进程 (macOS版本)
// 注意：macOS不支持Linux容器的namespace功能，这里提供一个基础实现
func NewParentProcess(tty bool, command string) *exec.Cmd {
	args := []string{"init", command}
	log.Infof("start run proc/self.exe")
	cmd := exec.Command("/proc/self/exe", args...)
	// macOS不支持Linux的namespace，所以不设置特殊的SysProcAttr
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	return cmd
}

// RunContainerInitProcess 运行容器初始化进程 (macOS版本)
func RunContainerInitProcess(command string, args []string) error {
	log.Infof("command %s", command)

	path, err := exec.LookPath(command)
	if err != nil {
		log.Errorf("Exec loop path error %v", err)
		return err
	}
	log.Infof("Find path %s", path)

	// 在macOS上使用常规的exec执行
	cmd := exec.Command(path, args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
