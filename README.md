# IPFS REPL
A bare bones REPL for IPFS 

## Install
You can build the binary by running make build
```
make build
```

Code is currently setup to use `localhost:5001` as node source. 

## Usage
### Define
You can define variables in the REPL through standard assignment
```bash
myVar = 123
```

These can be used in other commands

### Print
You can print the value of a variable by typing its name
```bash
myVar
123
```

### Add 
You can add a file to IPFS by using `add` with the full path to the file to add
```bash
ipfs > add /Users/myCoolUser/ipfs-repl/test.txt

QmZLRFWaz9Kypt2ACNMDzA5uzACDRiCqwdkNSP1UZsu56D
```

The returned value is the content identifier hash for IPFS. You can use the 'special' `$$` variable to reprint the value too.
```bash
ipfs > myVar = $$

ipfs > myVar
QmZLRFWaz9Kypt2ACNMDzA5uzACDRiCqwdkNSP1UZsu56D
```

### Cat
You can use `cat` to print the value of the file from IPFS:
```bash
ipfs > cat $$
This is a test file stored on IPFS.
```

## Licence
This code is licensed under MIT open source licence.


