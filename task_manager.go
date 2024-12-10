package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

type AnyStruct interface{}

type Config struct {
	ListDirectory          string `json:"ListDirectory"`
	FileNameSameAsListName bool   `json:"FileNameSameAsListName"`
	DeleteCompletedTasks   bool   `json:"DeleteCompletedTasks"`
	SortAttribute          string `json:"SortAttribute"`
}

type TaskList struct {
	Name  string `json:"name"`
	Tasks []Task `json:"tasks"`
}

type Task struct {
	Name             string `json:"name"`
	Severity         int    `json:"severity"`
	CompletionStatus bool   `json:"completionStatus"`
}

var configData Config
var scanner *bufio.Scanner
var activeList TaskList
var activeFileName string

func main() {
	configFile, err := os.Open("./task_manager.config.json")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer configFile.Close()
	configFileInfo, err := configFile.Stat()
	if err != nil {
		log.Fatal(err)
		return
	}
	buffer := make([]byte, configFileInfo.Size())
	_, err = configFile.Read(buffer)
	if err != nil {
		fmt.Println("Error reading config file: ", err)
		return
	}

	err = json.Unmarshal(buffer, &configData)
	if err != nil {
		fmt.Println("Error unmarshalling json config data: ", err)
		return
	}

	os.Chdir(configData.ListDirectory)

	scanner = bufio.NewScanner(os.Stdin)
	fmt.Println("You are now running Dominic Camill's \"Task Manager\"\nThere are many operations you can perform. To learn about what you can do execute the help command.")
	for {
		runCommand()
	}
}

func help(args []string) {
	var name string
	if len(args) == 1 {
		name = ""
	} else {
		name = args[1]
	}

	fmt.Println("For more help with a specific command name, type \"help <command name>\"")
	switch name {
	case "help":
		fmt.Println("help: List out all commands with descriptions defining what the command does. Can be refined to a single more in depth description of the command using the format \"help <command-name>\". This will run the help command but will only give the description of the command defined in the \"<command-name>\" section. If the command does not exist, the help command will just return an error to the")
	case "mkls":
		fmt.Println("mkls: Creates up and opens a new list file using the format: \"mkls\" or \"mkls <name>\". The command \"mkls\" creates a new list with the name of \"New_List_(10 random alphanumerical digits)\". The command \"mkls <name>\" creates a list where \"<name>\" becomes the name of the list. This means that \"mkls\" can be run without using arguments which result in a list with a random name or if the name argument is given, a named list can be created.")
	case "opls":
		fmt.Println("opls: Opens a designated list as the active list and prepares said list for addition, deletition, or editing of tasks. Uses the format of \"opls <list-name>\" where \"<list-name>\" is the name of the list that the user wishes to be opened. This command loads the file for the list stored in the list direcotry into the program to be edited.")
	case "vwls":
		fmt.Println("vwls: Prints out all elements of the active list file in a readable format. The print starts with the name of the list, then iterates through all of the tasks of the list, printing each attribute of the task including name, severity, and completion status. The \"vwls\" never takes any arguments, but will not error if any are provided, they are just ignored. The \"vwls\" command is occasionally called inside other commands which may print out the same information as \"vwls\" would.")
	case "mktk":
		fmt.Println("mktk: Promps the ")
	case "edls":
		fmt.Println("No in-depth description is currently available.")
	case "dells":
		fmt.Println("No in-depth description is currently available.")
	case "deltk":
		fmt.Println("No in-depth description is currently available.")
	case "lsls":
		fmt.Println("No in-depth description is currently available.")
	case "cdls":
		fmt.Println("No in-depth description is currently available.")
	case "edconf":
		fmt.Println("edconf: Uses the format: \"edconf <attibute> <new-value>\", where \"<attribute>\" is the attribute of the config file that you wish to change and \"<new-value>\" is the new value that you are giving the the attribute you previously chose. The other format that the \"edconf\" command uses is: \"edconf\" which will then prompt you with a series of inputs.")
	case "end":
		fmt.Println("end: Immediately ends the program, ending any running command (a.k.a. goroutines)")
	default:
		fmt.Println("help: Lists out all commands with descriptions")
		fmt.Println("mkls: Creates a .json file in the list directory defined with in the config file.")
		fmt.Println("opls: Opens the designated list and prepares it for the addition, deletion, or editing of tasks.")
		fmt.Println("vwls: Prints the list out in an understandable format to the console.")
		fmt.Println("mktk: Creates a task in the currently opened list. If no list is opened, it returns an error.")
		fmt.Println("edls: Allows changes to be made to the active / opened list.")
		fmt.Println("dells: No description is currently available.")
		fmt.Println("deltk: No description is currently available.")
		fmt.Println("lsls: Lists out all list files in the list directory.")
		fmt.Println("cdls: Either prints the list directory or changes the list directory based off user input.")
		fmt.Println("edconf: Allows the user to edit the config file based of their inputs.")
		fmt.Println("end: Immediately ends the program.")
	}

}

