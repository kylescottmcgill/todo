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
	done_mark1 = "\u2610"
	done_mark2 = "\u2611"
)

func make_cmd_list(filename string) *commander.Command {
	cmd_list := func(cmd *commander.Command, args []string) error {
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
			if strings.HasPrefix(line, "-") {
				if !nflag {
					fmt.Printf("%s %03d: %s\n", done_mark2, n, strings.TrimSpace(string(line[1:])))
				}
			} else {
				fmt.Printf("%s %03d: %s\n", done_mark1, n, strings.TrimSpace(line))
			}
			n++

		}
		return nil
	}

	flg := *flag.NewFlagSet("list", flag.ExitOnError)
	flg.Bool("n", false, "only not done")

	return &commander.Command{
		Run:       cmd_list,
		UsageLine: "list [options]",
		Short:     "show list index",
		Flag:      flg,
	}
}
