package ratelimit

import (
	"sync"
	"time"
)

// singleRule 单组用户访问控制策略
type singleRule struct {
	defaultExpiration            time.Duration               // 表示计时周期,每条访问记录需要保存的时长，超过这个时长的数据记录将会被清除
	numberOfAllowedAccesses      int                         // 在计时周期内最多允许访问的次数
	estimatedNumberOfOnlineUsers int                         // 在计时周期内预计有多少个用户会访问网站，建议选用一个稍大于实际值的值，以减少内存分配次数
	cleanupInterval              time.Duration               // 默认多长时间需要执行一次清除过期数据操作
	visitorRecords               []*autoGrowCircleQueueInt64 // 用于存储用户的每一条访问记录
	usedVisitorRecordsIndex      sync.Map                    // 存储visitorRecords中已使用的数据索引,key代表用户名或IP,为文本或数字类型,value代表visitorRecords中的下标位置
	notUsedVisitorRecordsIndex   map[int]struct{}            // 对应visitorRecords中未使用的数据的下标位置，其自身非并发安全，其并发安全由locker实现,因sync.Map计算长度不优
	lockerForKeyIndex            *sync.RWMutex               // 只用于分配用户KEY，即只需保证用户KEY正确的分配在usedVisitorRecordsIndex与notUsedVisitorRecordsIndex
}

/*
 newSingleRule 初始化一个条单组用户访问控制控制策略,例：
 	vc := newSingleRule(time.Minute*30, 50) 或者 vc := newSingleRule(time.Minute*30, 50, 1000)
	它表示:
	在30分钟内每个用户最多允许访问50次,并且我们预计在这30分钟内大致有1000个用户会访问我们的网站
	1000为可选字段，此参数可默认不填写，主要是用于提升性能，类似于声明切片时的cap,绝大部分情况下无需关注此参数。
*/
func newSingleRule(defaultExpiration time.Duration, numberOfAllowedAccesses int, estimatedNumberOfOnlineUserNum ...int) *singleRule {
	// 规范化numberOfAllowedAccesses
	// 若参数numberOfAllowedAccesses设置是否合理，在此被强行修改为1
	if numberOfAllowedAccesses <= 0 {
		numberOfAllowedAccesses = 1
	}

	// 规范化estimatedNumberOfOnlineUsers
	// estimatedNumberOfOnlineUsers没填写,或者是乱填写的,就默认用numberOfAllowedAccesses
	estimatedNumberOfOnlineUsers := 0
	if len(estimatedNumberOfOnlineUserNum) > 0 {
		estimatedNumberOfOnlineUsers = estimatedNumberOfOnlineUserNum[0]
	}

	if estimatedNumberOfOnlineUsers <= 0 {
		estimatedNumberOfOnlineUsers = numberOfAllowedAccesses
		// 普遍而言，某一段时间内在线用户数达到1000已经较大，所以除非用户指定estimatedNumberOfOnlineUserNum，否则最大值定义为1000
		// 在线用户数是指在某一段时间内访问过的唯一用户总数
		if estimatedNumberOfOnlineUsers > 1000 {
			estimatedNumberOfOnlineUsers = 1000
		}
	}

	// 规范化defaultExpiration
	cleanupInterval := defaultExpiration / 100
	//强行修正清除过期数据的最长时间间隔与最短时间间隔
	if cleanupInterval < time.Second*1 {
		cleanupInterval = time.Second * 1
	}
	if cleanupInterval > time.Second*60 {
		cleanupInterval = time.Second * 60
	}
	vc := createSingleRule(defaultExpiration, cleanupInterval, numberOfAllowedAccesses, estimatedNumberOfOnlineUsers)

	// 定期清除过期数据,并定期清理内存
	go vc.deleteExpired()
	return vc
}

func createSingleRule(defaultExpiration, cleanupInterval time.Duration, numberOfAllowedAccesses, estimatedNumberOfOnlineUsers int) *singleRule {
	var vc singleRule
	vc.defaultExpiration = defaultExpiration
	vc.cleanupInterval = cleanupInterval
	vc.numberOfAllowedAccesses = numberOfAllowedAccesses
	vc.estimatedNumberOfOnlineUsers = estimatedNumberOfOnlineUsers
	vc.notUsedVisitorRecordsIndex = make(map[int]struct{})
	vc.lockerForKeyIndex = new(sync.RWMutex)
	// 根据在线用户数量初始化用户访问记录数据
	vc.visitorRecords = make([]*autoGrowCircleQueueInt64, vc.estimatedNumberOfOnlineUsers)
	for i := range vc.visitorRecords {
		vc.visitorRecords[i] = newAutoGrowCircleQueueInt64(vc.numberOfAllowedAccesses)
		// 刚刚开始时，所有数据都未使用，放入未使用索引中
		vc.notUsedVisitorRecordsIndex[i] = struct{}{}
	}
	return &vc
}

