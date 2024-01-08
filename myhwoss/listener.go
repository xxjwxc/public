package myhwoss

import (
	obs "github.com/xxjwxc/public/myhwoss/obs"

	"github.com/xxjwxc/public/mylog"
)

// 定义进度条监听器。
type ObsProgressListener struct {
}

// 定义进度变更事件处理函数。
func (listener *ObsProgressListener) ProgressChanged(event *obs.ProgressEvent) {
	switch event.EventType {
	case obs.TransferStartedEvent:
		mylog.Infof("Transfer Started, ConsumedBytes: %d, TotalBytes %d.\n",
			event.ConsumedBytes, event.TotalBytes)
	case obs.TransferDataEvent:
		mylog.Infof("\rTransfer Data, ConsumedBytes: %d, TotalBytes %d, %d%%.\n",
			event.ConsumedBytes, event.TotalBytes, event.ConsumedBytes*100/event.TotalBytes)
	case obs.TransferCompletedEvent:
		mylog.Infof("\nTransfer Completed, ConsumedBytes: %d, TotalBytes %d.\n",
			event.ConsumedBytes, event.TotalBytes)
	case obs.TransferFailedEvent:
		mylog.Infof("\nTransfer Failed, ConsumedBytes: %d, TotalBytes %d.\n",
			event.ConsumedBytes, event.TotalBytes)
	default:
	}
}
