package wordsfilter

import (
	"bytes"
	"strings"
)

type Node struct {
	Child        map[string]*Node
	Placeholders string
}

// New creates a node.
func NewNode(child map[string]*Node, placeholders string) *Node {
	return &Node{
		Child:        child,
		Placeholders: placeholders,
	}
}

// Add sensitive words to specified sensitive words Map.
func (node *Node) add(text string, root map[string]*Node, placeholder string) {
	if text == "" {
		return
	}
	textr := []rune(text)
	end := len(textr) - 1
	for i := 0; i <= end; i++ {
		word := string(textr[i])
		if n, ok := root[word]; ok { // contains key
			if i == end { // the last
				n.Placeholders = strings.Repeat(placeholder, end+1)
			} else {
				if n.Child != nil {
					root = n.Child
				} else {
					root = make(map[string]*Node)
					n.Child = root
				}
			}
		} else {
			placeholders, child := "", make(map[string]*Node)
			if i == end {
				placeholders = strings.Repeat(placeholder, end+1)
			}
			root[word] = NewNode(child, placeholders)
			root = child
		}
	}
}

// Remove specified sensitive words from sensitive word map.
func (node *Node) remove(text string, root map[string]*Node) {
	textr := []rune(text)
	end := len(textr) - 1
	for i := 0; i <= end; i++ {
		word := string(textr[i])
		if n, ok := root[word]; ok {
			if i == end {
				n.Placeholders = ""
			} else {
				root = n.Child
			}
		} else {
			return
		}
	}
}

// Replace sensitive words in strings and return new strings.
// Follow the principle of maximum matching.
func (node *Node) replace(text string, root map[string]*Node) string {
	if root == nil || text == "" {
		return text
	}
	textr := []rune(text)
	i, s, e, l := 0, 0, 0, len(textr)
	bf := bytes.Buffer{}
	words := make(map[string]*Node)
	var back []*Node
loop:
	for e < l {
		words = root
		i = e
		// Maximum Matching Principle, Matching Backwards First
		for ; i < l; i ++ {
			word := string(textr[i])
			if n, ok := words[word]; ok {
				back = append(back, n)
				if n.Child != nil {
					words = n.Child
				} else if n.Placeholders != "" {
					bf.WriteString(string(textr[s:e]))
					bf.WriteString(n.Placeholders)
					i++
					s, e = i, i
					continue loop
				} else {
					break
				}
			} else if n != nil && n.Placeholders != "" {
				bf.WriteString(string(textr[s:e]))
				bf.WriteString(n.Placeholders)
				s, e = i, i
				continue loop
			} else {
				break
			}
		}
		// Backward match fails, backtracking.
		for ; i > e; i-- {
			bl := len(back)
			if bl == 0 {
				break
			}
			last := back[bl-1]
			back = back[:bl-1]
			if last.Placeholders != "" {
				bf.WriteString(string(textr[s:e]))
				bf.WriteString(last.Placeholders)
				s, e = i, i
				continue loop
			}
		}

		e++
		back = back[:0]
	}
	bf.WriteString(string(textr[s:e]))

	return bf.String()
}

// Whether the string contains sensitive words.
func (node *Node) contains(text string, root map[string]*Node) bool {
	if root == nil || text == "" {
		return false
	}
	textr := []rune(text)
	end := len(textr) - 1
	for i := 0; i <= end; i++ {
		word := string(textr[i])
		if n, ok := root[word]; ok {
			if i == end {
				return n.Placeholders != ""
			} else {
				if len(n.Child) == 0 { // last
					return true
				}
				root = n.Child
			}
		} else {
			continue
		}
	}
	return false
}
