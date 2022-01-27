package watcher

import (
	"log"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Watcher watcher
type Watcher struct {
	locker   sync.RWMutex
	logger   Logger
	notifier *fsnotify.Watcher
	e        map[string]*Event
	p        map[string]struct{}
	o        func(*Watcher, *Event)
	d        func(*Event, *Event) bool
}

// NewWatcher new watcher
func NewWatcher(paths []string, logger Logger) (*Watcher, error) {
	var (
		notifier *fsnotify.Watcher
		watcher  *Watcher
		err      error
	)

	if notifier, err = fsnotify.NewWatcher(); err == nil {
		if logger == nil {
			logger = DefaultLogger
		}

		watcher = &Watcher{
			logger:   logger,
			notifier: notifier,
			e:        make(map[string]*Event),
			p:        make(map[string]struct{}),
		}

		for _, p := range paths {
			watcher.Add(p)
		}
	}
	return watcher, err
}

// Prepare prepare watcher
func (w *Watcher) Prepare(paths []string) *Watcher {
	for _, p := range paths {
		w.Add(p)
		log.Println("watch path", p)
	}
	return w
}

// Run run watcher
func (w *Watcher) Run() *Watcher {
	// run
	go func() {
		var (
			kind   Kind
			action Action
			e      fsnotify.Event
			ok     bool
			err    error
		)

		for {
			select {
			case e, ok = <-w.notifier.Events:
				if !ok {
					break
				}
				// operation
				switch e.Op {
				case fsnotify.Create:
					action = CREATE
				case fsnotify.Remove:
					action = REMOVE
				case fsnotify.Write:
					action = UPDATE
				case fsnotify.Rename:
					// action = Remove
					// for _, p := range w.List() {
					// 	if strings.HasPrefix(e.Name, p) {
					// 		action = UPDATE
					// 		break
					// 	}
					// }
					action = UPDATE
				case fsnotify.Chmod:
					action = UPDATE
				default:
					action = NONE
					continue
				}
				// kind
				switch {
				case !Exist(e.Name):
					kind = UNKNOWN
				case IsDir(e.Name):
					kind = FOLDER
				default:
					kind = FILE
				}

				if w.o != nil {
					var (
						event = &Event{
							Date:   time.Now().Unix(),
							Name:   e.Name,
							Kind:   kind,
							Action: action,
						}
					)
					if w.trigger(event) {
						w.o(w, event)
					}
				}
			case err = <-w.notifier.Errors:
				if err != nil {
					log.Println("fsnotify", err)
				}
			}
		}
	}()
	return w
}

// Exit exit watcher
func (w *Watcher) Exit() {
	w.locker.Lock()
	defer w.locker.Unlock()

	w.notifier.Close()
}

// Operate reply to operation
func (w *Watcher) Operate(o func(*Watcher, *Event)) *Watcher {
	w.o = o
	return w
}

// Diff diff event
func (w *Watcher) Diff(d func(*Event, *Event) bool) *Watcher {
	w.d = d
	return w
}

// Add add watched path
func (w *Watcher) Add(path string) *Watcher {
	w.locker.Lock()
	defer w.locker.Unlock()
	// update dict
	w.p[path] = struct{}{}
	// update watcher
	w.notifier.Add(path)
	return w
}

// Remove remove watched path
func (w *Watcher) Remove(path string) *Watcher {
	w.locker.Lock()
	defer w.locker.Unlock()
	// update dict
	delete(w.p, path)
	// update watcher
	w.notifier.Remove(path)
	return w
}

// List list watched path
func (w *Watcher) List() []string {
	w.locker.RLock()
	defer w.locker.RUnlock()

	var (
		list = make([]string, 0)
	)

	for p := range w.p {
		list = append(list, p)
	}
	return list
}

func (w *Watcher) trigger(e *Event) bool {
	w.locker.RLock()
	defer w.locker.RUnlock()

	var (
		act bool
	)

	if exist, ex := w.e[e.Name]; ex {
		if w.d != nil {
			act = w.d(e, exist)
		} else {
			act = e.Date > exist.Date ||
				e.Date == exist.Date && e.Action != exist.Action
		}
	} else {
		act = true
	}
	// update
	w.e[e.Name] = e
	return act
}
