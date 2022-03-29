package richie
//
//import (
//"bufio"
//"bytes"
//"encoding/binary"
//"io"
//"os"
//"path/filepath"
//"strings"
//"sync"
//"time"
//)
//
//// LoadingAndAutoSaveToDisc 如果有历史备份文件，则加载，无历史备份文件则后续自动生成，并且开启自动保存，默认60秒完成一次存盘
//func (r *Rule) LoadingAndAutoSaveToDisc(backupFileName string, backUpInterval ...time.Duration) {
//	r.loadBackupFileOnce.Do(func() {
//		r.lockerForBackup = new(sync.Mutex)
//		if len(r.rules) == 0 {
//			//panic("rule is empty，please add rule by AddRule")
//		}
//		r.needBackup = true
//		r.backupFileName = strings.Split(backupFileName, ".")[0]
//		if len(backUpInterval) == 0 {
//			//默认60秒存盘一次
//			r.backUpInterval = time.Second * 60
//		} else {
//			r.backUpInterval = backUpInterval[0]
//		}
//		if r.backupFileName == "" {
//			//panic("backupFileName err:" + backupFileName)
//		}
//		//初次运行程序时，无备份文件，不认为是错误
//		// err := r.loading()
//		// if err != nil {
//		// 	if !strings.HasPrefix(err.Error(), "Open backup file fail") {
//		// 		panic(err.Error() + ` please repair or remove the backup file:"` + r.backupFileName + `.ratelimit" and then restart this program.`)
//		// 	}
//		// }
//		var err error
//		go func() {
//			finished := true
//			for range time.Tick(r.backUpInterval) {
//				//如果数据量较大，那么在一个时间周期内不一定会完成存盘操作,所以要判断上一轮次的存盘是否完成
//				if finished {
//					finished = false
//					err = r.SaveToDiscOnce()
//					finished = true
//				}
//			}
//		}()
//	})
//}
//
////把数据保存到硬盘上,仅支持key为string,int,int64等类型数据的缓存
//func (r *Rule) SaveToDiscOnce() (err error) {
//	r.lockerForBackup.Lock()
//	defer r.lockerForBackup.Unlock()
//	if len(r.rules) == 0 {
//		//panic("rule is empty，please add rule by AddRule")
//	}
//	if !r.needBackup {
//		//panic("If you want't to SaveToDiscOnce,you should use LoadingAndAutoSaveToDisc after AddRule.")
//	}
//	f, err := os.Create(r.backupFileName + ".ratelimit_temp")
//	if err != nil {
//		return err
//	}
//	defer os.Remove(r.backupFileName + ".ratelimit_temp")
//	buf := bufio.NewWriterSize(f, 40960)
//	//1 先写规则数量
//	_, err = buf.Write(uint64ToByte(uint64(len(r.rules))))
//	if err != nil {
//		return err
//	}
//	//2 依次写入每一组数据
//	for i := range r.rules {
//		curRuleData := new(bytes.Buffer)
//		tempBuf := bufio.NewWriterSize(curRuleData, 40960)
//		curRuleKeyNum := 0
//		r.rules[i].usedVisitorRecordsIndex.Range(func(key, Index interface{}) bool {
//			index := Index.(int)
//			//备份过程中，不允许其它操作，加锁
//			r.rules[i].visitorRecords[index].locker.Lock()
//			//有效的才能加进去
//			//2.3.1 写入key，key指用户名IP等，只能是数字或string
//			switch key.(type) {
//			case string:
//				//与其它类型不同，KEY长度是不定长的
//				tempBuf.Write([]byte{0x00})
//				tempBuf.Write(uint64ToByte(uint64(len(key.(string)))))
//				tempBuf.WriteString(key.(string))
//			case int:
//				tempBuf.Write([]byte{0x01})
//				tempBuf.Write(uint64ToByte(uint64(key.(int))))
//			case int8:
//				tempBuf.Write([]byte{0x02})
//				tempBuf.Write(uint64ToByte(uint64(key.(int8))))
//			case int16:
//				tempBuf.Write([]byte{0x03})
//				tempBuf.Write(uint64ToByte(uint64(key.(int16))))
//			case int32:
//				tempBuf.Write([]byte{0x04})
//				tempBuf.Write(uint64ToByte(uint64(key.(int32))))
//			case int64:
//				tempBuf.Write([]byte{0x05})
//				tempBuf.Write(uint64ToByte(uint64(key.(int64))))
//			case uint:
//				tempBuf.Write([]byte{0x06})
//				tempBuf.Write(uint64ToByte(uint64(key.(uint))))
//			case uint8:
//				tempBuf.Write([]byte{0x07})
//				tempBuf.Write(uint64ToByte(uint64(key.(uint8))))
//			case uint16:
//				tempBuf.Write([]byte{0x08})
//				tempBuf.Write(uint64ToByte(uint64(key.(uint16))))
//			case uint32:
//				tempBuf.Write([]byte{0x09})
//				tempBuf.Write(uint64ToByte(uint64(key.(uint32))))
//			case uint64:
//				tempBuf.Write([]byte{0x0A})
//				tempBuf.Write(uint64ToByte(key.(uint64)))
//			default:
//				//panic("key type can only be string,int,int8,int16,int32,int64,uint,uint8,uint16,uint32,uint64")
//			}
//			r.rules[i].visitorRecords[index].tailForCopy = r.rules[i].visitorRecords[index].tail
//			r.rules[i].visitorRecords[index].headForCopy = r.rules[i].visitorRecords[index].head
//			size := r.rules[i].visitorRecords[index].usedSize()
//			//2.3.2写下当前key对应的有效访问记录数,为了简单，不判断其是否过期
//			tempBuf.Write(uint64ToByte(uint64(size)))
//			if size > 0 {
//				for ii := 0; ii < size; ii++ {
//					val, _ := r.rules[i].visitorRecords[index].tempQueuePopForCopy()
//					//2.3.3写下每条访问数据的时间点
//					tempBuf.Write(uint64ToByte(uint64(val)))
//				}
//			}
//			curRuleKeyNum++
//			r.rules[i].visitorRecords[index].locker.Unlock()
//			return true
//		})
//		//2.1 //先写当前下标
//		buf.Write(uint64ToByte(uint64(i)))
//		//2.2 再写当前键的个数
//		buf.Write(uint64ToByte(uint64(curRuleKeyNum)))
//		//2.3再写某个键下面的所有数据，如果无数据，则不写
//		//tempBuf由上面提前算出
//		if curRuleKeyNum > 0 {
//			tempBuf.Flush()
//			b := curRuleData.Bytes()
//			buf.Write(b)
//		}
//	}
//	buf.Flush()
//	err = f.Close()
//	if err != nil {
//		return
//	}
//	//成功生成临时文件后，成替换正式文件
//	_, err = copyFile(r.backupFileName+".ratelimit", r.backupFileName+".ratelimit_temp")
//	return
//}
//func uint64ToByte(i uint64) []byte {
//	b := make([]byte, 8)
//	binary.LittleEndian.PutUint64(b, i)
//	return b
//}
//
////复制文件，目标文件所在目录不存在，则创建目录后再复制
////Copy(`d:\test\hello.txt`,`c:\test\hello.txt`)
//func copyFile(dstFileName, srcFileName string) (w int64, err error) {
//	//打开源文件
//	srcFile, err := os.Open(srcFileName)
//	if err != nil {
//		return 0, err
//	}
//	defer srcFile.Close()
//	// 创建新的文件作为目标文件
//	dstFile, err := os.Create(dstFileName)
//	if err != nil {
//		//如果出错，很可能是目标目录不存在，需要先创建目标目录
//		err = os.MkdirAll(filepath.Dir(dstFileName), 0666)
//		if err != nil {
//			return 0, err
//		}
//		//再次尝试创建
//		dstFile, err = os.Create(dstFileName)
//		if err != nil {
//			return 0, err
//		}
//	}
//	defer dstFile.Close()
//	//通过bufio实现对大文件复制的自动支持
//	dst := bufio.NewWriter(dstFile)
//	defer dst.Flush()
//	src := bufio.NewReader(srcFile)
//	w, err = io.Copy(dst, src)
//	if err != nil {
//		return 0, err
//	}
//	return w, err
//}
