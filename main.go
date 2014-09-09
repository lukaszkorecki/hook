package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/howeyc/fsnotify"
)

func getAllDirs(root string) []string {
	var dirs []string

	fn := func(path string, info os.FileInfo, err error) error {
		if err == nil {
			if info.IsDir() {
				// log.Printf("Adding > %s", path)
				dirs = append(dirs, path)
			}
		}

		return nil
	}
	filepath.Walk(root, fn)

	return dirs
}

func runCmd(c string, file string, root string, passArgs bool) {
	command := strings.Split(strings.Trim(c, " "), " ")
	cmd := exec.Command(command[0])

	var args []string
	if passArgs {
		args = append(command, file)
	} else {
		args = command
	}

	cmd.Args = args
	cmd.Dir = root

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Start()

	err := cmd.Wait()

	if err != nil {
		log.Printf("ERR\n%v", err)
	}

	fmt.Println("Done")
}

var (
	dir      string
	command  string
	passArgs bool
)

func main() {
	flag.StringVar(&dir, "dir", ".", "Directory to monitor")
	flag.StringVar(&command, "cmd", "echo", "Command to run")
	flag.BoolVar(&passArgs, "args", false, "Pass file path to command?")

	flag.Parse()

	watchDir, err := filepath.Abs(dir)
	if err != nil {
		log.Fatal(err)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)

	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				// log.Printf("> %v %v has changed", ev.Name, ev)
				if ev.IsModify() || ev.IsCreate() {
					runCmd(command, ev.Name, watchDir, passArgs)
				}
			case err := <-watcher.Error:
				log.Println("> error:", err)
			}
		}
	}()

	os.Chdir(watchDir)
	for _, d := range getAllDirs(watchDir) {
		err = watcher.Watch(d)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Printf("> Watching '%v' for changes.", watchDir)

	<-done

	watcher.Close()
}
