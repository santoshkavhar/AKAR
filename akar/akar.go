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

// func main() {

// 	time.Sleep(1 * time.Second)
// 	fmt.Println(os.Getpid(), "Hello World")
//     fmt.Println(int(C.getpid()))
//    //      char *args[]={"./v1",NULL};
//         C.newProc();
// }


// var goExecutable = "my_program"
var AKAR = 0

func compileAndRun() error {
	// fmt.Println("Compiling...")
	// cmd := exec.Command("go", "build")
	// out, err := cmd.CombinedOutput()
	// if err != nil {
	// 	fmt.Println("Compilation failed:", string(out))
	// 	return err
	// }

	fmt.Println("Running...")
	// cmd = exec.Command("./v1")
	C.newProc();
	// cmd.Stdin = os.Stdin
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	return nil
}

func isGoFile(filename string) bool {
	return strings.HasSuffix(filename, ".go")
}

func MonitorChanges() {
	var lastCompileTime time.Time

	for {
		// Sleep for a short duration to avoid too frequent checks
		time.Sleep(5 * time.Second)

		if AKAR == 1{
			fmt.Println("Compiling...")
		}
		cmd := exec.Command("go", "build")
		out, err := cmd.CombinedOutput()
		if err != nil  && AKAR == 1{
			fmt.Println("Compilation failed:", string(out))
		}
	
		if AKAR == 1{
			fmt.Println("Running...")
		}
		// cmd = exec.Command("./v1")
		C.newProc();

		if AKAR != 1 {
			continue
		}

		filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}

			if !info.IsDir() && isGoFile(info.Name()) {
				modTime := info.ModTime()
				if modTime.After(lastCompileTime) {
					fmt.Println("File change detected.")// Killing previous instance...")
					//exec.Command("pkill", "v1").Run()
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


