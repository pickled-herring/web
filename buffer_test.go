package web

import (
	"math/rand"
	"testing"
)
/* What properties are there with the trie?
   Persistence
   It's also a vector so a read should return the write in the correct order
   Reads at the same index should give the same data
   What else...

   Whats the spec of our library? All that tests do is check that our code does
   not fundamentally change the functionality of things

 */

func TestSize(t *testing.T) {
	a := NewTrie[int](0)
	if a.Size() != 0 {
		t.Fatalf("a.Size() = %d != 0", a.Size())
	}
	b := a.AppendSlice([]int{1,2,3})
	if a.Size() != 0 {
		t.Fatalf("a.Size() = %d != 0", a.Size())
	}
	if b.Size() != 3 {
		t.Fatalf("b.Size() = %d != 3", b.Size())
	}
}

func TestAppendFullTrie(t *testing.T) {
	a := NewTrans[int](0,1234)
	for i := 0; i < 1024; i ++ {
		a = a.Append(1234, rand.Int())
	}
	if a.Size() != 1024{
		t.Errorf("a.Size() = %d != 1024", a.Size())
	}
	if a.height != 1 {
		t.Errorf("a.height = %d != 1", a.height)
	}
	a = a.Append(1234, 12)
}

func TestAppend(t *testing.T){
	a := NewTrans[byte](0,1234)
	s := []byte("This things what else is there to know")
	for i := 0; i < 100; i++ {
		a = a.AppendSliceTrans(s)
	}
	if a.Size() != len(s) * 100 {
		t.Errorf("a.Size() = %d != len(s)*100 = %d", a.Size(), len(s)*100)
	}
	for i := 0; i< 1000; i++ {
		index := rand.Intn(len(s)*100-1)
		read_slice,err := a.ReadSlice(index)
		if err != nil {
			t.Errorf("a.ReadSlice(%d) failed with error %v", index, err)
			return
		}
		if read_slice[0] != s[index%len(s)] {
			t.Errorf("a.ReadSlice(%d)[0] == %c in %s; expected s[%d], %c",
				index,read_slice[0], read_slice, index%len(s), s[index%len(s)])
			return
		}
	}
}

func BenchmarkAppendTrans(b *testing.B){
	a := NewTrans[byte](0,1234)
	s := []byte("This things what else is there to know")
	for i:=0 ; i<b.N; i++ {
		a = a.Append(1234,s[i%len(s)])
	}
}

func BenchmarkAppendPersis(b *testing.B){
	a := NewTrie[byte](0)
	s := []byte("This things what else is there to know")
	for i:=0 ; i<b.N; i++ {
		a = a.Append(0,s[i%len(s)])
	}
}

func BenchmarkAppendSlice(b *testing.B){
	a := NewTrie[byte](0)
	s := []byte("This things what else is there to know")
//	b.Logf("len(s): %d", len(s))
	for i:=0 ; i<b.N; i++ {
		a = a.AppendSlice(s)
	}
}

func BenchmarkAppendSliceTrans(b *testing.B){
	a := NewTrans[byte](0,1234)
	s := []byte("This things what else is there to know")
	for i:=0 ; i<b.N; i++ {
		a = a.AppendSliceTrans(s)
	}
}

func BenchmarkReadSlice(b *testing.B){
	a := NewTrie[byte](0)
	s := []byte("This things what else is there to know")
	for i := 0; i < 1000; i++ {
		a = a.AppendSlice(s)
	}
	for i:=0 ; i<b.N; i++ {
		a.ReadSlice(i%(1000*len(s)))
	}
}
