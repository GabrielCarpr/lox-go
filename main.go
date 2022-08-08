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

	}
	source := string(content)

	run(source)
}

func runPrompt() {
	for {
		fmt.Print("> ")
		bio := bufio.NewReader(os.Stdin)
		line, err := bio.ReadString('\n')
		if err != nil || line == "" {
			report(err, 0)
		}
		err = run(line)
		if err != nil {
			report(err, 0)
		}
	}
}

func run(source string) error {
	scanner := NewScanner(source)
	scanner.Scan()
	tokens := scanner.tokens

	for _, token := range tokens {
		fmt.Println(token.String())
	}

	return nil
}

func report(message error, line int) {
	if line != 0 {
		fmt.Printf("[line %d] Error: %s\n", line, message.Error())
	} else {
		fmt.Printf("Error: %s\n", message.Error())
	}
}
