package main

import (
	"fmt"
	"github.com/mihard/selective-build-helper/action"
	"github.com/mihard/selective-build-helper/vcs"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringSliceFlag{
			Name:  "trigger",
			Usage: "list directories which trigger full build",
		},
		cli.StringSliceFlag{
			Name:  "dir",
			Usage: "list all directories for the full build",
		},
		cli.StringSliceFlag{
			Name:  "xdir",
			Usage: "list all directories to exclude from the build",
		},
		cli.StringFlag{
			Name:  "root",
			Usage: "optional, project root",
		},
		cli.StringFlag{
			Name:  "base",
			Usage: "base path",
		},
		cli.StringFlag{
			Name:  "commit",
			Usage: "use a specific commit",
		},
	}

	app.Action = func(c *cli.Context) (err error) {
		rp := c.String("root")
		if rp == "" {
			rp, err = os.Getwd()
			if err != nil {
				return errors.Wrapf(err, "Unable to understand, where do we are")
			}
		}

		bp := c.String("base")
		git := vcs.MakeGit(rp, bp)
		commit := c.String("commit")

		uniqueFolders, err := action.CollectDirectories(bp, commit, git)
		if err != nil {
			return errors.Wrapf(err, "Unable to collect changed directories from VCS")
		}

		triggers := c.StringSlice("trigger")
		directories := c.StringSlice("dir")

		if len(directories) < 1 {
			directories = []string{}
			_directories, err := action.CollectAllDirectories(rp, bp)
			if err != nil {
				return errors.Wrapf(err, "Unable to collect all directories to build")
			}

			for _, d := range _directories {
				exclude := false
				for _, xd := range c.StringSlice("xdir") {
					if d == xd {
						exclude = true
						break
					}
				}
				if !exclude {
					directories = append(directories, d)
				}
			}
		}

		fullBuild := false

		for _, p := range uniqueFolders {
			for _, tr := range triggers {
				if tr == p {
					fullBuild = true
					break
				}
			}
		}

		if fullBuild {
			for _, p := range directories {
				fmt.Println(p)
			}
		} else {
			for p := range uniqueFolders {
				fmt.Println(p)
			}
		}

		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
