package wordsfilter

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWordsFilter(t *testing.T) {
	texts := []string{
		"爱女人",
		"爱液",
		"你他妈",
		"ความรุ่งโรจน์",
	}
	wf := New()
	root := wf.Generate(texts)

	// 命中:文本中包含完整敏感词
	if !wf.Contains("你的爱液", root) {
		t.Error("Contains 应命中敏感词『爱液』")
	}

	// 未命中:正常文本
	if wf.Contains("这是完全正常的句子", root) {
		t.Error("Contains 不应命中正常文本")
	}

	// 未命中:仅覆盖敏感词前缀,不是完整词
	if wf.Contains("爱", root) {
		t.Error("Contains 不应命中不完整的敏感词前缀")
	}

	// Replace:StripSpace 默认开启,先去掉空格,敏感词按字符数替换为 *
	if got, want := wf.Replace("你他妈", root), "***"; got != want {
		t.Errorf("Replace 结果不符,期望 %q,实际 %q", want, got)
	}
	if got, want := wf.Replace("你 他妈", root), "***"; got != want {
		t.Errorf("Replace 应先去空格,期望 %q,实际 %q", want, got)
	}

	// Remove 之后不再命中
	wf.Remove("你他妈", root)
	if wf.Contains("你他妈", root) {
		t.Error("Remove 之后 Contains 不应再命中")
	}
}

func TestWordsFilterGenerateWithFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "words.txt")
	if err := os.WriteFile(path, []byte("测试词\n敏感\n\n  \n"), 0644); err != nil {
		t.Fatal(err)
	}

	wf := New()
	root, err := wf.GenerateWithFile(path)
	if err != nil {
		t.Fatalf("GenerateWithFile 失败: %v", err)
	}
	if !wf.Contains("这里包含测试词", root) {
		t.Error("从文件加载的词表应能命中『测试词』")
	}
	if wf.Contains("这里没有命中", root) {
		t.Error("不应命中词表中不存在的词")
	}
}
