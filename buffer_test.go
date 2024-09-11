package web

import (
	"math/rand"
	"testing"
)
// TODO: Write tests... preferably property based testing
/* What properties are there with the trie?
   Persistence
   It's also a vector so a read should return the write in the correct order
   Reads at the same index should give the same data
   What else...

   Whats the spec of our library? All that tests do is check that our code does
   not fundamentally change the functionality of things
 */

func TestAppend(t *testing.T){
	a := NewTrie[byte](0)
	s := []byte("This things what else is there to know")
	for i := 0; i < 100; i++ {
		a = a.AppendSlice(s)
	}
	for i := 0; i< 1000; i++ {
		index := rand.Intn(len(s)*100-1)
		read_slice,err := a.ReadSlice(index)
		if err != nil {
			t.Errorf("a.ReadSlice(%d) failed with error %v", index, err)
		}
		if read_slice[0] != s[index%len(s)] {
			t.Errorf("a.ReadSlice(%d)[0] == %c in %s; expected s[%d], %c",
				index,read_slice[0], read_slice, index%len(s), s[index%len(s)])
		}
	}
}

func BenchmarkAppend(b *testing.B){
	a := NewTrans[byte](0,1234)
	s := []byte("This things what else is there to know")
	for i:=0 ; i<b.N; i++ {
		a = a.Append(s[i%len(s)])
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
