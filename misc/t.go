package main

import (
	"fmt"
	"math/rand"
)


type trie struct {
	height, id int
	contents []byte
	subtries []*trie
	subsizes []int
}

const (
	b = 5
	m = 1<<b
)

type plan struct{
	ms, minus_ones, leftover int
}

var (
	rebalance = make([]plan,2*m*m)
)

func sum(p plan) int {
	return p.ms + p.minus_ones
}

func ComputeRebalancePlan(){
	for i:=0; i<m-1; i++ {
		rebalance[i] = plan{0,0,i}
	}
	rebalance[m-1] = plan{0,1,0}
	rebalance[m] = plan{1,0,0}
	for i:=m+1; i<2*m*m; i++ {
		plan_ms := rebalance[i-m]
		plan_minus_ones := rebalance[i-m+1]
		if plan_ms.leftover == 0 {
			rebalance[i] = plan{plan_ms.ms+1, plan_ms.minus_ones, 0}
		} else if plan_minus_ones.leftover == 0 {
			rebalance[i] = plan{plan_minus_ones.ms, plan_minus_ones.minus_ones+1, 0}
		} else {
			rebalance[i] = rebalance[i-1]
			rebalance[i].leftover += 1
		}
	}
}

func NewTrans(h,id int) *trie {
	t := &trie{}
	t.height = h
	t.id = id
	if h == 0 {
		t.contents = make([]byte, 0, m)
	} else {
		t.subtries = make([]*trie, 0, m)
		t.subsizes = make([]int, 0, m)
	}
	return t
}

func NewTrie(h int) *trie {
	return NewTrans(h,0)
}

func Size(t *trie) int {
	if t.height == 0 {
		return len(t.contents)
	} else {
		sum := 0
		for _,s := range t.subsizes {
			sum += s
		}
		return sum
	}
}

func Len(t *trie) int {
	if t.height == 0 {
		return len(t.contents)
	} else {
		return len(t.subtries)
	}
}

func Full(t *trie) bool {
	if t.height > 0 {
		last := len(t.subtries) - 1
		if last < m-1 {
			return false
		} else {
			return Full(t.subtries[last])
		}
	} else {
		return len(t.contents) == m
	}
}

func CloneTrans(t *trie, id int) *trie {
	if t.id == id {
		return t
	}
	tt := *t
	if tt.height == 0 {
		tt.contents = append([]byte(nil), t.contents...)
	} else {
		tt.subtries = append([]*trie(nil), t.subtries...)
		tt.subsizes = append([]int(nil), t.subsizes...)
	}
	return &tt
}

func Clone(t *trie) *trie {
	return CloneTrans(t,0)
}

func Append_h(t *trie, h, id int, v byte) *trie {
	// t is never full
	if t == nil {
		if h == 0 {
			n := NewTrans(0, id)
			n.contents = append(n.contents, v)
			return n
		} else {
			n := NewTrans(h, id)
			n.subtries = append(n.subtries, Append_h(nil, h-1, id, v))
			n.subsizes = append(n.subsizes, 1)
			return n
		}
	}

	if h == 0 {
		n := CloneTrans(t, id)
		n.contents = append(n.contents, v)
		return n
	} else {
		last := len(t.subtries)-1
		if Full(t.subtries[last]) {
			n := CloneTrans(t, id)
			ns := Append_h(nil, h-1, id, v)
			n.subtries = append(n.subtries, ns)
			n.subsizes = append(n.subsizes, Size(ns))
			return n
		} else {
			n := CloneTrans(t, id)
			ns := Append_h(n.subtries[last], h-1, id, v)
			n.subtries[last] = ns
			n.subsizes[last] = Size(ns)
			return n
		}
	}
}

func Append(t *trie, v byte) *trie {
	if Full(t) {
		r := NewTrans(t.height+1, t.id)
		r.subtries = append(r.subtries, t)
		r.subsizes = append(r.subsizes, Size(t))
		return Append_h(r, r.height, r.id, v)
	} else {
		return Append_h(t, t.height, t.id, v)
	}
}

func Trans(t *trie) *trie {
	if t.id != 0 {
		return t
	}
	tt := Clone(t)
	tt.id = rand.Int()
	return tt
}

func Persistent(t *trie) *trie {
	t.id = 0
	return t
}

func AppendSlice(t *trie, b []byte) *trie {
	var tt *trie
	if t == nil {
		tt = NewTrans(0, rand.Int())
	} else {
		tt = Trans(t)
	}
	for _, bb := range b {
		tt = Append(tt, bb)
	}
	return Persistent(tt)
}

