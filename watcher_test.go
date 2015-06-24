package watcher_test

import (
	"log"
	"path/filepath"
	"testing"

	"github.com/minodisk/go-watcher"
)

func TestWatch(t *testing.T) {
	changed := make(chan string)
	finish := make(chan bool)
	finished := make(chan bool)
	err := watcher.Watch([]string{"fixtures"}, changed, finish, finished)
	if err != nil {
		t.Fatal(err)
	}
	for {
		select {
		case filename := <-changed:
			log.Printf("file changed: %s", filename)
			target, err := filepath.Abs("fixtures/foo")
			if err != nil {
				t.Fatal(err)
			}
			log.Println(filename)
			log.Println(target)
			if filename == target {
				log.Println("=====")
				finish <- true
				log.Println("-----")
			}
		}
	}
	<-finished
	log.Println("done!!")
}
