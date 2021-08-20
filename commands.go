package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	ipfs "github.com/ipfs/go-ipfs-api"
)

var (
	cmdAdd *ShellCommand = &ShellCommand{
		Name:        "Add",
		Description: "Adds a piece of content to IPFS network",
		Identifier:  "add",
		Run: func(args []string, sh *ipfs.Shell) (string, error) {
			pathToFile := args[0]

			// Use variable
			setIfDefined(&pathToFile)

			file, err := os.Open(pathToFile)
			if err != nil {
				return "", fmt.Errorf("Unable to read file: %v", err)
			}
			reader := bufio.NewReader(file)

			cid, err := sh.Add(reader)
			if err != nil {
				return "", fmt.Errorf("error: %v", err)
			}
			// ipfs add
			res := fmt.Sprintf(cid)

			setShellVariable("$$", cid)

			return res, nil
		},
	}

	cmdCat *ShellCommand = &ShellCommand{
		Name:        "cmdCat",
		Description: "Reads an IPFS CID",
		Identifier:  "cat",
		Run: func(args []string, sh *ipfs.Shell) (string, error) {
			cid := args[0]
			// Use variable
			setIfDefined(&cid)

			reader, err := sh.Cat(cid)
			defer reader.Close()

			if err != nil {
				return "", err
			}

			buf := new(strings.Builder)
			_, err = io.Copy(buf, reader)
			if err != nil {
				return "", err
			}

			return buf.String(), nil
		},
	}

	cmdLs *ShellCommand = &ShellCommand{
		Name:        "List",
		Description: "Lists currently defined variables in shell",
		Identifier:  "ls",
		Run: func(args []string, sh *ipfs.Shell) (string, error) {
			if empty(locals) {
				// Don't print anything
				return "", nil
			}

			var str string

			for _, a := range locals {
				str = fmt.Sprintf("%s:%s\n", a.Name, a.Value)
			}

			return str, nil
		},
	}

	cmdPrint *ShellCommand = &ShellCommand{
		Name:        "Print",
		Description: "Prints the variable identifier's value",
		Identifier:  "p",
		Run: func(args []string, sh *ipfs.Shell) (string, error) {
			ident := args[0]

			i, found := contains(locals, ident)

			if !found {
				return "", fmt.Errorf("%s is not defined", ident)
			}

			return locals[i].Value, nil
		},
	}

	cmdDefine *ShellCommand = &ShellCommand{
		Name:        "Define",
		Description: "Defines a variable with the value given",
		Identifier:  "define",
		Run: func(args []string, sh *ipfs.Shell) (string, error) {
			ident := args[0]
			val := args[1]

			setShellVariable(ident, val)

			// Print the value to cli
			return val, nil
		},
	}

	cmdExit *ShellCommand = &ShellCommand{
		Name:        "Exit",
		Description: "Exits",
		Identifier:  "exit",
		Run: func(args []string, sh *ipfs.Shell) (string, error) {
			os.Exit(0)
			return "", nil
		},
	}
)
