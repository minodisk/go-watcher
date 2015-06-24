package watcher

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/fsnotify.v1"
	"gopkg.in/minodisk/go-walker.v1"
)

func Watch(pathes []string, changed chan string, finish chan bool, finished chan bool) (err error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	dirBoolMap := make(map[string]bool)
	for i, path := range pathes {
		path, err = filepath.Abs(path)
		if err != nil {
			return err
		}
		pathes[i] = path
		fi, err := os.Stat(path)
		if err != nil {
			return err
		}
		if fi.IsDir() {
			dirs := walker.FindDirs(path)
			for _, dir := range dirs {
				dirBoolMap[dir] = true
			}
		} else {
			dirBoolMap[filepath.Dir(path)] = true
		}
	}

	var dirs []string
	for dir, _ := range dirBoolMap {
		log.Printf("[watcher] start watching dir: %s", dir)
		dirs = append(dirs, dir)
		watcher.Add(dir)
	}

	// for {
	// 	select {
	// case event := <-watcher.Events:
	// 	if event.Op&fsnotify.Create == fsnotify.Create || event.Op&fsnotify.Write == fsnotify.Write {
	// 		if watches(dirs, event.Name) {
	// 			log.Printf("[watcher] detect modified: %s", event.Name)
	// 			changed <- event.Name
	// 		}
	// 	}
	// case err := <-watcher.Errors:
	// 	log.Println("[watcher] error:", err)
	// case <-finish:
	// log.Println("[watcher] finishing watching")
	// finished <- true
	// log.Println("==============")
	// err := watcher.Close()
	// if err != nil {
	// 	log.Printf("[watcher] fail to close: %v", err)
	// }
	// 	}
	// }

	log.Println("[watcher] finished watching")

	return nil
}

func watches(dirs []string, filename string) bool {
	for _, dir := range dirs {
		if contains(dir, filename) {
			return true
		}
	}
	return false
}

func contains(dir, filename string) bool {
	return strings.Index(filename, dir) == 0
}
