package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gonuts/commander"
	"github.com/gonuts/flag"
)

const (
	red       = "\033[1;31m"
	green     = "\033[1;32m"
	yellow    = "\033[1;33m"
	gray      = "\033[1;90m"
	reset     = "\033[0m"
	doneMark1 = red + "\u2610" + reset
	doneMark2 = green + "\u2611" + reset
)

func makeCmdList(filename string) *commander.Command {
	cmdList := func(cmd *commander.Command, args []string) error {
		nflag := cmd.Flag.Lookup("n").Value.Get().(bool)
		f, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer f.Close()
		br := bufio.NewReader(f)
		n := 1
		for {
			b, _, err := br.ReadLine()
			if err != nil {
				if err != io.EOF {
					return err
				}
				break
			}
			line := string(b)

			if strings.Contains(line, "[") {
				var h strings.Builder

				for _, s := range strings.Split(line, " ") {
					if strings.HasPrefix(s, "[") {
						h.WriteString(fmt.Sprintf("%s%s%s ", gray, s, reset))
					} else {
						h.WriteString(fmt.Sprintf("%s ", s))
					}
					line = h.String()
				}

				if len(args) > 0 {
					for _, t := range args {
						line = strings.ReplaceAll(line, fmt.Sprintf("[%s", t), fmt.Sprintf("%s[%s", yellow, t))
					}
				}
			}

			if strings.HasPrefix(line, "-") {
				line = fmt.Sprintf("%s %03d: %s", doneMark2, n, strings.TrimSpace(line[1:]))
			} else {
				line = fmt.Sprintf("%s %03d: %s", doneMark1, n, strings.TrimSpace(line))
			}

			if !nflag {
				fmt.Println(line)
			}

			n++
		}
		return nil
	}

	flg := *flag.NewFlagSet("list", flag.ExitOnError)
	flg.Bool("n", false, "only not done")

	return &commander.Command{
		Run:       cmdList,
		UsageLine: "list [options]",
		Short:     "show list index",
		Flag:      flg,
	}
}
