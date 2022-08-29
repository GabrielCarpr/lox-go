package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	args := os.Args

	if len(args) > 2 {
		fmt.Print("Usage: glox [script]\n")
	} else if len(args) == 2 {
		runFile(args[1])
	} else {
		runPrompt()
	}
}

func runFile(path string) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		report(err, 0)
		os.Exit(1)
	}

	lox := NewLox()
	err = lox.Run(string(content))
	if err != nil {
		os.Exit(1)
	}
}

func runPrompt() {
	lox := NewLox()
	for {
		fmt.Print("> ")
		bio := bufio.NewReader(os.Stdin)
		line, err := bio.ReadString('\n')
		if err != nil || line == "" {
			report(err, 0)
		}
		lox.Run(line)
	}
}

func report(message error, line int) {
	errorName := "Error"
	if loxError, ok := message.(LoxError); ok {
		errorName = loxError.Type()
	}

	if line != 0 {
		fmt.Printf("[line %d] %s: %s\n", line, errorName, message.Error())
	} else {
		fmt.Printf("%s: %s\n", errorName, message.Error())
	}
}
