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

func NewWatcher() *Watcher {
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
		dirs = append(dirs, dir)
		w.fswatcher.Add(dir)
	}

	go func() {
		for {
			select {
			case event := <-w.fswatcher.Events:
				if watches(dirs, event.Name) {
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
