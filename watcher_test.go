package watcher_test

import (
	"log"
	"testing"

	"github.com/minodisk/go-watcher"
)

func TestWatch(t *testing.T) {
	// go func() {
	// 	time.Sleep(time.Second * 1)
	// 	err := ioutil.WriteFile("fixtures/foo", []byte("foo"), 0644)
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}
	// }()

	finished := make(chan bool)
	changed := make(chan string)
	finish := make(chan bool)

	go func() {
		err := watcher.Watch([]string{"fixtures"}, changed, finish, finished)
		if err != nil {
			t.Fatal(err)
		}

		for {
			// 	select {
			// 	case filename := <-changed:
			// 		log.Printf("file changed: %s", filename)
			// 		target, err := filepath.Abs("fixtures/foo")
			// 		if err != nil {
			// 			t.Fatal(err)
			// 		}
			// 		log.Println(filename)
			// 		log.Println(target)
			// 		if filename == target {
			// 			log.Println("=====")
			// 			finish <- true
			// 			log.Println("-----")
			// 		}
			// 	}
		}
		// finished <- 1
	}()

	<-finished
	log.Println("done!!")
}
