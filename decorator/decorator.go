package decorator

import (
	"bufio"
	"bytes"
	"log"
	"os"
)

func Deco(out *bytes.Buffer) {
	file, err := os.Open("./res/godown.css")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		out.WriteString(line)
		out.WriteString("\n")
	}
}
