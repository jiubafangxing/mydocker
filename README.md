# MyDocker

一个简单的Docker实现示例，用于学习容器技术。

## 项目结构

```
mydocker/
├── cmd/                    # 命令行工具
├── pkg/
│   └── container/         # 容器核心功能
│       ├── container_linux.go   # Linux版本实现
│       └── container_darwin.go  # macOS版本实现
├── internal/              # 内部包
├── main.go               # 主程序入口
├── go.mod                # Go模块定义
├── go.work               # Go工作空间配置
└── README.md             # 项目说明
```

## 功能特性

- ✅ 基本的容器命令框架
- ✅ 跨平台支持 (Linux/macOS)
- ✅ 命令行接口 (基于urfave/cli)
- ✅ 结构化日志 (基于logrus)
- ✅ Go模块管理

## 编译和运行

### 编译项目

```bash
# 设置Go代理（如果在国内网络环境）
export GOPROXY=https://goproxy.cn,direct

# 下载依赖
go mod tidy

# 编译
go build -o mydocker .
```

### 运行程序

```bash
# 查看帮助
./mydocker --help

# 查看run命令帮助
./mydocker run --help

# 示例运行（注意：在macOS上namespace功能有限）
./mydocker run -ti /bin/sh
```

## 主要依赖

- `github.com/urfave/cli/v3` - 命令行框架
- `github.com/sirupsen/logrus` - 结构化日志库
- `golang.org/x/sys/unix` - 系统调用接口

## 注意事项

1. **Linux vs macOS**: 
   - Linux版本支持完整的namespace功能
   - macOS版本提供基础的进程执行功能（macOS不支持Linux容器namespace）

2. **权限要求**:
   - 在Linux上运行需要适当的权限来创建namespace

3. **开发状态**:
   - 这是一个学习项目，不应该在生产环境中使用
   - 容器功能还在持续开发中

## 下一步开发计划

- [ ] 添加cgroup资源限制
- [ ] 实现容器文件系统隔离
- [ ] 添加网络隔离功能
- [ ] 支持容器镜像管理
- [ ] 添加更多的容器运行时选项

## 许可证

MIT License