func mkls(args []string) {
	var listName string
	if len(args) == 1 {
		listName = ""
	} else {
		listName = args[1]
	}
	if listName == "" {
		listName = "New_List_" + randHexDec(10)
	}

	file, err := os.Create(listName + ".json")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()
	fileinfo, err := file.Stat()
	if err != nil {
		log.Fatal(err)
		return
	}
	buffer := make([]byte, fileinfo.Size())
	_, err = file.Read(buffer)
	if err != nil {
		fmt.Println("Error reading list file: ", err)
		return
	}

	var data TaskList
	data.Name = listName
	data.Tasks = make([]Task, 0)
	activeList = data
	activeFileName = file.Name()
	saveActiveListToJSON()
}

func opls(args []string) {
	var filename string
	if len(args) == 1 {
		filename = ""
	} else {
		filename = args[1]
	}
	if filename == "" {
		fmt.Println("Filename required: Correct formatting of the \"opls\" command is: \"opls <file-name (no file extension))>\"")
		return
	}
	filename = filename + ".json"
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
		return
	}
	fileinfo, err := file.Stat()
	if err != nil {
		log.Fatal(err)
		return
	}
	buffer := make([]byte, fileinfo.Size())
	_, err = file.Read(buffer)
	if err != nil {
		fmt.Println("Error reading config file: ", err)
		return
	}

	var data TaskList
	err = json.Unmarshal(buffer, &data)
	if err != nil {
		fmt.Println("Error unmarshalling json config data: ", err)
		return
	}

	buffer, err = json.Marshal(data)
	if err != nil {
		log.Fatal(err)
		return
	}

	os.WriteFile(filename, buffer, 0760)

	activeFileName = file.Name()
	activeList = data
	file.Close()
	vwls()
}

func vwls() {
	fmt.Println(activeList.Name, ":")
	for i, n := range activeList.Tasks {
		fmt.Println("  Task", i+1, ":")
		fmt.Println("    Name:", n.Name)
		fmt.Println("    Importance (1-5):", n.Severity)
		fmt.Println("    Is completed (true or false):", n.CompletionStatus)
		if i+1 < len(activeList.Tasks) {
			fmt.Println()
		}
	}
}

func mktk(args []string) {
	if activeFileName == "" {
		fmt.Println("No list file open: A Task Manager list file must be opened for this command to function properly")
		return
	}
	if len(args) != 3 {
		fmt.Println("Incorrect Format: The command \"mktk\" uses the format: \"mktk <new-task-name> <new-task-severity>\". Where \"<new-task-name>\" is replaced witht the name of the task that you are trying to create and \"<new-task-severity>\" is replaced witht the severity/importance of the task that you are trying to create.")
		return
	}
	name := args[1]
	severity, err := strconv.ParseInt(args[2], 10, 64)
	if err != nil {
		log.Fatal(err)
		return
	}

	if name == "" || severity == 0 {
		fmt.Println("Name and Severity (1-5) variables expected: Correct format for \"mkts\" command is: \"mkts <name> <severity>\". Name can be any string as long as it is alphanumerical. Severity can be any number between 1 and 5.")
		return
	}
	if activeFileName == "" {
		fmt.Println("A list file needs to be opened before a task can be created. You can open a list file by using the \"opls\" command. If you need to learn how to use this command you can use the \"help\" by typing in \"help opls\" into the command input.")
		return
	}
	activeList.Tasks = append(activeList.Tasks, Task{name, int(severity), false})
	saveActiveListToJSON()
}

func edls(args []string) {
	if activeFileName == "" {
		fmt.Println("No list file open: A Task Manager list file must be opened for this command to function properly")
		return
	}
	if len(args) < 2 || len(args) > 3 {
		fmt.Println("Incorrect Format: The command \"edls\" uses the format: \"edls <name or tasks> <new-name (if previous input was name)\". Where \"<name or tasks>\" gets replaced with whatever attribute you are trying to edit, either the name of a list or the tasks of a list. If you chose to edit the name of the list, you will need to replace \"<new-name>\" with the new name that you are giving to the list. If you are editing the taskks of a list, you will be prompted with a series of inputs to edit whatever tasks and their attributes you want to edit")
		return
	}
	if args[1] == "name" {
		editListName(args[2])
	} else if args[1] == "tasks" {
		editListTasks()
	}
}

