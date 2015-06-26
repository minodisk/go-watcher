package watcher_test

import (
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"

	"gopkg.in/fsnotify.v1"

	"github.com/minodisk/go-watcher"
)

type Task struct {
	Text     string
	Filename string
	Duration time.Duration
}

func TestWatch(t *testing.T) {
	tasks := []*Task{
		&Task{"foo", "fixtures/foo", time.Second * 1},
		&Task{"baz", "fixtures/bar/baz", time.Second * 2},
		&Task{"quux", "fixtures/bar/qux/quux", time.Second * 3},
	}

	for _, task := range tasks {
		filename, err := filepath.Abs(task.Filename)
		if err != nil {
			t.Fatal(err)
		}
		task.Filename = filename

		go func(task *Task) {
			time.Sleep(task.Duration)
			err := ioutil.WriteFile(task.Filename, []byte(task.Text), 0644)
			if err != nil {
				t.Fatal(err)
			}
		}(task)
	}

	done := make(chan bool)
	w := watcher.NewWatcher()
	go func() {
		i := 0
		for {
			select {
			case event := <-w.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					task := tasks[i]
					i++

					actual := event.Name
					expected := task.Filename
					if actual == expected {
						if i == 3 {
							w.Close()
							done <- true
							return
						}
					} else {
						t.Errorf("filename is expected %s, but actual %s", expected, actual)
					}
				}
			}
		}
	}()
	w.Watch([]string{"fixtures"})
	<-done
}
