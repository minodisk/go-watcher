package watcher

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/fsnotify.v1"
	"gopkg.in/minodisk/go-walker.v1"
)

type Watcher struct {
	Events    chan fsnotify.Event
	Errors    chan error
	fswatcher *fsnotify.Watcher
	done      chan bool
}

func New() *Watcher {
	w := Watcher{
		Events: make(chan fsnotify.Event),
		Errors: make(chan error),
		done:   make(chan bool),
	}
	return &w
}

func (w *Watcher) Watch(pathes []string) (err error) {
	w.fswatcher, err = fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	targetsMap := make(map[string]bool)
	dirsMap := make(map[string]bool)
	filesMap := make(map[string]bool)
	for _, path := range pathes {
		if path == "" {
			continue
		}
		path, err = filepath.Abs(path)
		if err != nil {
			return err
		}
		fi, err := os.Stat(path)
		if err != nil {
			continue
		}
		if fi.IsDir() {
			dirs := walker.FindDirs(path)
			for _, dir := range dirs {
				targetsMap[dir] = true
				dirsMap[dir] = true
			}
		} else {
			targetsMap[filepath.Dir(path)] = true
			filesMap[path] = true
		}
	}

	// var targets []string
	var dirs []string
	var files []string
	for target, _ := range targetsMap {
		// targets = append(targets, target)
		w.fswatcher.Add(target)
	}
	for dir, _ := range dirsMap {
		dirs = append(dirs, dir)
	}
	for file, _ := range filesMap {
		files = append(files, file)
	}

	go func() {
		for {
			select {
			case event := <-w.fswatcher.Events:
				if watch(dirs, files, event.Name) {
					w.Events <- event
				}
			case err := <-w.fswatcher.Errors:
				w.Errors <- err
			case <-w.done:
				return
			}
		}
	}()

	return nil
}

func (w *Watcher) Close() error {
	close(w.done)
	return w.fswatcher.Close()
}

func watch(dirs, files []string, filename string) bool {
	for _, dir := range dirs {
		if in(dir, filename) {
			return true
		}
	}
	for _, file := range files {
		if is(file, filename) {
			return true
		}
	}
	return false
}

func in(dir, filename string) bool {
	return strings.Index(filename, dir) == 0
}

func is(file, filename string) bool {
	return file == filename
}
