package rated

type Watcher struct {}

func NewWatcher() *Watcher {
	return &Watcher{}
}

func (w *Watcher) Watch() error {
	return nil
}