// getIndexFrom 根据用户key返回其数据在visitorRecords中的下标
func (s *singleRule) getIndexFrom(key interface{}) int {
	// 大部分情况下是读，只有少部分情况下是写
	// 只需要用到读锁
	s.lockerForKeyIndex.RLock()
	// 现有访问记录中有，则直接返回
	if index, exist := s.usedVisitorRecordsIndex.Load(key); exist {
		s.lockerForKeyIndex.RUnlock()
		return index.(int)
	}
	s.lockerForKeyIndex.RUnlock()

	// 以下需要用到互斥锁
	s.lockerForKeyIndex.Lock()
	defer s.lockerForKeyIndex.Unlock()
	// visitorRecords有闲置空间，则从闲置空间中获取一条来返回
	if len(s.notUsedVisitorRecordsIndex) > 0 {
		for index := range s.notUsedVisitorRecordsIndex {
			delete(s.notUsedVisitorRecordsIndex, index)
			s.usedVisitorRecordsIndex.Store(key, index)
			s.visitorRecords[index].key = key
			return index
		}
	}

	// visitorRecords没有闲置空间时，则需要插入一条新数据到visitorRecords中
	queue := newAutoGrowCircleQueueInt64(s.numberOfAllowedAccesses)
	queue.key = key
	s.visitorRecords = append(s.visitorRecords, queue)
	index := len(s.visitorRecords) - 1 // 最后一条的位置即为新的索引位置
	s.usedVisitorRecordsIndex.Store(key, index)
	return index
}

// updateIndexOf 经过一段时间无访问数据时，从usedVisitorRecordsIndex中删除用户Key
func (s *singleRule) updateIndexOf(key interface{}) {
	s.lockerForKeyIndex.Lock()
	defer s.lockerForKeyIndex.Unlock()
	if index, exist := s.usedVisitorRecordsIndex.Load(key); exist {
		s.usedVisitorRecordsIndex.Delete(key)                  // 删除完过期数据之后，如果该用户的所有访问记录均过期了，那么就删除该用户
		s.notUsedVisitorRecordsIndex[index.(int)] = struct{}{} // 并把该空间返还给notUsedVisitorRecordsIndex以便下次重复使用
	}
}

// allowVisit 是否允许访问,允许访问则往访问记录中加入一条访问记录
func (s *singleRule) allowVisit(key interface{}) bool {
	return s.add(key) == nil
}

// remainingVisits 剩余访问次数
func (s *singleRule) remainingVisits(key interface{}) int {
	index := s.getIndexFrom(key)
	return s.visitorRecords[index].unUsedSize()
}

// add 增加一条访问记录
func (s *singleRule) add(key interface{}) (err error) {
	index := s.getIndexFrom(key)
	s.visitorRecords[index].deleteExpired(key)
	return s.visitorRecords[index].pushWithConcurrencysafety(s.defaultExpiration)
}

// addFromBackUpFile 增加一条访问记录,从备份文件中增加,从备份文件中过来的数据不可信，有可能被不小心修改过，需要做校检
func (s *singleRule) addFromBackUpFile(key interface{}, reordFromBackUpFile int64) (err error) {
	index := s.getIndexFrom(key)
	s.visitorRecords[index].deleteExpired(key)
	return s.visitorRecords[index].push(reordFromBackUpFile)
}

// manualEmptyVisitorRecordsOf 清除访问记录
func (s *singleRule) manualEmptyVisitorRecordsOf(key interface{}) {
	index := s.getIndexFrom(key)
	for {
		_, err := s.visitorRecords[index].pop()
		if err != nil {
			break
		}
	}
}

// deleteExpired 删除过期数据
func (s *singleRule) deleteExpired() {
	finished := true
	for range time.Tick(s.cleanupInterval) {
		// 如果数据量较大，那么在一个清除周期内不一定会把所有数据全部清除,所以要判断上一轮次的清除是否完成
		if finished {
			finished = false
			s.deleteExpiredOnce()
			finished = true
		}
	}
}

// deleteExpiredOnce 在特定时间间隔内执行一次删除过期数据操作
func (s *singleRule) deleteExpiredOnce() {
	s.usedVisitorRecordsIndex.Range(func(key, indexVal interface{}) bool {
		index := s.getIndexFrom(key)
		s.visitorRecords[index].deleteExpired(key)
		if s.visitorRecords[index].usedSize() == 0 {
			// 返回数据前，检察空间大小，太大的话，需要清理空间,把空间缩小到默认大小
			s.visitorRecords[index].reSet()
			s.updateIndexOf(key)
		}
		return true
	})
}
