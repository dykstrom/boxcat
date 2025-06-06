package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/dykstrom/boxcat/internal/app/bci"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: boxcat <sourcefile>")
		return
	}
	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	var program []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		program = append(program, scanner.Text())
	}

	interpreter := bci.NewInterpreter(program)
	interpreter.Run()
}
