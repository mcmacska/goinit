package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func getLatestVersion() string {
	res, err := http.Get("https://go.dev/VERSION?m=text")
	if err != nil {
		fmt.Println("error making http request: ", err)
	}

	if res.StatusCode != 200 {
		fmt.Println("Could not get latest go version")
		return "1.22.2"
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	// trim the "go" from the beginning of the body, and trim everything after the "time"
	result := strings.Split(strings.TrimPrefix(string(body), "go"), "time")[0]

	fmt.Println("latest version: ", result)

	return result
}

func readFile(filename string) string {
	var dat, err = os.ReadFile(filename)

	if err != nil {
		fmt.Println("Could not read file: " + filename)
	}

	return string(dat)
}

func writeFile(path string, content string) {
	contentByteArray := []byte(content)

	err := os.WriteFile(path, contentByteArray, 0644)

	if err != nil {
		fmt.Println(err)
	}
}

// sets up a simple go project based on a config
func setProject(configString string) {
	var config map[string]string = make(map[string]string)

	var configLines []string = strings.Split(configString, "\n")

	for i := 0; i < len(configLines); i++ {
		var keyValue = strings.Split(string(configLines[i]), ":")
		if len(keyValue) < 2 {
			continue
		}

		config[keyValue[0]] = keyValue[1]
	}

	if config[name] != "" && config[path] != "" {
		// make project directory
		err := os.MkdirAll(config[path]+"/"+config[name], os.ModePerm)
		if err != nil {
			fmt.Println("Error creating directory: ", err)
			return
		} else {
			fmt.Println("Directory created successfully")
		}

		// get imports
		var toImport []string = strings.Split(config[imports], ",")

		// join the iports as a string with the correct format and apostrophes
		var toImportString = "\t\"" + strings.Join(toImport[:], "\"\n\t\"") + "\""

		// create a main.go inside directory
		var mainContent = "package main\n\nimport (\n" + toImportString + "\n)\n\nfunc main() {\n\tfmt.Println()\n}\n"

		writeFile(config[path]+"/"+config[name]+"/main.go", mainContent)

		// get latest stable go version
		var latesVersion = getLatestVersion()

		// create a go.mod
		var goModContent = "module " + config[name] + "\n\ngo " + latesVersion
		writeFile(config[path]+"/"+config[name]+"/go.mod", goModContent)
	}
}

func main() {
	fmt.Println("Initializing go project...")
	var config = readFile("config.txt")

	setProject(config)
}
