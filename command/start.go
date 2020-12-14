package command

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli/v2"
)

type startProjectConfig interface {
	GetProjectMainContainer() string
	GetStartCommand() string
	SaveContainerNameToProjectConfig(string) error
	SaveStartCommandToProjectConfig(string) error
}

type startDialog interface {
	SetMainContaner([]string) (int, string, error)
	SetStartCommand() (string, error)
}

func defineProjectMainContainer(cfg startProjectConfig, d startDialog, containerlist []string) (err error) {
	if cfg.GetProjectMainContainer() == "" {
		_, container, err := d.SetMainContaner(containerlist)

		if err != nil {
			return err
		}

		if container == "" {
			return errors.New("Container name is empty. Set the container name")
		}

		return cfg.SaveContainerNameToProjectConfig(container)
	}

	return err
}

func defineStartCommand(cfg startProjectConfig, d startDialog, containerlist []string) (err error) {
	if cfg.GetStartCommand() == "" {
		startCommand, err := d.SetStartCommand()

		if err != nil {
			return err
		}

		if startCommand == "" {
			return errors.New("Start command cannot be empty")
		}

		return cfg.SaveStartCommandToProjectConfig(startCommand)
	}

	return err
}

func runStartProject(c *cli.Context, cfg startProjectConfig, args []string) error {
	commandSlice := strings.Split(cfg.GetStartCommand(), " ")

	var binary = commandSlice[0]
	var initArgs = commandSlice[1:]

	extraInitArgs := c.Args().Slice()

	args = append(initArgs, args...)
	args = append(args, extraInitArgs...)

	log.Printf("Called: %s %s", binary, strings.Join(args, " "))

	cmd := exec.Command(binary, args...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// CallStartProjectBasic runs docker project
func CallStartProjectBasic(initf func(), cfg startProjectConfig, d startDialog, containerlist []string) *cli.Command {
	cmd := cli.Command{
		Name:            "start",
		Aliases:         []string{"st"},
		Usage:           `runs defined command: {docker-compose -f docker-compose.yml up} [custom parameters]`,
		Description:     `It's possible to use any custom parameters coming after "up"`,
		SkipFlagParsing: true,
		Action: func(c *cli.Context) (err error) {
			initf()

			if err = defineProjectMainContainer(cfg, d, containerlist); err != nil {
				return err
			}

			if err = defineStartCommand(cfg, d, containerlist); err != nil {
				return err
			}

			return runStartProject(c, cfg, []string{})
		},
	}

	return &cmd
}

// CallStartProjectForceRecreate runs docker project
func CallStartProjectForceRecreate(initf func(), cfg startProjectConfig, d startDialog, containerlist []string) *cli.Command {
	cmd := cli.Command{
		Name:    "start:force",
		Aliases: []string{"s:f"},
		Usage:   `runs defined command: {docker-compose -f docker-compose.yml up --force-recreat} [custom parameters]`,
		Description: `
		--force-recreate - Recreate containers even if their configuration and image haven't changed
		It's possible to use any custom parameters coming after "up"`,
		SkipFlagParsing: true,
		Action: func(c *cli.Context) (err error) {
			initf()

			if err = defineProjectMainContainer(cfg, d, containerlist); err != nil {
				return err
			}

			if err = defineStartCommand(cfg, d, containerlist); err != nil {
				return err
			}

			return runStartProject(c, cfg, []string{"--force-recreate"})
		},
	}

	return &cmd
}

// CallStartProjectOrphans runs docker project
func CallStartProjectOrphans(initf func(), cfg startProjectConfig, d startDialog, containerlist []string) *cli.Command {
	cmd := cli.Command{
		Name:    "start:orphans",
		Aliases: []string{"s:o"},
		Usage:   `runs defined command: {docker-compose -f docker-compose.yml up --remove-orphans} [custom parameters]`,
		Description: `
		--remove-orphans - Remove containers for services not defined in the Compose file
		It's possible to use any custom parameters coming after "up"`,
		SkipFlagParsing: true,
		Action: func(c *cli.Context) (err error) {
			initf()

			if err = defineProjectMainContainer(cfg, d, containerlist); err != nil {
				return err
			}

			if err = defineStartCommand(cfg, d, containerlist); err != nil {
				return err
			}

			return runStartProject(c, cfg, []string{"--remove-orphans"})
		},
	}

	return &cmd
}

// CallStartProjectForceOrphans runs docker project
func CallStartProjectForceOrphans(initf func(), cfg startProjectConfig, d startDialog, containerlist []string) *cli.Command {
	cmd := cli.Command{
		Name:    "start:force-orphans",
		Aliases: []string{"s:fo"},
		Usage:   `runs defined command: {docker-compose -f docker-compose.yml up --force-recreate --remove-orphans} [custom parameters]`,
		Description: `
		--force-recreate - Recreate containers even if their configuration and image haven't changed
		--remove-orphans - Remove containers for services not defined in the Compose file
		It's possible to use any custom parameters coming after "up"`,
		SkipFlagParsing: true,
		Action: func(c *cli.Context) (err error) {
			initf()

			if err = defineProjectMainContainer(cfg, d, containerlist); err != nil {
				return err
			}

			if err = defineStartCommand(cfg, d, containerlist); err != nil {
				return err
			}

			return runStartProject(c, cfg, []string{"--force-recreate", "--remove-orphans"})
		},
	}

	return &cmd
}
