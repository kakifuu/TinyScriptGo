package util

import (
	"bufio"
	"strings"
	"testing"
)

func compareSlices(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := 0; i < len(s1); i++ {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

func TestIterator_Next(t *testing.T) {
	reader := bufio.NewReader(strings.NewReader("ABCDEFGHIJKL"))
	it := NewRuneIterator(reader)
	wants := []string{
		"A", "B", "C",
		"D", "E", "F",
		"G", "H", "I",
		"J", "K", "L",
		string(0),
	}
	for _, want := range wants {
		element, _ := it.Next()
		got := element.(string)
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	}
}

func TestIterator_Peek(t *testing.T) {
	reader := bufio.NewReader(strings.NewReader("ABCDE"))
	it := NewRuneIterator(reader)
	wants := []string{"A", "A", "A", "A", "A"}
	for _, want := range wants {
		element, _ := it.Peek()
		got := element.(string)
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	}
	it.Next()
	wants = []string{"B", "B", "B", "B", "B"}
	for _, want := range wants {
		element, _ := it.Peek()
		got := element.(string)
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	}
}

func TestIterator_PeekN(t *testing.T) {
	reader := bufio.NewReader(strings.NewReader("ABCDE"))
	it := NewRuneIterator(reader)
	tests := []struct {
		n    int
		want []string
	}{
		{1, []string{"A"}},
		{2, []string{"A", "B"}},
		{3, []string{"A", "B", "C"}},
	}
	for _, test := range tests {
		elements, _ := it.PeekN(test.n)
		got := MapToString(elements)
		if !compareSlices(got, test.want) {
			t.Errorf("got %v, want %v", got, test.want)
		}
	}
	// test panic
	func() {
		defer func() {
			r := recover()
			err, ok := r.(error)
			if !ok || err != ErrOverPeek {
				t.Errorf("got %v, want %v", err, ErrOverPeek)
			}
		}()
		_, _ = it.PeekN(11)
	}()
}

func TestIterator_PutBack(t *testing.T) {
	reader := bufio.NewReader(strings.NewReader("ABC"))
	it := NewRuneIterator(reader)
	func() {
		defer func() {
			r := recover()
			err, ok := r.(error)
			if !ok || err != ErrOverPutBack {
				t.Errorf("got %v, want %v", err, ErrOverPutBack)
			}
		}()
		it.PutBack(1)
	}()
	tests := []struct {
		action func() string
		want   string
	}{
		{
			func() string {
				it.Next()
				it.PutBack(1)
				element, _ := it.Next()
				return element.(string)
			},
			"A",
		},
		{
			func() string {
				it.Next()
				it.Next()
				it.PutBack(2)
				element, _ := it.Next()
				return element.(string)
			},
			"B",
		},
	}
	for _, test := range tests {
		got := test.action()
		if got != test.want {
			t.Errorf("got %v, want %v", got, test.want)
		}
	}
}

func TestIterator_HasNext(t *testing.T) {
	reader := bufio.NewReader(strings.NewReader("ABC"))
	it := NewRuneIterator(reader)
	tests := []struct {
		preAction func()
		want      bool
	}{
		{
			func() {
				it.Next()
			},
			true,
		},
		{
			func() {
				it.Next()
			},
			true,
		},
		{
			func() {
				it.Next()
			},
			false,
		},
	}
	for _, test := range tests {
		test.preAction()
		got := it.HasNext()
		if got != test.want {
			t.Errorf("got %v, want %v", got, test.want)
		}
	}
}