func editListName(newName string) {
	activeList.Name = newName
	saveActiveListToJSON()
}

func editListTasks() {
	for {
		fmt.Print("Name of task you want to edit (\"none\" if you are done editing): ")
		scanner.Scan()
		taskName := scanner.Text()
		if taskName == "none" {
			break
		}
		for i := 0; i < len(activeList.Tasks); i++ {
			if activeList.Tasks[i].Name == taskName {
				breakBool := false
				for !breakBool {
					fmt.Print("Attribute to edit (\"none\" if you are done editing): ")
					scanner.Scan()
					attributeEdited := scanner.Text()
					switch attributeEdited {
					case "name":
						fmt.Print("New task name: ")
						scanner.Scan()
						activeList.Tasks[i].Name = scanner.Text()
						saveActiveListToJSON()
					case "severity":
						fmt.Print("New task severity: ")
						scanner.Scan()
						severity, err := strconv.ParseInt(scanner.Text(), 10, 64)
						if err != nil {
							log.Fatal()
							return
						}
						activeList.Tasks[i].Severity = int(severity)
						saveActiveListToJSON()
					case "completionStatus":
						fmt.Print("New task completion status: ")
						scanner.Scan()
						completionStatus, err := strconv.ParseBool(scanner.Text())
						if err != nil {
							log.Fatal()
							return
						}
						activeList.Tasks[i].CompletionStatus = completionStatus
						saveActiveListToJSON()
					case "none":
						breakBool = true
					default:
						fmt.Println("That is not an attribute of a Task, the attribute that are available to edit for a Task are: name, severity, or completionStatus (spelled and capitalized in that manner)(you can type none in as an attribute if you are done editing the attributes of this task).")
					}
				}
			}
		}
	}
}

func dells(args ...string) {
	if activeFileName == "" {
		fmt.Println("No list file open: A Task Manager list file must be opened for this command to function properly")
		return
	} else if len(args) > 1 {
		fmt.Println("Incorrect Format: use \"help dells\" command to check format for \"deltk\" command.")
		return
	} else if len(args) == 0 {
		fmt.Print("List to be deleted: ")
		scanner.Scan()
		args = append(args, scanner.Text())
	}

	currentFileName := activeFileName
	currentList := activeList
	saveActiveListToJSON()

	opls([]string{"opls", args[0]})
	opls([]string{"opls", strings.Trim(currentFileName, ".json")})

	if currentList.Name == activeList.Name {
		err := os.Remove(args[0] + ".json")
		if err != nil {
			log.Fatal(err)
			return
		}
	} else {
		log.Fatal("File validation failed somewhere. You might have tryed to delete a file that wasn't a list file for the program.")
		os.Exit(1)
	}

}

func deltk(args ...string) {
	if activeFileName == "" {
		fmt.Println("No list file open: A Task Manager list file must be opened for this command to function properly")
		return
	} else if len(args) > 1 {
		fmt.Println("Incorrect Format: use \"help deltk\" command to check format for \"deltk\" command.")
		return
	} else if len(args) == 0 {
		fmt.Print("Task to be deleted: ")
		scanner.Scan()
		args = append(args, scanner.Text())
	}

	taskName := args[0]
	nameExists := false

	var nameIndex int
	for i, n := range activeList.Tasks {
		if n.Name == taskName {
			nameExists = true
			nameIndex = i
			break
		}
	}

	if !nameExists {
		fmt.Println("A task with that name does not exist in this list file. Try using \"vwls\" to see what tasks you can delete or open a different list using \"opls <list-name>\".")
	}

	activeList.Tasks = append(activeList.Tasks[:nameIndex], activeList.Tasks[nameIndex+1:]...)
	saveActiveListToJSON()
}

func lsls() {
	files, err := os.ReadDir(configData.ListDirectory)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Directory of ", configData.ListDirectory)
	fmt.Println()
	for _, file := range files {
		if file.IsDir() {
			fmt.Println(file.Name(), "\t<DIR>")
		} else {
			fmt.Println(file.Name())
		}

	}
}

func cdls(args []string) {
	if len(args) == 1 {
		fmt.Println(configData.ListDirectory)
	} else if len(args) == 2 {
		fmt.Println("Old List Directory: ", configData.ListDirectory)
		configData.ListDirectory = args[1]
		fmt.Println("New List Directory: ", configData.ListDirectory)
	} else {
		fmt.Println("Incorrect Format: \"cdls\" command has the following format: \"cdls\" <optional-new-directory-path>\". Where \"<optional-new-directory-path>\" is the directory where lists get stored and accessed. The path must be an absolute filename (starts from the root of the drive ex: \"C:/Projects/VSCode/TaskManager/Lists/)\"")
	}
}

