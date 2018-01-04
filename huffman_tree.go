package simu

import (
	"sort"
	"strings"
)

type Pair struct {
	Key   string
	Value int
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func rankByWordCount(wordFrequencies map[string]int) PairList {
	pl := make(PairList, len(wordFrequencies))
	i := 0
	for k, v := range wordFrequencies {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

func CharCount(characters string) map[string]int {
	results := make(map[string]int)
	s := strings.Split(characters, "")
	for _, c := range s {
		if _, ok := results[c]; ok {
			results[c] = results[c] + 1
		} else {
			results[c] = 1
		}
	}
	return results
}

func hash(p Pair) int {
	return p.Value
}

func HuffmanTree(characters string) error {
	charCount := rankByWordCount(CharCount(characters))
	bt, err := NewBinaryTree(charCount[len(charCount)-1].Key)
	if err != nil {
		return err
	}
	for k, _ := range charCount {
		bt.Insert(k)
	}
	return nil
}
