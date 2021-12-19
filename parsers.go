package main

import "fmt"

func ParseCommandArguments(cmd *ShellCommand, tokens []string) ([]string, error) {
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

func ParseShellCommand(tokens []string) (*ShellCommand, error) {
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
