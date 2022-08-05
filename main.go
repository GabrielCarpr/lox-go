package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	args := os.Args

	var err error
	if len(args) > 2 {
		fmt.Print("Usage: glox [script]\n")
	} else if len(args) == 2 {
		err = runFile(args[1])
	} else {
		err = runPrompt()
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runFile(path string) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	source := string(content)

	err = run(source)
	return err
}

func runPrompt() error {

	for {
		fmt.Print("> ")
		var line string
		_, err := fmt.Scanln(&line)
		if err != nil || line == "" {
			return err
		}
		err = run(line)
		if err != nil {
			report(err, 0)
		}
	}
}

func run(source string) error {
	fmt.Println(source)

	return nil
}

func report(message error, line int) {
	if line != 0 {
		fmt.Printf("[line %d] Error: %s\n", line, message.Error())
	} else {
		fmt.Printf("Error: %s\n", message.Error())
	}
}
