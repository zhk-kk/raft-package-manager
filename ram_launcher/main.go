package main

// const APPS_SYMLINKS_SUB_PATH = "apps"
// const DATABASE_SUB_PATH = "database"

func main() {

	// userHomeDir, err := os.UserHomeDir()
	// if err != nil {
	// 	log.Fatalln("Unable to obtain the path of the user home directory.")
	// }

	// ramlRunPath, err := filepath.Abs(os.Args[0])
	// if err != nil {
	// 	log.Fatalln("Unable to obtain the path of the file, from which the raml was ran.")
	// }
	// ramlExePath, err := os.Executable()
	// if err != nil {
	// 	log.Fatalln("Unable to obtain the raml executable path.")
	// }

	// isRamlDirectRun := false
	// if ramlRunPath == ramlExePath {
	// 	isRamlDirectRun = true
	// }

	// ramStorePath := path.Join(userHomeDir, ".local", "share", "ram-store")
	// if err != nil {
	// 	log.Fatalln("Unable to obtain the ram store path.")
	// }

	// if _, err := os.Stat(ramStorePath); errors.Is(err, os.ErrNotExist) {
	// 	fmt.Println("`ramstore` doesn't seem to have been initialized. Creating the ramstore.")
	// 	initRamStore(ramStorePath)
	// func initRamStore(ramStorePath string) {
	// 	os.MkdirAll(ramStorePath, os.ModePerm)
	// 	os.MkdirAll(path.Join(ramStorePath, APPS_SYMLINKS_SUB_PATH), os.ModePerm)
	// 	os.MkdirAll(path.Join(ramStorePath, DATABASE_SUB_PATH), os.ModePerm)
	// }
	// }

	// fmt.Println("Raml was run from:     ", ramlRunPath)
	// fmt.Println("Raml executable is at: ", ramlExePath)
	// fmt.Println(isRamlDirectRun)
}

// func initRamStore(ramStorePath string) {
// 	os.MkdirAll(ramStorePath, os.ModePerm)
// 	os.MkdirAll(path.Join(ramStorePath, APPS_SYMLINKS_SUB_PATH), os.ModePerm)
// 	os.MkdirAll(path.Join(ramStorePath, DATABASE_SUB_PATH), os.ModePerm)
// }
