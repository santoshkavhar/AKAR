package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var goExecutable = "my_program"

func compileAndRun() error {
	fmt.Println("Compiling...")
	cmd := exec.Command("go", "build", "-o", goExecutable)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Compilation failed:", string(out))
		return err
	}

	fmt.Println("Running...")
	cmd = exec.Command("./" + goExecutable)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func isGoFile(filename string) bool {
	return strings.HasSuffix(filename, ".go")
}

func monitorChanges() {
	var lastCompileTime time.Time

	for {
		// Sleep for a short duration to avoid too frequent checks
		time.Sleep(500 * time.Millisecond)

		filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}

			if !info.IsDir() && isGoFile(info.Name()) {
				modTime := info.ModTime()
				if modTime.After(lastCompileTime) {
					fmt.Println("File change detected. Killing previous instance...")
					exec.Command("pkill", goExecutable).Run()
					lastCompileTime = modTime

					if err := compileAndRun(); err != nil {
						fmt.Println("Error:", err)
					}
				}
			}

			return nil
		})
	}
}

func main() {
	// Initial compile and run
//	if err := compileAndRun(); err != nil {
//		panic(err)
//	}

	monitorChanges()
}

