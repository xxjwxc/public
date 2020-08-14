package mysort

// 队列名
type base struct {
	items []interface{}
}

func (b *base) Add(item interface{}) {
	b.items = append(b.items, item)
}

// EqualAt 获取相等位置下标,不重复返回-1
func (b *base) EqualAt(item interface{}) int {
	for i, v := range b.items {
		if v == item {
			return i
		}
	}

	return -1
}

// GetItems 获取
func (b *base) GetItems() interface{} {
	return b.items
}

// PushBack 尾部添加数据
func (b *base) PushBack(item interface{}) {
	b.items = append(b.items, item) // 没有就添加
}

// PushFront 头部添加数据
func (b *base) PushFront(item interface{}) {
	b.Insert(item, 0)
}

// Insert 插入元素
func (b *base) Insert(item interface{}, i int) {
	if i < 0 {
		i = 0
	}
	if i > len(b.items) {
		i = len(b.items)
	}

	// 高效插入
	b.items = append(b.items, nil)   // 切片扩展1个空间
	copy(b.items[i+1:], b.items[i:]) // 向后移动1个位置
	b.items[i] = item                // 设置新添加的元素
}

// Replace 替换
func (b *base) Replace(item interface{}, i int) {
	b.items[i] = item
}

// ReplaceFront 替换i位置且将item插入到头部
func (b *base) ReplaceFront(item interface{}, i int) {
	copy(b.items[1:], b.items[0:i]) // 向后整体移动
	b.items[0] = item
}

// ReplaceBack 替换i位置且将item插入到头尾部
func (b *base) ReplaceBack(item interface{}, i int) {
	copy(b.items[i:], b.items[i+1:]) // 向后整体移动
	b.items[0] = item
}