func PrintTrie_h(t *trie, d int) string {
	if t == nil {
		return "<nil>\n"
	}
	s := ""
	for i:=0; i<d; i++ {
		s += " "
	}
	if t.height == 0 {
		s += fmt.Sprintf("%p h:%d contents:%p \"%s\"\n",
			t, t.height, &t.contents, t.contents)
	} else {
		s += fmt.Sprintf("%p h:%d %p %#v sizes:%#v\n", 
			t, t.height, &t.subtries, t.subtries, t.subsizes)
	}
	if t.height > 0 {
		for _, ts := range t.subtries {
			s += PrintTrie_h(ts, d+1)
		}
	}
	return s
}

func PrintTrie(t *trie) string {
	return PrintTrie_h(t, 0)
}

func Read_h(t *trie, p[]byte) (int,bool) {
	if t.height == 0{
		c := copy(p,t.contents)
		return c, c<len(t.contents)
	} else {
		n := 0
		for _, st := range t.subtries {
			c,end := Read_h(st, p[n:])
			n += c
			if end {
				return n,end
			}
		}
		return n, false
	}
}

func (t *trie) Read(p []byte) (int, error) {
	n,_ := Read_h(t, p)
	return n, nil
}

func Take(t *trie, i int) *trie {
	if t.height == 0 {
		n := Clone(t)
		n.contents = n.contents[:i]
		return n
	}
	var cutoff int
	for si,s := range t.subsizes {
		if i < s {
			cutoff = si
			break
		} else {
			i -= s
		}
	}

	if cutoff == 0 {
		return Take(t.subtries[0],i)
	} else {
		n := Clone(t)
		n.subtries = n.subtries[:cutoff+1]
		n.subsizes = n.subsizes[:cutoff+1]
		n.subtries[cutoff] = Take(n.subtries[cutoff],i)
		n.subsizes[cutoff] = i
		return n
	}
}

func Drop(t *trie, i int) *trie {
	if t.height == 0 {
		n := Clone(t)
		n.contents = n.contents[i:]
		return n
	}
	var cutoff int
	for si,s := range t.subsizes {
		if i < s {
			cutoff = si
			break
		} else {
			i -= s
		}
	}
	
	if cutoff == len(t.subtries)-1 {
		return Drop(t.subtries[cutoff],i)
	} else {
		n := Clone(t)
		n.subtries = n.subtries[cutoff:]
		n.subsizes = n.subsizes[cutoff:]
		n.subtries[0] = Drop(n.subtries[0],i)
		n.subsizes[0] = n.subsizes[0] - i
		return n
	}
}

