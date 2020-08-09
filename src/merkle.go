package main

import (
	b64 "encoding/base64"
	"fmt"
	"golang.org/x/crypto/sha3"
	"io"
	"os"
)

const BufferSize = 64

const (
	Yellow  = "\033[1;33m%s\033[0m"
	Blue    = "\033[0;34m%s\033[0m"
	Magenta = "\033[0;35m%s\033[0m"
	Cyan    = "\033[0;36m%s\033[0m"
)

func prettyPrint(fingerPrint string) {
	for i := 0; i < 4; i++ {
		startIdx := i * 22
		endIdx := (i + 1) * 22
		sum := 0
		for j := startIdx; j < endIdx; j++ {
			sum += int(fingerPrint[j])
		}
		remainder := sum % 4
		output := fingerPrint[startIdx:endIdx]
		var color string
		switch remainder {
		case 0:
			color = Yellow
		case 1:
			color = Blue
		case 2:
			color = Magenta
		case 3:
			color = Cyan
		}
		fmt.Printf(color, output)
		fmt.Println()
	}
}

func fingerPrint(table [][]byte) []byte {
	if len(table) == 1 {
		return table[0]
	}
	hasher := sha3.New512()
	children := make([][]byte, 0)
	for i, j := 0, 1; j < len(table); i, j = i+1, j+1 {
		hasher.Write(table[i])
		hasher.Write(table[j])
		hashValue := hasher.Sum(nil)
		hasher.Reset()
		children = append(children, hashValue)
	}
	return fingerPrint(children)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Not enough arguments!")
		return
	}
	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	hasher := sha3.New512()
	defer file.Close()
	buffer := make([]byte, BufferSize)
	children := make([][]byte, 0)
	for {
		bytesRead, err := file.Read(buffer)
		if err != nil {
			if err != io.EOF {
				panic(err)
			}
			break
		}
		if bytesRead < 64 {
			for i := range buffer {
				if i >= bytesRead {
					buffer[i] = 0
				}
			}
		}
		hasher.Write(buffer)
		hashValue := hasher.Sum(nil)
		hasher.Reset()
		children = append(children, hashValue)
	}
	if len(children)%2 == 1 {
		for i := range buffer {
			buffer[i] = 0
		}
		hasher.Write(buffer)
		hashValue := hasher.Sum(nil)
		hasher.Reset()
		children = append(children, hashValue)
	}
	prettyPrint(b64.StdEncoding.EncodeToString(fingerPrint(children)))
}
