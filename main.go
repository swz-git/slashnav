package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"

	_ "embed"

	. "github.com/logrusorgru/aurora"
)

//go:embed toptext.txt
var topText string

const tmpfile = "/tmp/slashnav-tomove.txt"

func ClearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func HideStdin() {
	// disable input buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
}

func GetCurrentDirInfo() (string, []string) {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fileList, err := os.ReadDir(currentDir)
	if err != nil {
		panic(err)
	}
	return currentDir, DirEntries2Strings(fileList)
}

func RemoveIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}

func DirEntries2Strings(entries []os.DirEntry) []string {
	strings := []string{".."}
	for _, entry := range entries {
		if entry.IsDir() {
			strings = append(strings, entry.Name())
		}
	}
	return strings
}

func EnableSigHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			currentDir, _ := GetCurrentDirInfo()
			// write cwd to file
			err := ioutil.WriteFile(tmpfile, []byte(currentDir), 0777)
			if err != nil {
				panic(err)
			}
			fmt.Printf("%s, exiting...\n", sig)
			ClearScreen()
			os.Exit(0)
		}
	}()
}

type TermSize struct {
	Width  int
	Height int
}

func getTermSize() TermSize {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	split := strings.Split(strings.Split(string(out), "\n")[0], " ")
	height, err := strconv.Atoi(split[0])
	if err != nil {
		panic(err)
	}
	width, err := strconv.Atoi(split[1])
	if err != nil {
		panic(err)
	}
	return TermSize{width, height}
}

func Frame(itemsToList []string, activeItem int, currentDir string, searchString string) {
	ClearScreen()
	newlineCount := 0

	fmt.Print(topText)
	newlineCount += strings.Count(topText, "\n")

	fmt.Println("📂", Underline(currentDir))
	newlineCount++

	if searchString != "" {
		fmt.Println("🔎", Underline(searchString))
		newlineCount++
	}

	avaliabeLines := getTermSize().Height - newlineCount - 1

	startPos := int(math.Max(float64(activeItem-avaliabeLines/2), 0))
	startPos = int(math.Min(float64(startPos), float64(len(itemsToList)-avaliabeLines)))

	itemsListed := 0
	for i, item := range itemsToList {
		if i < startPos || i >= startPos+avaliabeLines {
			continue
		}
		itemsListed++
		if activeItem == i {
			fmt.Println(" >", Blue(Bold(item)))
		} else {
			fmt.Println("  ", item)
		}
	}
	if itemsListed == avaliabeLines && len(itemsToList)-avaliabeLines != startPos {
		fmt.Print("↓ ...")
	}
}

func main() {
	if os.Getenv("SLASHNAV_WRAPPER") != "1" {
		panic("Please use the wrapper script called \"slash\"")
	}
	EnableSigHandler()
	HideStdin()
	last3 := [3]int{}
	var fileIndex uint16 = 0
	searchString := ""

	fileListSearched := []string{}

	currentDir, fileList := GetCurrentDirInfo()

	Frame(fileList, int(fileIndex), currentDir, searchString)
	for {
		var b []byte = make([]byte, 1)
		os.Stdin.Read(b)

		currentDir, fileList = GetCurrentDirInfo()

		if searchString == "" {
			fileListSearched = fileList
		}
		last3[0] = last3[1]
		last3[1] = last3[2]
		last3[2] = int(b[0])

		switch last3 {
		case [3]int{27, 91, 65}:
			// up arrow
			if fileIndex != 0 {
				fileIndex--
			}
		case [3]int{27, 91, 66}:
			// down arrow
			if fileIndex != uint16(len(fileList)-1) {
				fileIndex++
			}
		case [3]int{27, 91, 67}:
			// right arrow
		case [3]int{27, 91, 68}:
			// left arrow
		default: // any other key
			switch b[0] {
			case 10: // enter
				os.Chdir(currentDir + "/" + fileListSearched[fileIndex])
				currentDir, fileList = GetCurrentDirInfo()
				fileIndex = 0
				searchString = ""
			case 127: // backspace
				if len(searchString) > 0 {
					searchString = searchString[:len(searchString)-1]
				}
			case 27: // escape
			case 91: // escape
			default:
				if string(b) != "" {
					fmt.Println(b)
					searchString += string(b)
					fileIndex = 0
				}
			}
		}

		fileListSearched = []string{}
		for _, item := range fileList {
			if strings.Contains(strings.ToLower(item), strings.ToLower(searchString)) || item == ".." {
				fileListSearched = append(fileListSearched, item)
			}
		}

		Frame(fileListSearched, int(fileIndex), currentDir, searchString)
	}
}
