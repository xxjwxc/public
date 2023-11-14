package mytriton

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestMain(t *testing.T) {
	triton, err := NewTritonServer("192.155.1.93:18001", "model_ensemble_pre", "", 10*time.Second)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(triton.ServerLive(context.Background()))
	out, _ := triton.RequestFromText(context.Background(), "texts", "你好", "sentence_embedding")
	f32 := triton.BytesToFloat64(out)
	fmt.Println(f32)
}

func TestMain1(t *testing.T) {
	triton, err := NewTritonServer("192.155.1.93:18001", "dev_ftt_nlu_332", "", 10*time.Second)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(triton.ServerLive(context.Background()))
	out, _ := triton.RequestFromText(context.Background(), "TEXT", "推荐食谱", "intent_classification")
	_strs := triton.BytesToString(out)
	fmt.Println(_strs)
}
