package mytriton

import (
	"bytes"
	"context"
	"encoding/binary"
	"time"

	triton "github.com/xxjwxc/public/mytriton/grpc-client"

	"google.golang.org/grpc"
)

type TritonInfo struct {
	modelName    string
	modelVersion string
	conn         *grpc.ClientConn
	client       triton.GRPCInferenceServiceClient
	timeout      time.Duration
}

func NewTritonServer(url, modelName, modelVersion string, timeout time.Duration) (*TritonInfo, error) {
	// Connect to gRPC server
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	// Create client from gRPC server connection
	client := triton.NewGRPCInferenceServiceClient(conn)

	return &TritonInfo{
		modelName: modelName,
		conn:      conn,
		client:    client,
		timeout:   timeout,
	}, nil
}

func (t *TritonInfo) Close() {
	t.conn.Close()
}

// ServerLive 心跳检测
func (t *TritonInfo) ServerLive(ctx context.Context) (bool, error) {
	// Create context for our request with 10 second timeout
	ctx, cancel := context.WithTimeout(ctx, t.timeout)
	defer cancel()

	serverLiveRequest := triton.ServerLiveRequest{}
	// Submit ServerLive request to server
	serverLiveResponse, err := t.client.ServerLive(ctx, &serverLiveRequest)
	if err != nil {
		return false, err
	}
	return serverLiveResponse.Live, nil
}

func (t *TritonInfo) RequestFromText(ctx context.Context, name, text string, outTensorsName string) ([]byte, error) {
	l := len(text)
	b := make([]byte, 4)
	binary.LittleEndian.PutUint16(b, uint16(l))
	b = append(b, []byte(text)...)

	ctx, cancel := context.WithTimeout(ctx, t.timeout)
	defer cancel()

	inferInputs := []*triton.ModelInferRequest_InferInputTensor{
		{
			Name:     name,
			Datatype: "BYTES",
			Shape:    []int64{1, 1},
		},
	}

	// Create request input output tensors
	inferOutputs := []*triton.ModelInferRequest_InferRequestedOutputTensor{
		{
			Name: outTensorsName,
		},
	}

	// Create inference request for specific model/version
	modelInferRequest := triton.ModelInferRequest{
		ModelName:    t.modelName,
		ModelVersion: t.modelVersion,
		Inputs:       inferInputs,
		Outputs:      inferOutputs,
	}

	modelInferRequest.RawInputContents = append(modelInferRequest.RawInputContents, b)

	// Submit inference request to server
	modelInferResponse, err := t.client.ModelInfer(ctx, &modelInferRequest)
	if err != nil {
		return nil, err
	}

	return modelInferResponse.RawOutputContents[0], nil
}

func (t *TritonInfo) BytesToFloat32(outputBytes []byte) []float32 {
	ff := 4
	size := len(outputBytes) / ff

	outputData0 := make([]float32, size)
	// outputData1 := make([]int64, outputSize)
	for i := 0; i < size; i++ {
		buf := bytes.NewBuffer(outputBytes[i*ff : i*ff+ff])
		var retval float32
		binary.Read(buf, binary.LittleEndian, &retval)
		outputData0[i] = retval
	}
	return outputData0
}

// func (t *TritonInfo) BytesToFloat64(outputBytes []byte) []float64 {
// 	ff := 4
// 	size := len(outputBytes) / ff

// 	outputData0 := make([]float64, size)
// 	// outputData1 := make([]int64, outputSize)
// 	for i := 0; i < size; i++ {
// 		buf := bytes.NewBuffer(outputBytes[i*ff : i*ff+ff])
// 		var retval float32
// 		binary.Read(buf, binary.LittleEndian, &retval)
// 		outputData0[i] = float64(retval)
// 	}
// 	return outputData0
// }

func (t *TritonInfo) BytesToString(outputBytes []byte) (out []string) {
	ff := 4
	size := len(outputBytes)
	i := 0
	for i < size {
		buf := bytes.NewBuffer(outputBytes[i : i+ff])
		var _len int32
		binary.Read(buf, binary.LittleEndian, &_len)
		if _len < int32(size) {
			out = append(out, string(outputBytes[i+ff:i+ff+int(_len)]))
		}

		i += ff + int(_len)
	}
	return
}

func (t *TritonInfo) RequestFromTexts(ctx context.Context, texts []string, outTensorsName string) ([]float32, error) {
	var _bytes []byte
	for _, text := range texts {
		l := len(text)
		b := make([]byte, 4)
		binary.LittleEndian.PutUint16(b, uint16(l))
		b = append(b, []byte(text)...)
		_bytes = append(_bytes, b...)
	}

	ctx, cancel := context.WithTimeout(ctx, t.timeout)
	defer cancel()

	inferInputs := []*triton.ModelInferRequest_InferInputTensor{
		{
			Name:     "texts",
			Datatype: "BYTES",
			Shape:    []int64{1, int64(len(texts))},
		},
	}

	// Create request input output tensors
	inferOutputs := []*triton.ModelInferRequest_InferRequestedOutputTensor{
		{
			Name: outTensorsName,
		},
	}

	// Create inference request for specific model/version
	modelInferRequest := triton.ModelInferRequest{
		ModelName:    t.modelName,
		ModelVersion: t.modelVersion,
		Inputs:       inferInputs,
		Outputs:      inferOutputs,
	}

	modelInferRequest.RawInputContents = append(modelInferRequest.RawInputContents, _bytes)

	// Submit inference request to server
	modelInferResponse, err := t.client.ModelInfer(ctx, &modelInferRequest)
	if err != nil {
		return nil, err
	}

	outputBytes0 := modelInferResponse.RawOutputContents[0]
	ff := 4
	size := len(outputBytes0) / ff

	outputData0 := make([]float32, size)
	// outputData1 := make([]int64, outputSize)
	for i := 0; i < size; i++ {
		buf := bytes.NewBuffer(outputBytes0[i*ff : i*ff+ff])
		var retval float32
		binary.Read(buf, binary.LittleEndian, &retval)
		outputData0[i] = retval
	}
	return outputData0, nil
}
