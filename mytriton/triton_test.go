package mytriton

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestMain(t *testing.T) {
	triton, err := NewTritonServer("192.155.1.93:18001", "model_pipe", "", 10*time.Second)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(triton.ServerLive(context.Background()))
	fmt.Println(triton.RequestFromText(context.Background(), "两点钟,下午3点肉沫豆腐怎么做", "embeddings"))
}