func edconf(args []string) {
	if len(args) == 3 {
		switch args[1] {
		case "ListDirectory", "listdirectory", "listdir", "lsdir":
			configData.ListDirectory = args[2]
		case "DeleteCompletedTasks", "deletecompletedtasks", "delcomptk":
			value, err := strconv.ParseBool(args[2])
			if err != nil {
				log.Fatal(err)
				return
			}
			configData.DeleteCompletedTasks = value
		case "SortAttribute", "sortattribute", "statt":
			configData.SortAttribute = args[2]
		}
		saveConfigtoJson()
	} else if len(args) == 1 {
		breakBool := false
		for !breakBool {
			fmt.Println("Editable attributes of the config file include:\n1. ListDirectory(through \"cdls\" command)\n2. DeleteCompletedTasks\n3. SortAttribute\n4. none (to quit out of command)")
			fmt.Print("Attribute to edit: ")
			scanner.Scan()
			attribute := scanner.Text()
			switch attribute {
			case "ListDirectory", "List Directory", "listdirectory", "list directory", "listdir", "lsdir":
				fmt.Print("New value for ListDirectory: ")
				scanner.Scan()
				configData.ListDirectory = scanner.Text()
				break
			case "DeleteCompletedTasks", "Delete Completed Tasks", "deletecompletedtasks", "delete completed tasks", "delcomptk":
				fmt.Print("New value for DeleteCompletedTasks: ")
				scanner.Scan()
				value, err := strconv.ParseBool(scanner.Text())
				if err != nil {
					log.Fatal(err)
					return
				}
				configData.DeleteCompletedTasks = value
				break
			case "SortAttribute", "Sort Attribute", "sortattribute", "sort attribute", "statt":
				fmt.Print("New value for SortAttribute: ")
				scanner.Scan()
				configData.SortAttribute = scanner.Text()
				break
			case "none", "done", "no":
				breakBool = true
				break
			default:
				fmt.Println("That is not a valid attribute of config file.")
			}
			saveConfigtoJson()
		}
	} else {
		fmt.Println("Incorrect Format: Use the command \"help edconf\" to learn more.")
		return
	}

}

func runCommand() {
	fmt.Print("\nTask-Manager|-> $ ")
	scanner.Scan()
	word := strings.Split(scanner.Text(), " ")

	switch word[0] {
	case "help":
		help(word)
	case "mkls":
		mkls(word)
	case "opls":
		opls(word)
	case "vwls":
		if activeFileName == "" {
			fmt.Println("No list file open: A Task Manager list file must be opened for this command to function properly")
			return
		}
		vwls()
	case "mktk":
		mktk(word)
	case "edls":
		edls(word)
	case "dells":
		dells(word[1:]...)
	case "deltk":
		deltk(word[1:]...)
	case "lsls":
		lsls()
	case "cdls":
		cdls(word)
	case "edconf":
		edconf(word)
	case "end":
		os.Exit(0)
	default:
		fmt.Println("That is not a valid command. To see a list of all valid commands and their actions, use the \"help\" command.")
	}
}

func randHexDec(size int) string {
	var letterBytes = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, size)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func saveActiveListToJSON() {
	if configData.DeleteCompletedTasks {
		completedTaskIndex := make([]int, 0)
		for i, n := range activeList.Tasks {
			if n.CompletionStatus {
				completedTaskIndex = append(completedTaskIndex, i)
			}
		}
		for _, n := range completedTaskIndex {
			activeList.Tasks = append(activeList.Tasks[:n], activeList.Tasks[n+1:]...)
		}
	}

	if configData.FileNameSameAsListName {
		if strings.Trim(activeFileName, ".json") != activeList.Name {
			fmt.Println("list name does not equal file name")
		}
	}

	activeFile, err := os.OpenFile(activeFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer activeFile.Close()

	jsonData, err := json.Marshal(activeList)
	if err != nil {
		log.Fatal(err)
		return
	}

	_, err = activeFile.Write(jsonData)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func saveConfigtoJson() {
	configFile, err := os.OpenFile("../task_manager.config.json", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer configFile.Close()

	jsonData, err := json.Marshal(configData)
	if err != nil {
		log.Fatal(err)
		return
	}

	_, err = configFile.Write(jsonData)
	if err != nil {
		log.Fatal(err)
		return
	}
}
