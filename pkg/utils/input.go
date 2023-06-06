package utils

import (
	"bufio"
	"os"
)

// ReadInput reads input from stdin then pushes the string result down the input channel
func ReadInput(inputChan chan string) {
	reader := bufio.NewReader(os.Stdin)
	for {
		command, _ := reader.ReadString('\n')
		inputChan <- command
		reader.Reset(os.Stdin)
	}
}
