package akar

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"unsafe"
)

// #include <stdio.h>
// #include <errno.h>
// #include <unistd.h>
// #include <stdlib.h>
// void runSameProc(char *str){
//	char *args[]={str, NULL};
//	execvp(args[0], args);
//}
import "C"

// TODO: Get file name from the directory or check the mod file or something else
// var goExecutable = "my_program"
var AKAR = 0
var binaryName = ""
var lastCompileTime = time.Now()

func compileAndRun() error {

	if AKAR == 1 {
		fmt.Println("Compiling...")
	}

	cmd := exec.Command("go", "build")
	out, err := cmd.CombinedOutput()
	if err != nil && AKAR == 1 {
		fmt.Println("Compilation failed:", string(out))
		return err
	}
	//fmt.Println(out)
	if AKAR == 1 {
		fmt.Println("Running...")
	}

	// Convert Go string to C string
	binaryName = "./" + binaryName
	cBinaryName := C.CString(binaryName)
	defer C.free(unsafe.Pointer(cBinaryName)) // Free the C string when done

	// Call the C function with the C string
	C.runSameProc(cBinaryName)

	return nil
}

func isGoFile(filename string) bool {
	return strings.HasSuffix(filename, ".go")
}

func MonitorChanges() {

	getBinaryName()

	for {
		// Sleep for a short duration to avoid too frequent checks
		// TODO: Let this be user settable
		time.Sleep(5 * time.Second)

		if AKAR != 1 {
			continue
		}

		filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Println("error walking")
				return nil
			}

			if !info.IsDir() && isGoFile(info.Name()) {
				modTime := info.ModTime()
				if modTime.After(lastCompileTime) {
					fmt.Println("File change detected.")
					lastCompileTime = modTime

					if err := compileAndRun(); err != nil {
						fmt.Println("Error: ", err)
					}
				}
			}

			return nil
		})
	}
}

func getBinaryName() {

	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	//fmt.Println("Current Working Directory:", cwd)

	filePath := filepath.Join(cwd, "go.mod")

	// Open the go.mod file for reading
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening go.mod:", err)
		return
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Read the first line of the file
	if scanner.Scan() {
		firstLine := scanner.Text()
		// Split the first line by whitespace
		words := strings.Fields(firstLine)

		if len(words) >= 2 {
			// Select the second word
			secondWord := words[1]

			// Split the second word by "/"
			parts := strings.Split(secondWord, "/")

			if len(parts) > 0 {
				// Output the last part
				lastPart := parts[len(parts)-1]
				//fmt.Print("Last word from the second word(Binary name): ", lastPart, "\r")
				binaryName = lastPart
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading go.mod:", err)
	}

}
