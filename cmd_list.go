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
	doneMark1 = "\033[1;31m\u2610\033[0m"
	doneMark2 = "\033[1;32m\u2611\033[0m"
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
			if len(args) > 0 {
				for _, t := range args {
					line = strings.ReplaceAll(line, "#"+t, fmt.Sprintf("\033[1;31m#%s\033[0m", t))
				}
			}
			if strings.HasPrefix(line, "-") {
				if !nflag {
					fmt.Printf("%s %03d: %s\n", doneMark2, n, strings.TrimSpace(line[1:]))
				}
			} else {
				fmt.Printf("%s %03d: %s\n", doneMark1, n, strings.TrimSpace(line))
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
