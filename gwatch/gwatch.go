package gwatch

import (
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

type Event = fsnotify.Event

func OnChange(root string, fn func(Event)) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return w.Add(path)
		}
		return nil
	})
	for ev := range w.Events {
		fn(ev)
	}
}
