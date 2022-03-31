package wordsfilter

import (
	"testing"
)

func TestWordsFilter(t *testing.T) {
	texts := []string{
		"Miyamoto Musashi",
		"妲己",
		"アンジェラ",
		"ความรุ่งโรจน์",
	}
	wf := New()
	root := wf.Generate(texts)
	wf.Remove("shif", root)
	c1 := wf.Contains("アン", root) // 是否有敏感词
	if c1 != false {
		t.Errorf("Test Contains expect false, get %T, %v", c1, c1)
	}
	c2 := wf.Contains("->アンジェラ2333", root)
	if c2 != true {
		t.Errorf("Test Contains expect true, get %T, %v", c2, c2)
	}
	r1 := wf.Replace("Game ความรุ่งโรจน์ i like 妲己 heroMiyamotoMusashi", root)
	if r1 != "Game*************ilike**hero***************" {
		t.Errorf("Test Replace expect Game*************ilike**hero***************,get %T,%v", r1, r1)
	}
	// Test generated with file.
	root, _ = wf.GenerateWithFile("./words_test.txt")
	if wf.Contains("アンジェラ", root) != true {
		t.Errorf("Test Contains expect true, get %T, %v", c2, c2)
	}
}
