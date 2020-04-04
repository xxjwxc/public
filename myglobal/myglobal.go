package myglobal

import (
	"fmt"

	"github.com/xxjwxc/public/myglobal/snowflake"
	"github.com/xxjwxc/public/mylog"
)

// NodeInfo 节点信息
type NodeInfo struct {
	ID        int64
	snowflake *snowflake.Node // 雪花算法获取全局唯一id
}

var node NodeInfo

func init() {
	SetNodeID(1) //默认值
}

// SetNodeID 设置当前机器节点ID(一般程序启动时就需要设置,且初始化)
func SetNodeID(nodeID int64) {
	node.ID = nodeID
	Init()
}

// 初始化
func Init() {
	var err error
	node.snowflake, err = snowflake.NewNode(node.ID)
	if err != nil {
		mylog.Error(err)
	}
}

// GetNode 获取node节点
func GetNode() *NodeInfo {
	return &node
}

// GetID 获取全局唯一id
func (n *NodeInfo) GetID() int64 {
	return n.snowflake.Generate().Int64()
}

// GetIDStr 获取全局唯一id
func (n *NodeInfo) GetIDStr() string {
	return fmt.Sprintf("%v", n.snowflake.Generate().Int64())
}
