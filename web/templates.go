package web

import (
	"html/template"
	"log"
	"path"
)

import "github.com/JamesDunne/go-util/fs/notify"

// Watches the html/*.html templates for changes:
func WatchTemplates(name, templatePath, glob string, uiTmpl **template.Template) (watcher *notify.Watcher, deferClean func(), err error) {
	// Parse template files:
	tmplGlob := path.Join(templatePath, glob)
	ui, err := template.New(name).ParseGlob(tmplGlob)
	if err != nil {
		return nil, nil, err
	}
	*uiTmpl = ui

	// Watch template directory for file changes:
	watcher, err = notify.NewWatcher()
	if err != nil {
		return nil, nil, err
	}
	deferClean = func() { watcher.RemoveWatch(templatePath); watcher.Close() }

	// Process watcher events
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				if ev == nil {
					break
				}
				//log.Println("event:", ev)

				// Update templates:
				var err error
				ui, err := template.New(name).ParseGlob(tmplGlob)
				if err != nil {
					log.Println(err)
					break
				}
				*uiTmpl = ui
			case err := <-watcher.Error:
				if err == nil {
					break
				}
				log.Println("watcher error:", err)
			}
		}
	}()

	// Watch template file for changes:
	watcher.Watch(templatePath)

	return
}
