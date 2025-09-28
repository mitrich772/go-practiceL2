package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	args := os.Args[1:]
	fmt.Println("Args:")
	for _, arg := range args {
		fmt.Print(arg + " ")
	}
	fmt.Println()

	// if _, err := os.Create("data.txt"); err != nil {
	// 	fmt.Printf("file not created: %s", err)
	// 	os.Exit(1)
	// }
	file, err := os.Open("data.txt")
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	defer file.Close()
	data := make([]byte, 64)
	for {
		n, err := file.Read(data)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Ошибка чтения:", err)
			break
		}
		fmt.Print(string(data[:n]))
	}

}
