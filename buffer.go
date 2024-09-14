/* This trie structure is the relaxed radix balanced tree, which can be
   seen as a vector that supports persistence, and fast random order
   inserts, lookup and deletes.
   The trie structure contains either m elemeents, or m subtries. It's 
   basically a B-Tree with a branching factor of m, and to keep the height
   of the trie minimal and allow for fast random inserts, the trie that is
   full can sometimes have m-1 elements.
   The trie is also persistent, meaning that changing the data structure 
   results in a new data structure, and the old one is untouched. This trie
   structure allows us to keep the data to a minimum when a new trie is
   created.
   the data structure is from P. Bagwell (2011)
   the implementation for transience follows L'Orange (2014)
   although there are some changes made to fit go:
   + Both Bagwell and L'Orange embeds the array of elements/subtries inside
     the struct. JPB Puente (2017) also does this with an extra improvement:
     embedding more elements into the trie contents. This is harder to do in
     go (though not impossible, it requires the use of unsafe. bookkeeping 
     must be done to keep track of the amount of contents stored in the struct
     ). An obvious concern is taking up extra memory, as well as forcing the 
     use of more pointers. A slice in go is just a pointer to an underlying
     array. Though considering that trie and content/subtrie slice are 
     initialized together, and that slices are pre allocated the required
     capacity and are never reallocated, cache locality may be preserved and
     the performance loss may be negligible.
   + An aside with Puente (2017) is that jamming extra elements into the lowest
     layer of the program may hinder concat performance. A 32 branch-factor trie
     on a 64 bit machine could store 8 more bytes in place of a single subtrie,
     resulting in 256 bytes in a single trie. If we were to concat two large
     enough tries we would end up reshuffling 65536 bytes in the worst case. A
     single byte vector would have to be significantly larger than that to see
     the benefits of rrbt concatenation.

 OPEN: try JPB Puente's embedded implementation?
 */
package web
import (
	"fmt"
	"io"
	"math/rand"
)

func assert(cond bool){
	if !cond {
		panic("assert failed!")
	}
}


const (
	b = 5
	m = 1<<b
)

type Trie[T any] struct {
	height, length, id int
	content *[m]T
	subtrie *[m]*Trie[T]
	subsize *[m]int
}

func NewTrans[T any](h, id int) *Trie[T] {
	a := &Trie[T]{}
	a.id = id
	a.length = 0
	a.height = h
	if h == 0 {
		a.content = &[m]T{}
	} else {
		a.subtrie = &[m]*Trie[T]{}
		a.subsize = &[m]int{}
	}
	return a
}

func NewTrie[T any](h int) *Trie[T] {
	return NewTrans[T](h,0)
}

func (t *Trie[T])Full() bool {
	if t.height > 0 {
		if t.length < m {
			return false
		} else {
			return t.subtrie[m-1].Full()
		}
	}
	return t.length == m
}

// TODO: Fix size so that it stores the amount of elements in the trie up to
//       the subtrie
func (t *Trie[T])Size() int {
	if t.length == 0 {
		return 0
	}
	if t.height == 0 {
		return t.length
	}
	return t.subsize[t.length-1]
}

func clone_array[T any](arr *[m]T, length int) *[m]T {
	if arr == nil {
		return nil
	}
	new_arr := [m]T{}
	for i:=0; i<length; i++ {
		new_arr[i] = arr[i]
	}
	return &new_arr
}

func (t *Trie[T])CloneTrans(id int) *Trie[T] {
	if t.id != 0 && t.id == id {
		return t
	}
	tt := NewTrans[T](t.height, t.id)
	tt.length  = t.length
	tt.content = clone_array[T](t.content, t.length)
	tt.subtrie = clone_array[*Trie[T]](t.subtrie, t.length)
	tt.subsize = clone_array[int](t.subsize, t.length)
	return tt
}

func (t *Trie[T])Clone() *Trie[T] {
	return t.CloneTrans(0)
}

func (t *Trie[T])Trans() *Trie[T] {
	n := t.Clone()
	n.id = rand.Int()
	return n
}

func (t *Trie[T])String() string {
	return fmt.Sprintf("%p %d %d %d", t,t.height,
		t.length, t.id)
}

/* 
  Appending elements to Tries are simple enough
 */

func (t *Trie[T])AppendContent(id int, v T) *Trie[T] {
	n := t.CloneTrans(id)
	assert(n.height == 0)
	assert(n.length < m)
	n.content[n.length] = v
	n.length++
	return n
}

func (t *Trie[T])AppendSubTrie(id int, st *Trie[T]) *Trie[T] {
	n := t.CloneTrans(id)
	assert(n.height > 0)
	assert(n.length < m)
	n.subtrie[n.length] = st
	n.subsize[n.length] = t.Size() + st.Size()
	n.length++
	return n
}

func NewTrieWithElement[T any](h,id int, v T) *Trie[T] {
	n := NewTrans[T](h,id)
	if h == 0 {
		return n.AppendContent(id, v)
	}
	return n.AppendSubTrie(id, NewTrieWithElement(h-1,id,v))
}

func (t *Trie[T])Append(id int, v T) *Trie[T] {
	if t.Full() {
		root := NewTrans[T](t.height+1,id)
		return root.AppendSubTrie(id, t).Append(id, v)
	}
	if t.height == 0 {
		return t.AppendContent(id,v)
	}
	n := t.CloneTrans(id)
	if n.subtrie[n.length-1].Full() {
		return n.AppendSubTrie(id, NewTrieWithElement(n.height-1,id,v))
	}
	n.subtrie[n.length-1] = n.subtrie[n.length-1].Append(id, v)
	n.subsize[n.length-1]++
	return n
}

func (t *Trie[T])AppendSlice(vs []T) *Trie[T] {
	n := t.Trans()
	n.AppendSliceTrans(vs)
	return n
}

func (t *Trie[T])AppendSliceTrans(vs []T) *Trie[T] {
	for _,v := range vs {
		t = t.Append(t.id, v)
	}
	return t
}

// ReadSlice only returns the slice pointing to the underlying array in the trie
// Use this to implement io reads
func (t *Trie[T])ReadSlice(index int) ([]T,error) {
	if index > t.Size() {
		return nil, io.EOF
	} else if t.height == 0 {
		return t.content[index:], nil
	} else {
		i := index>>(b*t.height)
		for (index >= t.subsize[i]){
			i++
		}
		subtrie_starts := 0
		if i > 0 {
			subtrie_starts = t.subsize[i-1]
		}
		return t.subtrie[i].ReadSlice(index-subtrie_starts)
	}
}

/*
func (t *Trie[T])Take(index int) *Trie[T] {
}

func (t *Trie[T])Drop(index int) *Trie[T] {
}

func (t *Trie[T])Concat(t *Trie[T]) *Trie[T] {
}
*/


