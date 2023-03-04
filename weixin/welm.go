package weixin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type WeLMInfoReq struct {
	Prompt      string  `json:"prompt"`      // 可选，默认值空字符串，给模型的提示
	Model       string  `json:"model"`       // 必选，要使用的模型名称，当前支持的模型名称有medium、 large 和 xl
	MaxTokens   int     `json:"max_tokens"`  // 可选，最多生成的token个数，默认值 16
	Temperature float64 `json:"temperature"` // 可选 默认值 0.85，表示使用的sampling temperature，更高的temperature意味着模型具备更多的可能性。对于更有创造性的应用，可以尝试0.85以上，而对于有明确答案的应用，可以尝试0（argmax采样）。 建议改变这个值或top_p，但不要同时改变。
	TopP        float64 `json:"top_p"`       // 可选 默认值 0.95，来源于nucleus sampling，采用的是累计概率的方式。即从累计概率超过某一个阈值p的词汇中进行采样，所以0.1意味着只考虑由前10%累计概率组成的词汇。 建议改变这个值或temperature，但不要同时改变。
	TopK        float64 `json:"top_k"`       // 可选 默认值50，从概率分布中依据概率最大选择k个单词，建议不要过小导致模型能选择的词汇少。
	N           int     `json:"n"`           // 可选 默认值 1 返回的序列的个数
	Echo        bool    `json:"echo"`        // 可选 默认值false，是否返回prompt
	Stop        string  `json:"stop"`        // 可选 默认值 null，停止符号。当模型当前生成的字符为stop中的任何一个字符时，会停止生成。若没有配置stop，当模型当前生成的token id 为end_id或生成的token个数达到max_tokens时，停止生成。合理配置stop可以加快推理速度、减少quota消耗。
}

type WeLMInfoResp struct {
	Id      string    `json:"id"`
	Object  string    `json:"object"`
	Created int64     `json:"created"`
	Model   string    `json:"model"`
	Choices []Choices `json:"choices"`
}

type Choices struct {
	Text         string `json:"text"`
	Index        int    `json:"index"`
	Logprobs     int    `json:"logprobs"`
	FinishReason string `json:"finish_reason"`
}

func NewWeLM(config *WeLMInfoReq, authorization string) *weLM {
	return &weLM{
		info:          config,
		authorization: authorization,
	}
}

func GetDefaultWelm() WeLMInfoReq {
	return WeLMInfoReq{
		Prompt:      "",
		Model:       "xl",
		MaxTokens:   16,
		Temperature: 0.85,
		TopP:        0.95,
		TopK:        50,
		N:           1,
		Echo:        false,
		Stop:        "",
	}
}

type weLM struct {
	info          *WeLMInfoReq
	authorization string
}

func (w *weLM) GetAnswer(quest string) string {
	w.info.Prompt = fmt.Sprintf(`问题：%v
	回答：`, quest)

	postData, _ := json.Marshal(w.info)
	client := &http.Client{}
	req, _ := http.NewRequest("POST", "https://welm.weixin.qq.com/v1/completions", bytes.NewReader(postData))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	req.Header.Add("Authorization", w.authorization)
	resp, e := client.Do(req)
	if e != nil {
		return ""
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	// mylog.Debug(string(body))

	var out WeLMInfoResp
	json.Unmarshal(body, &out)
	if len(out.Choices) > 0 {
		return out.Choices[0].Text
	}

	return ""
}
