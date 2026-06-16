// Package hotconf provides a file watcher that hot-reloads configuration files without restarting the process.
package hotconf

import (
	"context"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	watcher   *fsnotify.Watcher
	debouncer *Debouncer
	// TOOD: make callback list
	fn func(path string)
}

func NewWatcher(debounceDuration time.Duration) (*Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return &Watcher{
		watcher:   watcher,
		debouncer: NewDebouncer(debounceDuration),
	}, nil
}

func (w *Watcher) Watch(path string, fn func(path string)) error {
	if _, err := os.Stat(path); err != nil {
		return err
	}

	if err := w.watcher.Add(path); err != nil {
		return err
	}

	w.fn = fn

	return nil
}

func (w *Watcher) Start(ctx context.Context) error {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-w.watcher.Events:
				if event.Has(fsnotify.Rename) || event.Has(fsnotify.Remove) {
					go func() {
						time.Sleep(10 * time.Millisecond)
						for _, p := range w.watcher.WatchList() {
							_ = w.watcher.Add(p)
						}
					}()
				}
				if event.Has(fsnotify.Create) || event.Has(fsnotify.Write) {
					w.debouncer.Do(func() {
						w.fn(event.Name)
					})
				}
			case <-w.watcher.Errors:
				continue
			}
		}
	}()
	return nil
}

func (w *Watcher) Stop() {
	_ = w.watcher.Close()
}
