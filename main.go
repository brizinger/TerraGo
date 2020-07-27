package main

import (
	"strings"
	"bufio"
	"os"
	"fmt"
	wr "io/ioutil";
	tm "github.com/buger/goterm" // Will be used for better looking terminal.
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func getDirectory() {
	consoleReader := bufio.NewReader(os.Stdin)
	fmt.Print("Project Directory: ")
	dir, _ := consoleReader.ReadString('\n')
	dir = strings.TrimSuffix(dir, "\n")
	
	fmt.Println(dir)
	
	/*originalDir, err := os.Getwd() -------- Original Dir of needed later
	check(err)
	*/
	os.Chdir(dir)
}

func dockerProviderCode() {
	line1:= []byte("provider \"docker\" { \n host = \"tcp://127.0.0.1:2376\" \n}")
	wr.WriteFile("main.tf", line1 , 0666)
}

func awsProviderCode() {
	
}

func main() {

	version := "0.0.1"
	fmt.Println(tm.Background(tm.Color(tm.Bold("TerraGo"), tm.RED), tm.BLACK))
	fmt.Println("Version:", version)

	fmt.Println("For help go to: https://github.com/brizinger/TerraGo")

	// TODO: Check if file is there, if file is there, create a main.tf.2
	// TODO: Check if directory is correct
	// TODO: Add possibility to also use terraform commands from within
	
	getDirectory()

	f, err := os.Create("main.tf")
	check(err)
	
	defer f.Close()

	fmt.Print("Provider: ")

	var provider string

	fmt.Scan(&provider)

	switch provider {
	case "docker":
		
		dockerProviderCode()

	case "aws":
		awsProviderCode()
	}
}