func Concat_h(l *trie, r *trie) []*trie {
	// both tries should be of the same height
	// return nodes of a height lower than l and r
	if l.height == 1 { // the subtries in l and r have height 0
		total_length := Size(l)
		total_length += Size(r)
		plan := rebalance[total_length]
		new_subtries_len := plan.ms + plan.minus_ones
		if plan.leftover != 0 {
			new_subtries_len++
		}
		new_subtries := make([]*trie, 0, new_subtries_len)

		l_copied := 0
		r_copied := 0
		reshuffle := false
		for i,ls := range l.subtries {
			l_copied = i
			if Size(ls) == m && plan.ms > 0 {
				new_subtries = append(new_subtries,ls)
				plan.ms--
			} else if Size(ls) == m-1 && plan.minus_ones > 0{
				new_subtries = append(new_subtries,ls)
				plan.minus_ones--
			} else {
				reshuffle = true
				break
			}
		}

		if !reshuffle {
			for i,ls := range r.subtries {
				r_copied = i
				if Size(ls) == m && plan.ms > 0 {
					new_subtries = append(new_subtries,ls)
					plan.ms--
				} else if Size(ls) == m-1 && plan.minus_ones > 0{
					new_subtries = append(new_subtries,ls)
					plan.minus_ones--
				} else {
					break
				}
			}
		}

		c := make([]byte, 0,
			plan.ms*m + plan.minus_ones*(m-1) + plan.leftover)
		for i := l_copied; i<len(l.subtries); i++ {
			c = append(c, l.subtries[i].contents...)
		}
		for i := r_copied; i<len(r.subtries); i++ {
			c = append(c, r.subtries[i].contents...)
		}

		counter := 0
		for i:=0; i<plan.ms; i++ {
			n := NewTrie(0)
			n.contents = c[counter:counter+m]
			new_subtries = append(new_subtries,n)
			counter += m
		}
		for i:=0; i<plan.minus_ones; i++ {
			n := NewTrie(0)
			n.contents = c[counter:counter+m-1]
			new_subtries = append(new_subtries,n)
			counter += m-1
		}
		if plan.leftover != 0{
			n := NewTrie(0)
			n.contents = c[counter:]
			new_subtries = append(new_subtries,n)
		}
		return new_subtries
	}
	
	last := len(l.subtries)-1
	llast := l.subtries[last]
	rfirst := r.subtries[0]
	middle_subsubtries := Concat_h(llast, rfirst)
	total_length := 0
	for i:=0; i<last; i++ {
		total_length += Len(l.subtries[i])
	}
	total_length += len(middle_subsubtries)
	for i:=1; i<len(r.subtries); i++ {
		total_length += Len(r.subtries[i])
	}
	plan := rebalance[total_length]
	new_subtries_len := plan.ms + plan.minus_ones
	if plan.leftover != 0 {
		new_subtries_len++
	}
	new_subtries := make([]*trie, 0, new_subtries_len)

	l_copied := 0
	for i:=0; i < last; i++ {
		l_copied = i
		ls := l.subtries[i]
		if Len(ls) == m && plan.ms > 0 {
			new_subtries = append(new_subtries,ls)
			plan.ms--
		} else if Len(ls) == m-1 && plan.minus_ones > 0{
			new_subtries = append(new_subtries,ls)
			plan.minus_ones--
		} else {
			break
		}
	}

	c := make([]*trie, 0,
		plan.ms*m + plan.minus_ones*(m-1) + plan.leftover)
	s := make([]int, 0,
		plan.ms*m + plan.minus_ones*(m-1) + plan.leftover)
	for i := l_copied; i<last-1; i++ {
		c = append(c, l.subtries[i].subtries...)
		s = append(s, l.subtries[i].subsizes...)
	}
	c = append(c, middle_subsubtries...)
	for i := 0; i<len(middle_subsubtries); i++ {
		s = append(s,Size(middle_subsubtries[i]))
	}
	for i := 1; i<len(r.subtries); i++ {
		c = append(c, r.subtries[i].subtries...)
		s = append(s, r.subtries[i].subsizes...)
	}

	counter := 0
	for i:=0; i<plan.ms; i++ {
		n := NewTrie(l.height - 1)
		n.subtries = c[counter:counter+m]
		n.subsizes = s[counter:counter+m]
		new_subtries = append(new_subtries,n)
		counter += m
	}
	for i:=0; i<plan.minus_ones; i++ {
		n := NewTrie(l.height - 1)
		n.subtries = c[counter:counter+m-1]
		n.subsizes = s[counter:counter+m-1]
		new_subtries = append(new_subtries,n)
		counter += m-1
	}
	if plan.leftover != 0 {
		n := NewTrie(l.height -1)
		n.subtries = c[counter:]
		n.subsizes = s[counter:]
		new_subtries = append(new_subtries,n)
	}
	return new_subtries
}

func Concat(l *trie, r *trie) *trie {
	if r.height == 0 {
		return AppendSlice(l,r.contents)
	}
	for l.height > r.height {
		rr := NewTrie(r.height + 1)
		rr.subtries = append(r.subtries, r)
		rr.subsizes = append(r.subsizes, Size(r))
		r = rr
	}
	for l.height < r.height {
		ll := NewTrie(l.height + 1)
		ll.subtries = append(l.subtries, l)
		ll.subsizes = append(l.subsizes, Size(l))
		l = ll
	}

	new_subtries := Concat_h(l,r)
	if len(new_subtries) <= m {
		new_root := NewTrie(l.height)
		new_root.subtries = new_subtries
		for _,ns := range new_subtries {
			new_root.subsizes = append(new_root.subsizes, Size(ns))
		}
		return new_root
	} else {
		new_lnode := NewTrie(l.height)
		new_lnode.subtries = new_subtries[:m]
		for _,ns := range new_lnode.subtries {
			new_lnode.subsizes = append(new_lnode.subsizes, Size(ns))
		}
		new_rnode := NewTrie(r.height)
		new_rnode.subtries = new_subtries[m:]
		for _,ns := range new_rnode.subtries {
			new_rnode.subsizes = append(new_rnode.subsizes, Size(ns))
		}
		new_root := NewTrie(l.height+1)
		new_root.subtries = append(new_root.subtries, new_lnode, new_rnode)
		new_root.subsizes = append(new_root.subsizes, Size(new_lnode), Size(new_rnode))
		return new_root
	}
}

func Lookup(t *trie, i int) byte {
	if t.height == 0 {
		return t.contents[i]
	}
	for si,s := range t.subsizes {
		if i < s {
			return Lookup(t.subtries[si],i)
		} else {
			i -= s
		}
	}
	panic(fmt.Errorf("Out of bounds with i = %d at trie %s", i, PrintTrie(t)))
}

func main(){
	ComputeRebalancePlan()
	ss := "Lorem Ipsum is a placeholder text commonly used to demonstrate the visual form of a document or a typeface without relying on meaningful content."
	s := AppendSlice(nil, []byte(ss))
	fmt.Println(PrintTrie(s))
	r := Concat(s,s)
	fmt.Print(PrintTrie(r))
}
