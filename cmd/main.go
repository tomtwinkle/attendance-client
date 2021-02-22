package main

import (
	"context"
	"fmt"
	attendanceclient "github.com/tomtwinkle/attendance-client/cli"
	"github.com/tomtwinkle/attendance-client/config"
	bugyoclient "github.com/tomtwinkle/bugyo-client-go"
	"github.com/urfave/cli"
	"log"
	"os"
)

var version = "unknown"
var revision = "unknown"

func main() {
	app := cli.NewApp()
	app.Name = "Attendance Client"
	app.Usage = "奉行クラウド打刻クライアント"
	app.Author = "tomtwinkle"
	app.Version = fmt.Sprintf("attendance-client cli version %s.rev-%s", version, revision)
	app.Commands = []cli.Command{
		{
			Name:  "init",
			Usage: "設定ファイルの作成",
			Action: func(c *cli.Context) error {
				cfg := config.NewConfig()
				if _, err := cfg.Init(); err != nil {
					return err
				}
				return nil
			},
		},
		{
			Name:      "punchmark",
			ShortName: "pm",
			Usage:     "タイムレコーダーの打刻を行う",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "type, t",
					Usage: "出勤: --type in or -t in" +
						"\n\t退出: --type out or -t out" +
						"\n\t外出: --type go or -t go" +
						"\n\t再入: --type return or -t return",
					Required: true,
					Value:    "",
				},
				cli.BoolFlag{
					Name:     "slackonly, so",
					Usage:    "打刻は行わずslackにのみ通知する --slackonly or -so",
					Required: false,
				},
				cli.BoolFlag{
					Name:     "slackskip, ss",
					Usage:    "slackに通知をせず打刻のみ行う --slackskip or -ss",
					Required: false,
				},
			},
			Action: func(c *cli.Context) error {
				ctx := context.Background()
				acli := attendanceclient.NewCLI()
				slackonly := c.Bool("slackonly")
				slackskip := c.Bool("slackskip")
				switch c.String("type") {
				case "in":
					return acli.PunchMark(ctx, bugyoclient.ClockTypeClockIn, slackonly, slackskip)
				case "out":
					return acli.PunchMark(ctx, bugyoclient.ClockTypeClockOut, slackonly, slackskip)
				case "go":
					return acli.PunchMark(ctx, bugyoclient.ClockTypeGoOut, slackonly, slackskip)
				case "return":
					return acli.PunchMark(ctx, bugyoclient.ClockTypeReturned, slackonly, slackskip)
				default:
					return cli.ShowSubcommandHelp(c)
				}
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
