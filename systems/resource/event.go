package resource

// EventLoadStart is trigger in gorge when a resource starts loading.
type EventLoadStart struct {
	Name     string
	Resource interface{}
}

// EventLoadComplete is triggered when a resource finished loading.
type EventLoadComplete struct {
	Name     string
	Resource interface{}
	Err      error
}

// EventOpen is triggered when a resource is opened which differs from LoadStart
// it was meant to create some kind of progress bar while loading resources.
type EventOpen struct {
	Name string
}
