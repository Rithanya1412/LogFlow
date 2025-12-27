package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	f, err := os.Open("app.log")
	if err != nil {
		fmt.Println("error opening file", err)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	fmt.Println("hello World")
}
