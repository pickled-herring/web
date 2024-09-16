package web

import (
	"math/rand"
	"testing"
	"io"
)
/* What properties are there with the trie?
   Persistence
   It's also a vector so a read should return the write in the correct order
   Reads at the same index should give the same data
   What else...

   Whats the spec of our library? All that tests do is check that our code does
   not fundamentally change the functionality of things

 */

// TestAppend(
// Reads are followed what is written or appended

var letters = []byte("qwertyuiopasdfghjklzxcvbnm" + 
	"QWERTYUIOPASDFGHJKLZXCVBNM" +
	"1234567890,.?! ")

func randSeq(n int) []byte {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return b
}

func TestAppendTrans(t *testing.T){
	num := 1024
	ref := randSeq(num)
	a := TrieFromSlice[byte](ref)
	if a.Size() != num {
		t.Fatalf("a.Size() = %d != %d", a.Size(), num)
	}
	if !a.Full() {
		t.Fatalf("a.Full() is false!")
	}
	num2 := 500
	ref2 := randSeq(num2)
	a = a.AppendSliceTrans(ref2)
	if a.Size() != num + num2 {
		t.Fatalf("a.Size() = %d != %d", a.Size(), num + num2)
	}
}

func TestAppend(t *testing.T){
	num := 1024
	ref := randSeq(num)
	a := TrieFromSlice[byte](ref)
	if a.Size() != num {
		t.Fatalf("a.Size() = %d != %d", a.Size(), num)
	}

	num2 := 503
	ref2 := randSeq(num2)
	b := a.AppendSlice(ref2)
	if b == a {
		t.Fatalf("a: %p, b: %p\n", a, b)
	}
	if a.Size() != num {
		t.Fatalf("a.Size() = %d != %d", a.Size(), num)
	}
	if b.Size() != num + num2 {
		t.Logf("%v\n", b.subsize)
		t.Fatalf("b.Size() = %d != %d", b.Size(), num+num2)
	}
}

func TestAppendRead(t *testing.T){
	num := 1024
	ref := randSeq(num)
	a := TrieFromSlice[byte](ref)
	for i := 0; i < 990; i++ {
		r,err := a.ReadSlice(i)
		if err != nil {
			t.Fatalf("a.ReadSlice(%d) returned err %v", i, err)
		}
		if r[0] != ref[i] {
			t.Fatalf("a.ReadSlice(%d) = %s != %s", i, r, ref[i:i+32])
		}
	}

	_, err := a.ReadSlice(len(ref)+1)
	if err != io.EOF {
		t.Fatalf("a.ReadSlice(%d) returned err %v", len(ref)+1, err)
	}
}

func TestSize(t *testing.T){
	a := NewTrie[byte](0)
	if a.Size() != 0 {
		t.Fatalf("a.Size() = %d != 0", a.Size())
	}
}

func TestTake(t *testing.T){
	num := 1000
	ref := randSeq(num)
	a := TrieFromSlice[byte](ref)
	for i := 2; i < num-1; i++ {
		b := a.Take(0,i)
		if b.Size() != i {
			t.Fatalf("a.Take(%d).Size() = %d, not %d", i, b.Size(), i)
		}

		readb, err := b.ReadSlice(i-2)
		if err != nil {
			t.Fatalf("a.Take(%d).ReadSlice(%d) returned err %v", 
				i, i-2, err)
		}
		if readb[0] != ref[i-2] {
			t.Fatalf("a.Take(%d).ReadSlice(%d) returned %s != %s", 
				i, i-2, readb, ref[i-2:i])
		}

		readb, err = b.ReadSlice(0)
		if err != nil {
			t.Fatalf("a.Take(%d).ReadSlice(%d) returned err %v", 
				i, 0, err)
		}
		if readb[0] != ref[0] {
			t.Fatalf("a.Take(%d).ReadSlice(%d) returned %s != %s", 
				i, 0, readb, ref[:5])
		}
	}
}

func TestDrop(t *testing.T) {
	num := 1000
	ref := randSeq(num)
	a := TrieFromSlice[byte](ref)
	for i := 2; i < num-1; i++ {
		b := a.Drop(0,i)
		if b.Size() != num-i {
			t.Fatalf("a.Take(%d).Size() = %d, not %d", i, b.Size(), num-i)
		}

		readb, err := b.ReadSlice(0)
		if err != nil {
			t.Fatalf("a.Take(%d).ReadSlice(%d) returned err %v", 
				i, 0, err)
		}
		if readb[0] != ref[i] {
			t.Fatalf("a.Take(%d).ReadSlice(%d) returned %s != %s", 
				i, 0, readb, ref[i:i+5])
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

func BenchmarkTake(b *testing.B){
	num := 10000
	a := TrieFromSlice[byte](randSeq(num))
	for i:=0 ; i<b.N; i++ {
		a.Take(0, rand.Intn(num))
	}
}
