package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {
	cmd := exec.Command("nslookup")
	filePath := flag.String("path", "./configs/kube.txt", "Path to addresses")
	flag.Parse()
	file, err := os.Open(*filePath)
	check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var addresses strings.Builder
	for scanner.Scan() {
		address := scanner.Text()
		if address[0] == '$' {
			address = os.Getenv(address[1:])
		}
		addresses.WriteString(address)
		addresses.WriteString("\n")
	}

	cmd.Stdin = strings.NewReader(addresses.String())
	var out bytes.Buffer
	cmd.Stdout = &out

	err = cmd.Run()
	//check(err)
	//fmt.Printf("data:\n%s\n", out.String())
	fmt.Println("Yeeted this.")
}
