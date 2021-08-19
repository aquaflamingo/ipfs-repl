package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Variable struct {
	Name  string
	Value string
}

type Runnable func(interface{}) (error, string)

type ShellCommand struct {
	Name        string
	Description string
	Identifier  string
	Run         Runnable
}

var locals []*Variable

const (
	add ShellCommand = &ShellCommand{
		Name:        "Add",
		Description: "Adds a piece of content to IPFS network",
		Identifier:  "add",
		Run: func(filePath string) (error, string) {
			return nil, nil
		},
	}

	cat ShellCommand = &ShellCommand{
		Name:        "Cat",
		Description: "Reads an IPFS CID",
		Identifier:  "cat",
		Run: func(cid string) (error, string) {
			return nil, nil
		},
	}

	ls ShellCommand = &ShellCommand{
		Name:        "List",
		Description: "Lists currently defined variables in shell",
		Identifier:  "ls",
		Run: func() (error, string) {
			return nil, nil
		},
	}

	p ShellCommand = &ShellCommand{
		Name:        "Print",
		Description: "Prints the variable identifier's value",
		Identifier:  "p",
		Run: func() (error, string) {
			return nil, nil
		},
	}

	define ShellCommand = &ShellCommand{
		Name:        "Define",
		Description: "Defines a variable with the value given",
		Identifier:  "define",
		Run: func() (error, string) {
			return nil, nil
		},
	}
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("IPFS Shell: ")
	fmt.Println("---------------------")

	for {
		fmt.Print("ipfs > ")
		text, _ := reader.ReadString('\n')

		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)

		tokens := strings.Fields(text)

		cmd, err := parseShellCommand(tokens[0])

		if err != nil {
			// not a shell command
			// try REPL
			cmd, err := parseAliasCommands(tokens)

		} else {
			// is a shell command
			// exec shell command
		}
	}
}

func parseAliasCommands(tokens []string) (string, error) {
	// if <identifier> =
	// 		> Define new local variable
	// elsif <identifier>
	// 		> Print new local variable
	// else
	// 		> Error
}

func parseShellCommand(str string) (ShellCommand, error) {
	switch str {
	case add.Identifier:
		return add, nil
	case cat.Identifier:
		return cat, nil
	case ls.Identifier:
		return ls, nil
	case p.Identifier:
		return p, nil
	case define.Identifier:
		return define, nil
	default:
		return "", error.Error("Invalid shell command")
	}
}
