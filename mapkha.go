/*
 * Copyright (c) 2015, Vee Satayamas
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice, this 
 *    list of conditions and the following disclaimer.
 * 
 * 2. Redistributions in binary form must reproduce the above copyright notice,
 * this list of conditions and the following disclaimer in the documentation
 * and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE 
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
 * LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR 
 * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF 
 * SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS 
 * INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN 
 * CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF
 * THE POSSIBILITY OF SUCH DAMAGE.
 */

package main

import ("fmt"
	"io/ioutil"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func LoadDict(path string) ([][]rune, error) {
	b_slice, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	data := string(b_slice)
	swords := strings.Split(data, "\n")
	rwords := make([][]rune, len(swords))
	for i, word := range swords {
		rwords[i] = []rune(word)
	}
	return rwords, nil
}

type DictAcceptor struct {
	l int
	r int
	final bool
	offset int
	valid bool
}

func (a *DictAcceptor) Transit(ch rune, dict [][]rune) {
	var found bool
	a.l, found = DictSeek(LEFT, dict, a.l, a.r, a.offset, ch)
	if found {
		a.r, _ = DictSeek(RIGHT, dict, a.l, a.r, a.offset, ch)
		a.offset++
		w := dict[a.l]
		wlen := len(w)
		a.final = (wlen == a.offset)
	} else {
		a.valid = false
	}
}

func DictSeek(policy int, dict [][]rune, l int, r int, offset int, ch rune) (int, bool) {
	ans := 0
	found := false
	m := 0
	
	if policy != LEFT && policy != RIGHT {
		return 0, found
	}
	
	for {
		if l > r {
			break
		}
		m = (l+r) / 2
		w := dict[m]
		wlen := len(w)
		if wlen <= offset {
			l = m + 1
		} else {
			ch_ := w[offset]
			if ch_ < ch {
				l = m + 1
			} else if ch_ > ch {
				r = m - 1
			} else {
				ans = m
				found = true
				switch policy {
				case LEFT: r = m - 1
				case RIGHT: l = m + 1
				}
			}			
		}
	}
	
	return ans, found
}

type TextRange struct {
	s int
	e int
}

type Edge struct {
	w int
	unk int
	p int
}

const (
	LEFT = 1
	RIGHT = 2
)

func TransitAll(acc []DictAcceptor, ch rune, dict [][]rune) []DictAcceptor {
	_acc := append(acc, DictAcceptor{0, len(dict)-1, false, 0, true})
	__acc := make([]DictAcceptor, 0, len(_acc))
	for _, a := range(_acc) {
		a.Transit(ch, dict)
		if a.valid {
			__acc = append(__acc, a)
		}
	}
	return __acc
}

func Better(a *Edge, b *Edge) bool {
	if a.unk < a.unk {
		return true
	} else {
		if a.w < b.w {
			return true
		}
	}
	return false
}

func BestEdge(edges []Edge) *Edge {
	l := len(edges)
	if l == 0 {
		return nil
	}

	e := &edges[0]

	for i := 1; i < l; i++ {
		if Better(&edges[i], e) {
			e = &edges[i];
		}
	}
	
	return e
}

func BuildEdges(i int, acc []DictAcceptor, g []Edge) []Edge {
	edges := make([]Edge, 0, len(acc))
	for _, a := range(acc) {
		if a.final {
			p := i - a.offset + 1
			src := g[p]
			edge := Edge{src.w + 1, src.unk, p}
			edges = append(edges, edge)			
		}
	}

	if len(edges) == 0 {
		edge := Edge{100, 100, 0}
		edges = append(edges, edge)
	}
	return edges
}

func BuildGraph(t []rune, dict [][]rune) []Edge {
	g := make([]Edge, len(t) + 1)
	g[0] = Edge{0, 0, -1}
	var acc []DictAcceptor
	for i, ch := range(t) {
		acc = TransitAll(acc, ch, dict)
		edges := BuildEdges(i, acc, g)
		e := BestEdge(edges)
		g[i+1] = *e 
	}
	return g
}

func GraphToRanges(g []Edge) []TextRange {
	_ranges := make([]TextRange, 0, len(g))
	e := len(g) - 1
	var s int
	for {
		if e <= 0 {
			break
		}

		s = g[e].p

		r := TextRange{s,e}

		_ranges = append(_ranges, r)
		
		e = s
	}

	l := len(_ranges)
	
	ranges := make([]TextRange, l)

	for i, r := range(_ranges) {
		ranges[l - i - 1] = r
	}

	return ranges
}

func FindRanges(t []rune, dict [][]rune) []TextRange {
	g := BuildGraph(t, dict)
	ranges := GraphToRanges(g)
	return ranges
}

func Segment(_t string, dict [][]rune) ([]string, error) {
	t := []rune(_t)
	ranges := FindRanges(t, dict)
	wlst := make([]string, len(ranges))
	for i, r := range ranges {
		wlst[i] = string(t[r.s:r.e])
	}
	return wlst, nil
}

func main() {
	dict, e := LoadDict("tdict-std.txt")
	check(e)
	wl, e := Segment("ตัดคำไหม", dict)
	fmt.Println(wl)
}