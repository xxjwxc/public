package kv

var kvConf = NewConfig()

// Init ...
func Init(opts ...Option) error {
	for _, o := range opts {
		o(kvConf)
	}

	return kvConf.Init()
}

// SetOptions ...
func SetOptions(opts ...Option) error {
	for _, o := range opts {
		o(kvConf)
	}

	return kvConf.reset()
}

// Put ...
func Put(path string, value interface{}) error {
	return kvConf.Put(path, value)
}

// Delete ...
func Delete(path string) error {
	return kvConf.Delete(path)
}

// Get ...
func Get(keys ...string) *Result {
	return kvConf.Get(keys...)
}

// Watch ...
func Watch(path string, handler func(*Result)) error {
	return kvConf.Watch(path, handler)
}

// StopWatch ...
func StopWatch(path ...string) {
	kvConf.StopWatch(path...)
}
