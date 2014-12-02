package oo

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"gopkg.in/fsnotify.v1"
)

type OO struct {
	mode     fsnotify.Op
	show     bool
	excludes []*regexp.Regexp
	path     Path
	tmpl     *template.Template
	watched  map[string]bool
}

func New(mode, dir, cmd string, show bool, excludes []string) (*OO, error) {

	path := Path(dir)

	tmpl, err := template.New("oo").Parse(cmd)
	if err != nil {
		return nil, err
	}

	regExcludes := make([]*regexp.Regexp, 0, len(excludes))
	for _, exclude := range excludes {
		regExclude, err := regexp.Compile(exclude)
		if err != nil {
			return nil, err
		}
		regExcludes = append(regExcludes, regExclude)
	}

	oo := &OO{
		path:     path,
		tmpl:     tmpl,
		show:     show,
		excludes: regExcludes,
		watched:  make(map[string]bool),
	}

	for _, ch := range mode {
		switch ch {
		case 'w':
			oo.mode |= fsnotify.Write
		case 'd':
			oo.mode |= fsnotify.Remove
		case 'r':
			oo.mode |= fsnotify.Rename
		case 'n':
			oo.mode |= fsnotify.Create
		case 'c':
			oo.mode |= fsnotify.Chmod
		}
	}

	return oo, nil
}

func (oo *OO) Watch() error {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	if err := filepath.Walk(oo.path.String(), oo.walkFunc(watcher)); err != nil {
		return err
	}

loop:
	for {
		select {
		case event := <-watcher.Events:

			file, err := filepath.Rel(oo.path.String(), event.Name)
			if err != nil {
				log.Println(err.Error())
				continue
			}

			for _, exclude := range oo.excludes {
				if exclude.MatchString(file) {
					log.Println("exclude", file)
					continue loop
				}
			}

			if event.Op&fsnotify.Create == fsnotify.Create {
				if stat, err := os.Stat(event.Name); err != nil {
					return err
				} else if stat.IsDir() {
					if err := filepath.Walk(event.Name, oo.walkFunc(watcher)); err != nil {
						log.Println(err.Error())
						continue
					}
				}
			}

			if event.Op&fsnotify.Remove == fsnotify.Remove {
				if _, ok := oo.watched[event.Name]; ok {
					watcher.Remove(event.Name)
					delete(oo.watched, event.Name)
				}
			}
			if oo.mode&event.Op == fsnotify.Op(0) {
				continue
			}

			pre := oo.path.cd()
			var buf bytes.Buffer
			if err := oo.tmpl.Execute(&buf, Path(file)); err != nil {
				log.Println(err.Error())
				continue
			}
			pre.cd()

			if oo.show {
				log.Println(buf.String())
			}
			args := strings.Split(buf.String(), " ")
			cmd := exec.Command(args[0], args[1:]...)
			cmd.Stdout = os.Stdout
			cmd.Stdin = os.Stdin
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				log.Println(err)
			}
		case err := <-watcher.Errors:
			return err
		}
	}

	return nil
}

func (oo *OO) walkFunc(watcher *fsnotify.Watcher) filepath.WalkFunc {
	return filepath.WalkFunc(func(path string, info os.FileInfo, err error) error {
		if _, ok := oo.watched[path]; !ok && info.IsDir() {
			if err := watcher.Add(path); err != nil {
				return err
			}
			log.Println("watch", path)
			oo.watched[path] = true
		}

		return nil
	})
}
