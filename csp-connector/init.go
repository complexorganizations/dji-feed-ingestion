package main

var (
	applicationConfigFile = "config.json"
)

func init() {
	// Check if the config file exists in the current directory
	if !fileExists(applicationConfigFile) {
		// Write a config file in the current directory if it doesn't exist
		writeToFile(applicationConfigFile, []byte(``))
	}
	// Check if the config file has not been modified
	if sha256OfFile(applicationConfigFile) == "d41d8cd98f00b204e9800998ecf8427e" {
		// The file has not been modified
		exitTheApplication("The config file has not been modified. Please modify it and try again.")
	}
	// Check if the required application are present in the system
	if commandExists("git") == false {
		exitTheApplication("Git is not installed in your system. Please install it and try again.")
	}
	// Check if the config has the correct format and all the info is correct.
}
