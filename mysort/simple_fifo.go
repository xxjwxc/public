package mysort

// Fifo 先进先出排序(去重)
type Fifo struct {
	items []interface{}
}

// Push 推送
func (f *Fifo) Push(item interface{}) {
	for _, v := range f.items {
		if v == item {
			return
		}
	}

	f.items = append(f.items, item) // 没有就添加
}

// Gets 获取
func (f *Fifo) Gets() []interface{} {
	return f.items
}
