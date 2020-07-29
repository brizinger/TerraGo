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
	text, _ := consoleReader.ReadString('\n')	// Read until new line (pressed enter)
	text = strings.Replace(text, " ", "", -1)	// Remove whitespace
	text = strings.TrimSuffix(text, "\n")	// Remove new line at end
	
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
	fmt.Printf("\nLocation of %s: %s\n", file.Name(), OriginalDirectory())
}

// ReadFile reads the requested file and puts it in a string
func ReadFile(file string) string {
	b, err := wr.ReadFile(file)
			check(err)
		s := string(b)
		return s
}

//AddQuotes to string and return it
func AddQuotes(text string) string{
	return string('"') + text + string('"')
}

func dockerProviderCode(host string) {
	// Re-writes the old file
	// TODO: backup the old file if there is any
	text := "provider \"docker\" { \n	host = " + host + " \n}"
	line1:= []byte(text) 
	wr.WriteFile("main.tf", line1 , 0666)
}

func dockerContainerCode(containers string, originalDir string, mainDir string, network bool, networkName string, ports bool, extP string, intrP string) {

	os.Chdir(originalDir)

	if strings.Contains(containers, ",") { // More than one container
		containersSplit := strings.Split(containers, ",")

		s := ReadFile("docker_containers.txt")

		if network {
		//Docker Network is added here as we have already read the docker_containers.txt and can change the dir
		dockerAddNetworkInterface(networkName, mainDir)
		}
		// TODO: Some containers may have a different name from their image. 
		for i := 0; i < len(containersSplit); i++ {

			if(strings.Contains(s, containersSplit[i])) {

				dockerAddContainerMain(containersSplit[i], mainDir, network, networkName, ports, extP, intrP)

				dockerAddContainerImage(containersSplit[i], mainDir)
			}
		}
	}else {
		s := ReadFile("docker_containers.txt")

		if network {
			//Docker Network is added here as we have already read the docker_containers.txt and can change the dir
			dockerAddNetworkInterface(networkName, mainDir)
			}

		if(strings.Contains(s, containers)) {
			//TODO: PORTS
			dockerAddContainerMain(containers, mainDir, network, networkName, ports, extP, intrP)

			dockerAddContainerImage(containers, mainDir) 
		}
	}
}

func dockerAddContainerImage(image string, mainDir string) {
	if OriginalDirectory() != mainDir{
		os.Chdir(mainDir)
	}

	imageQuote := string('"') + image + string('"')

	text := "\n resource \"docker_image\" " + imageQuote + " { \n	name = " + imageQuote + " \n} \n"
	fmt.Printf("\n Writing image for: %s!", image)
	defer AppendFile("main.tf", text)
}

func dockerAddContainerMain(image string, mainDir string, network bool, networkName string, ports bool, extP string, intrP string) {

	if OriginalDirectory() != mainDir {
		os.Chdir(mainDir)
	}

	imageQuote := string('"') + image + string('"')
	
	text := "\n resource \"docker_container\" " + imageQuote + " { \n	name = " + imageQuote + " \n	image = docker_image."+ image + ".latest \n" 
	fmt.Printf("\n Writing container block for container: %s!", image)
	if network {
		networkNameQuote := string('"') + networkName + string('"')
		text += "	networks_advanced { \n		name = " + networkNameQuote + "\n } \n"
		fmt.Printf("\n Writing networks_advanced block to: %s!", image)
	}
	if ports {
		text = dockerAddPortsToContainer(extP, intrP, text)
	}
	text += "} \n"
	
	defer AppendFile("main.tf", text)
}

func dockerAddNetworkInterface(name string, mainDir string) {
	if OriginalDirectory() != mainDir{
		os.Chdir(mainDir)
	}
	
	nameQuote := string('"') + name + string('"')

	text := "\n resource \"docker_network\" " + nameQuote + " { \n	name = " + nameQuote + "\n } \n"
	fmt.Printf("\n Writing Docker Network Resource!")
	defer AppendFile("main.tf", text)
} 

func dockerAddPortsToContainer(ext string, intr string, text string) string {
	text += "	ports { \n" + "		external = " + AddQuotes(ext) + "\n		internal = " + AddQuotes(intr) + "\n }\n"
	fmt.Printf("\n Writing Ports!")
	return text
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
		text := getTextFromTerminal()
		text = string('"') + text + string('"')
		dockerProviderCode(text)

		fmt.Print("What images will you be using? (separate with comma): ")
		text = getTextFromTerminal()

		fmt.Printf("Do you want to add a network interface?(y/n): ")
		network := getTextFromTerminal()
		if network == "y" || network == "yes" {
			fmt.Printf("Name of interface: ")
			network = getTextFromTerminal()
			//TODO: User input for external and internal ports
			dockerContainerCode(text, originalDir, mainDir, true, network, true, "8080", "80")
		}else {
			dockerContainerCode(text, originalDir, mainDir, false, "", true,"8080", "80")
		}

	default: 

	fmt.Printf("Error")
	}
}
