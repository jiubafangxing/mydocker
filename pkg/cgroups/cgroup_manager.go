package cgroups

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type ResourceConfig struct {
	MemoryLimit string
	CpuShare    string
	CpuSet      string
}

type CgroupManager interface {
	// 创建或获取 cgroup
	Create() error

	// 设置 cgroup 资源限制
	Set(res *ResourceConfig) error

	// 将进程添加到 cgroup
	AddProcess(pid int) error

	// 获取 cgroup 中的进程列表
	GetProcesses() ([]int, error)

	// 删除 cgroup
	Destroy() error
}

type GroupV2Manager struct {
	Path string // cgroup 路径
	Dir  string // cgroup 文件系统挂载点
}

func NewV2CgroupManager(path string) (*GroupV2Manager, error) {
	rootPath, err := findCgroup2Mountpoint()
	if nil != err {
		return nil, err
	}
	result := &GroupV2Manager{
		Path: rootPath,
		Dir:  fmt.Sprintf("%s/%s", rootPath, path),
	}
	return result, nil
}

func (self *GroupV2Manager) Create() error {
	if err := os.MkdirAll(self.Path, 0755); nil != err {
		return err
	}
	return nil
}

func (self *GroupV2Manager) Set(res *ResourceConfig) error {
	if nil == res {
		return fmt.Errorf("input param valid")
	}
	if res.MemoryLimit != "" {
		if err := writeFile(self.Dir, "memory.max", res.MemoryLimit); err != nil {
			return err
		}
	}
	if res.CpuSet != "" {
		if err := writeFile(self.Dir, "cpuset.cpus", res.CpuSet); err != nil {
			return err
		}
	}
	if res.CpuShare != "" {
		if err := writeFile(self.Dir, "cpu.weight", res.CpuShare); err != nil {
			return fmt.Errorf("set cpu.weight: %w", err)
		}
	}
	return nil
}

func (self *GroupV2Manager) AddProcess(pid int) error {
	return os.WriteFile(filepath.Join(self.Dir, "cgroup.procs"), []byte(strconv.Itoa(pid)), 0644)
}

func (self *GroupV2Manager) GetProcesses() ([]int, error) {
	filePath := filepath.Join(self.Dir, "cgroup.procs")

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %w", filePath, err)
	}
	defer file.Close()

	var pids []int
	scanner := bufio.NewScanner(file)

	// 逐行读取
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue // 跳过空行
		}

		// 将字符串转换为整数
		pid, err := strconv.Atoi(line)
		if err != nil {
			// 可以选择记录错误并继续，或者直接返回错误
			return nil, fmt.Errorf("invalid PID format '%s' in %s: %w", line, filePath, err)
		}

		pids = append(pids, pid)
	}

	// 检查扫描过程中是否有错误
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading %s: %w", filePath, err)
	}

	return pids, nil
}

func (self *GroupV2Manager) Destroy() error {
	err := os.Remove(self.Dir)
	return err
}

func findCgroup2Mountpoint() (string, error) {
	handler, err := os.Open("/proc/mounts")
	if nil != err {
		return "", err
	}
	defer handler.Close()
	scanner := bufio.NewScanner(handler)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) >= 3 && fields[2] == "cgroup2" {
			return fields[1], nil
		}
	}
	return "", fmt.Errorf("cgroup2 mountpoint not found")
}

func writeFile(dir string, filename string, info string) error {
	file, err := os.Create(fmt.Sprintf("%s/%s", dir, filename))
	if nil != err {
		return err
	}
	_, err = file.Write([]byte(info))
	if nil != err {
		return err
	}
	return nil
}
