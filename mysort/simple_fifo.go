package mysort

// Fifo 先进先出排序(去重)
type Fifo struct {
	base
}

// Push 推送
func (f *Fifo) Push(item interface{}) {
	if f.EqualAt(item) >= 0 {
		return
	}

	f.PushBack(item) // 没有就添加
}

// PushGrab 推送（重复插位到尾部）
func (f *Fifo) PushGrab(item interface{}) {
	index := f.EqualAt(item)
	if index >= 0 {
		f.ReplaceBack(item, index)
		return
	}

	f.PushBack(item) // 没有就添加
}

// Gets 获取
func (f *Fifo) Gets() []interface{} {
	return f.items
}
