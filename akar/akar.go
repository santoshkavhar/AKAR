package akar

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// #include <stdio.h>
// #include <errno.h>
// #include <unistd.h>
// void newProc(){
//	char *args[]={"./v1", NULL};
//	execvp(args[0], args);
//}
import "C"

// TODO: Get file name from the directory or check the mod file or something else
// var goExecutable = "my_program"
var AKAR = 0

func compileAndRun() error {

	if AKAR == 1{
		fmt.Println("Compiling...")
	}

	cmd := exec.Command("go", "build")
	out, err := cmd.CombinedOutput()
	if err != nil && AKAR == 1 {
		fmt.Println("Compilation failed:", string(out))
		return err
	}
	if AKAR == 1{
		fmt.Println("Running...")
	}
	C.newProc();

	return nil
}

func isGoFile(filename string) bool {
	return strings.HasSuffix(filename, ".go")
}

func MonitorChanges() {
	var lastCompileTime time.Time

	for {
		// TODO: Use signals to monitor file changes
		// Sleep for a short duration to avoid too frequent checks
		time.Sleep(5 * time.Second)

		// TOOD: Find effective way to avoid compilation and Running everytime
		compileAndRun() 

		if AKAR != 1 {
			continue
		}

		// TODO: Cleanup this code
		filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
			if err != nil {
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


