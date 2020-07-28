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
        fmt.Println("ERROR: ", e)
    }
}
// OriginalDirectory - finds the main.go directory and stores it in a variable in the main function.
func OriginalDirectory() string {
	originalDir, err := os.Getwd() // Original Dir
		check(err)

	return originalDir
}

func getDirectory() string {
	consoleReader := bufio.NewReader(os.Stdin)
	fmt.Print("Project Directory: ")
	dir, _ := consoleReader.ReadString('\n')
	dir = strings.TrimSuffix(dir, "\n")
	
	fmt.Println(dir)
	
	os.Chdir(dir)

	return dir
}

func getTextFromTerminal() string {
	consoleReader := bufio.NewReader(os.Stdin)
	text, _ := consoleReader.ReadString('\n')
	text = strings.TrimSuffix(text, "\n")
	return text
}
//AppendFile - Will write after the end of the last line
func AppendFile(maintffile string, text string) {     
    file, err := os.OpenFile(maintffile, os.O_WRONLY|os.O_APPEND, 0644)
    check(err)
	defer file.Close()

	len, err := file.WriteString(text)
    check(err)
    fmt.Printf("\nLength: %d bytes", len)
    fmt.Printf("\nFile Name: %s", file.Name())
}

func dockerProviderCode(host string) {
	// Re-writes the old file
	// TODO: backp the old file if there is any
	text := "provider \"docker\" { \n	host = " + host + " \n}"
	line1:= []byte(text) 
	wr.WriteFile("main.tf", line1 , 0666)
}

func dockerContainerCode(containers string, originalDir string, mainDir string) {

		os.Chdir(originalDir)
		
	if strings.Contains(containers, ",") { // More than one container
		containersSplit := strings.Split(containers, ",")

		b, err := wr.ReadFile("docker_containers.txt")
			check(err)
		s := string(b)

		// TODO: Some containers may have a different name from their image. 
		for i := 0; i < len(containersSplit); i++ {

			if(strings.Contains(s, containersSplit[i])) {
				
				dockerAddContainerImage(containersSplit[i], mainDir)
			}

		}
	}else {
		b, err := wr.ReadFile("docker_containers.txt")
			check(err)
		s := string(b)

		if(strings.Contains(s, containers)) {

		dockerAddContainerImage(containers, mainDir)

		}
	}
}

func dockerAddContainerImage(image string, mainDir string) {
	os.Chdir(mainDir)
	image = string('"') + image + string('"')
	text := "\n resource \"docker_image\" " + image + " { \n	name = " + image + " \n} \n" 
	defer AppendFile("main.tf", text)
}

func awsProviderCode() { // TODO
	
}

func main() {
	originalDir := OriginalDirectory()
	version := "0.0.1"
	fmt.Println(tm.Background(tm.Color(tm.Bold("TerraGo"), tm.RED), tm.BLACK))
	fmt.Println("Version:", version)

	fmt.Println("For help go to: https://github.com/brizinger/TerraGo")

	// TODO: Check if file is there, if file is there, create a main.tf.2
	// TODO: Check if directory is correct
	// TODO: Add possibility to also use terraform commands from within
	
	mainDir := getDirectory()

	f, err := os.Create("main.tf")
	check(err)
	
	defer f.Close()

	fmt.Print("Provider: ")

	var provider string

	fmt.Scan(&provider)

	switch provider {
	case "docker":
		
		fmt.Print("Host: ")
		text:= getTextFromTerminal()
		text = string('"') + text + string('"')
		dockerProviderCode(text)

		fmt.Print("What images will you be using? (separate with comma): ")
		text = getTextFromTerminal()
		strings.Replace(text, " ", "", -1)
		dockerContainerCode(text, originalDir, mainDir)

	case "aws":
		awsProviderCode()
	}
}
