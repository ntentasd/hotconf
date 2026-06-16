package hotconf

import (
	"os"
	"testing"
	"time"
)

func TestWatcherFiresOnWrite(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test.txt")
	if err != nil {
		t.Fatal(err)
	}

	watcher, err := NewWatcher(time.Second)
	if err != nil {
		t.Fatal(err)
	}

	ch := make(chan struct{}, 1)
	errCh := make(chan error, 1)

	if err = watcher.Watch(tmpFile.Name(), func(path string) {
		ch <- struct{}{}
	}); err != nil {
		t.Fatal(err)
	}

	if err = watcher.Start(t.Context()); err != nil {
		t.Fatal(err)
	}
	defer watcher.Stop()

	go func() {
		_, err := tmpFile.WriteString("test")
		if err != nil {
			errCh <- err
		}
		_ = tmpFile.Sync()
	}()

	select {
	case <-ch:
		// pass
	case err := <-errCh:
		t.Fatal(err)
	case <-time.After(time.Second * 2):
		t.Fatal("timeout")
	}
}
