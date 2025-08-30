package main

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"

	"mydocker/pkg/container"
)

const usage = "this is a demo docker "

// Run 执行具体command
func Run(tty bool, command string) {
	parent := container.NewParentProcess(tty, command)
	if err := parent.Start(); err != nil {
		log.Error(err)
	}
	parent.Wait()
	os.Exit(-1)
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
		cmd := ctx.Args().First()
		tty := ctx.Bool("ti")

		log.Infof("enable tty %v", tty)
		log.Infof("command %s", cmd)

		Run(tty, cmd)
		return nil
	},
}

var initCommand = cli.Command{
	Name:  "init",
	Usage: "Init container process run user's process in container. Do not call it outside",
	Action: func(ctx *cli.Context) error {
		log.Info("init come on")
		cmd := ctx.Args().First()
		log.Infof("command %s", cmd)
		err := container.RunContainerInitProcess(cmd, ctx.Args().Slice())
		return err
	},
}
