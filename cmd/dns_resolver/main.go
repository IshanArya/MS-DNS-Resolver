package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func main() {
	cmd := exec.Command("nslookup")
	cmd.Stdin = strings.NewReader("google.com")
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("data:\n%s\n", out.String())
}
