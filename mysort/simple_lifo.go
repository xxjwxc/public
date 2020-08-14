package mysort

// Lifo 先进后出排序(去重)
type Lifo struct {
	base
}

// Push 推送(重复抛弃)
func (f *Lifo) Push(item interface{}) {
	if f.EqualAt(item) >= 0 {
		return
	}

	f.PushFront(item) // 没有就添加
}

// PushGrab 推送（去重插位到头部）
func (f *Lifo) PushGrab(item interface{}) {
	index := f.EqualAt(item)
	if index >= 0 {
		f.ReplaceFront(item, index)
		return
	}

	f.PushFront(item) // 没有就添加
}

// Gets 获取
func (f *Lifo) Gets() []interface{} {
	return f.items
}
