package main

import (
	"bufio"
	"fmt"
	wr "io/ioutil"
	"os"
	"strings"

	tm "github.com/buger/goterm" // Will be used for better looking terminal.
	// yaml "gopkg.in/yaml.v2"
	//TODO: Add possibility for config files in yaml format
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

// GetDirectory - asks for project Directory and mounts it.
func GetDirectory() string {
	consoleReader := bufio.NewReader(os.Stdin)
	fmt.Print("Project Directory: ")
	dir, _ := consoleReader.ReadString('\n')
	dir = strings.TrimSuffix(dir, "\n")

	fmt.Println(dir)

	os.Chdir(dir)

	return dir
}

// GetTextFromTerminal - Returns the text entered by the user in the terminal as string
func GetTextFromTerminal(yesNo bool) string {
	consoleReader := bufio.NewReader(os.Stdin)
	text, _ := consoleReader.ReadString('\n') // Read until new line (pressed enter)
	text = strings.Replace(text, " ", "", -1) // Remove whitespace
	text = strings.TrimSuffix(text, "\n")     // Remove new line at end
	if yesNo {
		firstCharacter := strings.ToLower(text) //Make all lowercase
		if strings.Contains(firstCharacter, "yes") || strings.Contains(firstCharacter, "no") || firstCharacter == "y" || firstCharacter == "n" {
			firstCharacter = text[0:1] // returns only the first character (y or no)
			text = firstCharacter
		} else {
			err := fmt.Errorf("Incorrect answer on Yes/No question. Returning No")
			fmt.Println(err.Error())
			text = "n"
		}
	}

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
func AddQuotes(text string) string {
	return string('"') + text + string('"')
}

// TerraformFile - proides the generic structure for the .tf file
type TerraformFile struct {
	provider         string
	hostIP           string
	networkInterface bool
	networkName      string
	containers       []string
	ports            []string // example: ports[0] == 8080:80 (External Port:Internal Port)
	useNetwork       []bool   // Is the container using the network interface?
	name             []string
	image            []string
}

// The user input that will determine the main.tf file
func userInput() TerraformFile {
	fmt.Printf("\nWhat provider will you be using?: ")
	providerUser := GetTextFromTerminal(false)

	fmt.Printf("\nWhat host will you be using?: ")
	hostUser := GetTextFromTerminal(false)

	fmt.Printf("\nDo you need network interface? (y/n): ")
	networkInterfaceTerminal := GetTextFromTerminal(true)
	var networkInterface bool
	var networkNameTerminal string
	if networkInterfaceTerminal == "y" {
		networkInterface = true
		fmt.Printf("\nName of network: ")
		networkNameTerminal = GetTextFromTerminal(false)
	} else {
		networkInterface = false
	}

	fmt.Printf("\nContainers which you want to use (separate by comma): ")
	containers := GetTextFromTerminal(false)
	var containersSplit []string
	if strings.Contains(containers, ",") {
		containersSplit = strings.Split(containers, ",")
	} else {
		containersSplit = []string{containers}
	}

	var ports []string
	var containerNames []string
	var containerImages []string
	var usesNetwork []bool

	// Resizes the arrays to have the same length as containersSplit
	if len(containersSplit) >= 1 {

		portsM := make([]string, len(containersSplit))
		ports = portsM

		namesM := make([]string, len(containersSplit))
		containerNames = namesM

		images := make([]string, len(containersSplit))
		containerImages = images

		useOfNetworkM := make([]bool, len(containersSplit))
		usesNetwork = useOfNetworkM
	}

	for i := 0; i < len(containersSplit); i++ {

		fmt.Printf("\nWhat name do you want to give to %s?: ", containersSplit[i])
		containerNames[i] = GetTextFromTerminal(false)

		fmt.Printf("\nWhat image do you want to use for %s?: ", containersSplit[i])
		containerImages[i] = GetTextFromTerminal(false)

		if networkInterface {
			fmt.Printf("\nDo you want to add the container to the network? (y/n): ")
			networkInterfaceTerminal = GetTextFromTerminal(true)
			if networkInterfaceTerminal == "y" {
				usesNetwork[i] = true
			} else {
				usesNetwork[i] = false
			}
		}

		fmt.Printf("\nWhat ports do you want to use for %s?: ", containersSplit[i])
		ports[i] = GetTextFromTerminal(false)

	}

	tf := TerraformFile{provider: providerUser, hostIP: hostUser, networkInterface: networkInterface, networkName: networkNameTerminal, containers: containersSplit, ports: ports, useNetwork: usesNetwork, name: containerNames, image: containerImages}

	return tf
}

// Creates the main.tf file, needs the tf struct
func fileCreation(tf TerraformFile, originalDir string, mainDir string) {
	dockerProviderCode(tf.provider, tf.hostIP)
	if tf.networkInterface {
		dockerAddNetworkInterface(tf.networkName)
	}
	dockerContainerCode(tf.containers, tf.name, originalDir, mainDir, tf.networkInterface, tf.networkName, tf.ports, tf.image)
	fmt.Printf("\nMake sure to set-up any environmental variables needed to start using the container(s)!\n")
}

func dockerProviderCode(provider string, host string) {

	if _, err := os.Stat("main.tf"); err == nil { // File already exists
		os.Rename("main.tf", "main.tf.old") // Rename old file
		//BUG: Deletes contents of old file!
		text := "provider " + AddQuotes(provider) + " { \n	host = \"tcp://" + host + "\" \n}"
		line1 := []byte(text)
		wr.WriteFile("main.tf", line1, 0666) // Write new file
	} else if os.IsNotExist(err) {
		text := "provider " + AddQuotes(provider) + " { \n	host = \"tcp://" + host + "\" \n}"
		line1 := []byte(text)
		wr.WriteFile("main.tf", line1, 0666)
	} else { // File may or may not exist. Could be permission issues or something else.
		fmt.Printf(err.Error())
	}

}

func dockerContainerCode(containers []string, containerNames []string, originalDir string, mainDir string, network bool, networkName string, ports []string, image []string) {

	os.Chdir(originalDir)

	for i := 0; i < len(containers); i++ {

		port := strings.Split(ports[i], ":")

		dockerAddContainerMain(containers[i], image[i], containerNames[i], mainDir, network, networkName, port[0], port[1])

		dockerAddContainerImage(image[i], mainDir)
	}
}

func dockerAddContainerImage(image string, mainDir string) {
	if OriginalDirectory() != mainDir {
		os.Chdir(mainDir)
	}

	text := "\n resource \"docker_image\" " + AddQuotes(image) + " { \n	name = " + AddQuotes(image) + " \n} \n"
	fmt.Printf("\n Writing image for: %s!", image)
	defer AppendFile("main.tf", text)
}

func dockerAddContainerMain(container string, image string, containerName string, mainDir string, network bool, networkName string, extrP string, intrP string) {

	if OriginalDirectory() != mainDir {
		os.Chdir(mainDir)
	}
	//TODO: Add support for specific releases (maybe add option for advanced vs simple setup?)
	text := "\n resource \"docker_container\" " + AddQuotes(container) + " { \n	name = " + AddQuotes(containerName) + " \n	image = docker_image." + image + ".latest \n"
	fmt.Printf("\n Writing container block for container: %s!", image)
	if network {
		text += "	networks_advanced { \n		name = " + AddQuotes(networkName) + "\n } \n"
		fmt.Printf("\n Writing networks_advanced block to: %s!", image)
	}

	text = dockerAddPortsToContainer(extrP, intrP, text)

	text += "} \n"

	defer AppendFile("main.tf", text)
}

func dockerAddNetworkInterface(networkName string) {

	text := "\n resource \"docker_network\" " + AddQuotes(networkName) + " { \n	name = " + AddQuotes(networkName) + "\n } \n"
	fmt.Printf("\n Writing Docker Network Resource!")
	defer AppendFile("main.tf", text)
}

func dockerAddPortsToContainer(ext string, intr string, text string) string {
	text += "	ports { \n" + "		external = " + AddQuotes(ext) + "\n		internal = " + AddQuotes(intr) + "\n }\n"
	fmt.Printf("\n Writing Ports!")
	return text
}

func main() {

	version := "0.0.1"
	fmt.Println(tm.Background(tm.Color(tm.Bold("TerraGo"), tm.RED), tm.BLACK))
	fmt.Println("Version:", version)

	fmt.Println("For help go to: https://github.com/brizinger/TerraGo")

	originalDir := OriginalDirectory()
	mainDir := GetDirectory()
	fmt.Printf("Project Directory: " + mainDir)
	f, err := os.Create("main.tf")
	check(err)

	defer f.Close()

	var tf TerraformFile
	tf = userInput()
	fmt.Printf("Do you want to create a file with the following settings? (y/n):")
	fmt.Println(tf)
	answer := GetTextFromTerminal(true)

	if answer == "y" {
		fileCreation(tf, originalDir, mainDir)
	}
}
