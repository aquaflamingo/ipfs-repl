package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	ipfs "github.com/ipfs/go-ipfs-api"
)

type Variable struct {
	Name  string
	Value string
}

var locals []Variable

type Runnable func(args []string, sh *ipfs.Shell) (string, error)

type ShellCommand struct {
	Name        string
	Description string
	Identifier  string
	Run         Runnable
}

const localhost string = "localhost:5001"

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("IPFS REPL v0.1")

	sh := ipfs.NewShell(localhost)

	// REPL loop
	for {
		fmt.Print("ipfs > ")
		text, _ := reader.ReadString('\n')

		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)

		if text == "" {
			continue
		}

		tokens := strings.Fields(text)

		cmd, err := ParseShellCommand(tokens)

		if err != nil {
			fmt.Printf("%v", err.Error())
		} else {
			args, err := ParseCommandArguments(cmd, tokens)

			if err != nil {
				fmt.Printf("%v", err.Error())
			} else {
				res, err := cmd.Run(args, sh)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println(res)
				}
			}
		}
	}
}
