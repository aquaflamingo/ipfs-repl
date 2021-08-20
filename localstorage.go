package main

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
