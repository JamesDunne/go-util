package web

import (
	"html/template"
	"log"
	"path"
)

import "github.com/JamesDunne/go-util/fs/notify"
import "github.com/JamesDunne/go-util/base"

// Watches the html/*.html templates for changes:
func WatchTemplates(name, templatePath, glob string, preParse func(*template.Template) *template.Template, uiTmpl **template.Template) (watcher *notify.Watcher, deferClean func(), err error) {
	if preParse == nil {
		preParse = func(t *template.Template) *template.Template { return t }
	}

	// Parse template files:
	ui, err := preParse(template.New(name)).ParseGlob(path.Join(base.CanonicalPath(templatePath), glob))
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
				ui, err := preParse(template.New(name)).ParseGlob(path.Join(base.CanonicalPath(templatePath), glob))
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
