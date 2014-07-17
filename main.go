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
				dirs = append(dirs, path)
			}
		}

		return nil
	}
	filepath.Walk(root, fn)

	return dirs
}

func runCmd(c string, file string, passArgs bool) {
	command := strings.Split(c, " ")
	cmd := exec.Command(command[0])

	var args []string
	if passArgs {
		args = append(command[1:], file)
	} else {
		args = command[1:]
	}



	cmd.Args = args
	out, err := cmd.Output()
	if err != nil {
		log.Printf("ERR\n%v", err)
	}

	fmt.Println(string(out))
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

	log.Printf("%v %v %v", dir, command, passArgs)

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
				log.Printf("> %v has changed", ev.Name)
				runCmd(command, ev.Name, passArgs)
			case err := <-watcher.Error:
				log.Println("> error:", err)
			}
		}
	}()

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
