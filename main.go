package main

import (
	"bufio"
	"fmt"
	"io"
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

		cmd, err := parseShellCommand(tokens)

		if err != nil {
			fmt.Printf("%v", err.Error())
		} else {
			args, err := parseCommandArguments(cmd, tokens)

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

func parseCommandArguments(cmd *ShellCommand, tokens []string) ([]string, error) {
	switch cmd.Identifier {
	case cmdAdd.Identifier:
		// Add file to ipfs
		if len(tokens) < 2 {
			return []string{}, fmt.Errorf("Error: no path to file given\n")
		}
		pathToFile := tokens[1]
		return []string{pathToFile}, nil
	case cmdCat.Identifier:
		// Cat file from ipfs
		if len(tokens) < 2 {
			return []string{}, fmt.Errorf("Error: no content identifier given\n")
		}

		cid := tokens[1]
		return []string{cid}, nil
	case cmdDefine.Identifier:
		// Define new variable in locals
		if len(tokens) < 3 {
			return []string{}, fmt.Errorf("Error: invalid assignment expression\n")
		} else if len(tokens) > 2 && tokens[1] != "=" {
			return []string{}, fmt.Errorf("Error: invalid assignment expression\n")
		}

		ident := tokens[0]
		val := tokens[2]
		return []string{ident, val}, nil
	case cmdLs.Identifier, cmdExit.Identifier:
		// List all variables in locals
		return []string{}, nil
	case cmdPrint.Identifier:
		// Print a variable
		return []string{tokens[0]}, nil
	default:
		return nil, fmt.Errorf("Invalid arguments")
	}
}

func parseAliasCommands(tokens []string) (*ShellCommand, error) {
	// If first token is only value, check if defined
	_, found := contains(locals, tokens[0])

	if len(tokens) == 1 && found {
		// This is a variable that is defined, print it
		return cmdPrint, nil
		// If three tokens check for defintion
	} else if len(tokens) == 3 && tokens[1] == "=" {
		return cmdDefine, nil
	} else {
		// No alias
		return nil, fmt.Errorf("Invalid shell command\n")
	}
}

func parseShellCommand(tokens []string) (*ShellCommand, error) {
	// if <identifier> =
	// 		> Define new local variable
	// elsif <identifier>
	// 		> Print new local variable
	// else
	// 		> Error
	identifer := tokens[0]

	switch identifer {
	case cmdAdd.Identifier:
		return cmdAdd, nil
	case cmdCat.Identifier:
		return cmdCat, nil
	case cmdLs.Identifier:
		return cmdLs, nil
	case cmdPrint.Identifier:
		return cmdPrint, nil
	case cmdDefine.Identifier:
		return cmdDefine, nil
	case cmdExit.Identifier:
		return cmdExit, nil
	default:
		return parseAliasCommands(tokens)
	}
}

func empty(arr []Variable) bool {
	return len(arr) < 1
}

func contains(arr []Variable, str string) (int, bool) {
	for i, a := range arr {
		if a.Name == str {
			return i, true
		}
	}
	return -1, false
}

func setIfDefined(strPtr *string) {
	// Use variable
	valIndex, valFound := contains(locals, *strPtr)

	if valFound {
		*strPtr = locals[valIndex].Value
	}
}

func setShellVariable(ident string, value string) {
	tmpVal := value

	setIfDefined(&tmpVal)

	index, found := contains(locals, ident)

	if found {
		// Set the variable
		locals[index] = Variable{Name: ident, Value: tmpVal}
	} else {
		// Add the variable
		locals = append(locals, Variable{Name: ident, Value: tmpVal})
	}
}
