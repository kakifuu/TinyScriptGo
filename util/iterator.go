package util

import (
	"bufio"
	"errors"
)

const (
	CACHE_SIZE = 10
	END_SYM    = string(0)
)

var (
	ErrOverPeek    = errors.New("Peeked too much ")
	ErrOverPutBack = errors.New("Put back too much ")
)

type Iterator struct {
	srcReader  *bufio.Reader
	cache      []string
	ptr        int
	putBackCnt int
	srcReadCnt int
}

func NewIterator(reader *bufio.Reader) *Iterator {
	return &Iterator{
		srcReader:  reader,
		cache:      make([]string, CACHE_SIZE),
		ptr:        0,
		putBackCnt: 0,
		srcReadCnt: 0,
	}
}

func srcNext(srcReader *bufio.Reader) string {
	r, _, err := srcReader.ReadRune()
	if err != nil {
		return END_SYM
	}
	return string(r)
}

func (it *Iterator) Peek() string {
	s := it.Next()
	it.PutBack(1)
	return s
}

func (it *Iterator) PeekN(n int) []string {
	if n > CACHE_SIZE {
		panic(ErrOverPeek)
	}
	res := make([]string, n)
	for i := 0; i < n; i++ {
		res[i] = it.Next()
	}
	it.PutBack(n)
	return res
}

func (it *Iterator) Next() string {
	var s string
	if it.putBackCnt > 0 {
		s = it.cache[it.ptr]
		it.putBackCnt--
	} else {
		s = srcNext(it.srcReader)
		it.cache[it.ptr] = s
		it.srcReadCnt++
	}
	it.ptr++
	if it.ptr == CACHE_SIZE {
		it.ptr = 0
	}
	return s
}

func (it *Iterator) PutBack(n int) {
	if n > it.srcReadCnt || n > CACHE_SIZE {
		panic(ErrOverPutBack)
	}
	it.putBackCnt += n
	it.ptr = (it.ptr - n + CACHE_SIZE) % CACHE_SIZE
}

func (it *Iterator) HasNext() bool {
	lastPtr := it.ptr - 1
	if lastPtr < 0 {
		lastPtr = CACHE_SIZE - 1
	}
	return it.cache[lastPtr] != END_SYM
}
