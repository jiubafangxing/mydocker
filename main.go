package main

import (
	"context"
	"mydocker/pkg/cgroups"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"

	"mydocker/pkg/container"
)

const usage = "this is a demo docker "

// Run 执行具体command
func Run(tty bool, command []string, res *cgroups.ResourceConfig) {
	cmd, w := container.NewParentProcess(tty)
	if err := cmd.Start(); err != nil {
		log.Error(err)
	}
	cgroupManager, err := cgroups.NewV2CgroupManager("my-docker-cgroup")
	if err != nil {
		log.Errorf("Failed to create cgroup manager: %v", err)
	}
	cgroupManager.Set(res)
	cgroupManager.AddProcess(cmd.Process.Pid)
	sendInitCommand(command, w)
	cmd.Wait()
	os.Exit(-1)
}

func sendInitCommand(cmds []string, writePipe *os.File) {
	command := strings.Join(cmds, " ")
	log.Infof("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}

func main() {
	subCommands := []*cli.Command{
		&initCommand,
		&runCommand,
	}
	app := &cli.Command{
		EnableShellCompletion: true,
		Commands:              subCommands,
		Name:                  "mydocker",
		Usage:                 usage,
		Before: func(context *cli.Context) error {
			log.SetFormatter(&log.JSONFormatter{})
			log.SetOutput(os.Stdout)
			return nil
		},
	}
	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}

}

var runCommand = cli.Command{
	Name:  "run",
	Usage: `Create a container with namespace and cgroup limit mydocker run -ti [command]`,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "ti",
			Usage: "enable tty",
		},
	},
	Action: func(ctx *cli.Context) error {
		if ctx.Args().Len() < 1 {
			return cli.Exit("Missing container command", 1)
		}
		cmd := ctx.Args()
		tty := ctx.Bool("ti")

		log.Infof("enable tty %v", tty)
		log.Infof("command %s", cmd)
		var cmdArray []string
		for _, arg := range ctx.Args().Slice() {
			cmdArray = append(cmdArray, arg)
		}
		resConf := &cgroups.ResourceConfig{
			MemoryLimit: defaultIfEmpty(ctx.String("m"), "512m"),
			CpuSet:      defaultIfEmpty(ctx.String("cpuset"), "0"),
			CpuShare:    defaultIfEmpty(ctx.String("cpushare"), "1024"),
		}
		log.Infof("createTty %v", tty)
		//containerName := ctx.String("name")
		//volume := ctx.String("v")
		//network := ctx.String("net")
		//
		//envSlice := ctx.StringSlice("e")
		//portmapping := ctx.StringSlice("p")
		Run(tty, cmdArray[1:], resConf)
		return nil
	},
}

func defaultIfEmpty(val, def string) string {
	if val == "" {
		return def
	}
	return val
}

var initCommand = cli.Command{
	Name:  "init",
	Usage: "Init container process run user's process in container. Do not call it outside",
	Action: func(ctx *cli.Context) error {
		log.Info("init come on")
		cmd := ctx.Args().First()
		log.Infof("command %s", cmd)
		err := container.RunContainerInitProcess()
		return err
	},
}
