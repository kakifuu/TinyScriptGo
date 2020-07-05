package util

import (
	"TinyScriptGo/lexer"
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
	ErrEndOfTokens = errors.New("End of tokens ")
)

type Iterator interface {
	Peek() (interface{}, error)
	PeekN(n int) ([]interface{}, error)
	Next() (interface{}, error)
	HasNext() bool
	PutBack(n int)
}

type MyIterator struct {
	srcNext    func() (interface{}, error)
	cache      []interface{}
	ptr        int
	putBackCnt int
	srcReadCnt int
}

func CreateMyIterator(srcNext func() (interface{}, error)) *MyIterator {
	return &MyIterator{
		srcNext:    srcNext,
		cache:      make([]interface{}, CACHE_SIZE),
		ptr:        0,
		putBackCnt: 0,
		srcReadCnt: 0,
	}
}

func (myIt *MyIterator) Peek() (interface{}, error) {
	element, err := myIt.Next()
	myIt.PutBack(1)
	return element, err
}

func (myIt *MyIterator) PeekN(n int) ([]interface{}, error) {
	if n > CACHE_SIZE {
		panic(ErrOverPeek)
	}
	res := make([]interface{}, n)
	for i := 0; i < n; i++ {
		element, err := myIt.Next()
		if err != nil {
			return nil, err
		}
		res[i] = element
	}
	myIt.PutBack(n)
	return res, nil
}

func (myIt *MyIterator) HasNext() bool {
	_, err := myIt.Peek()
	return err == nil
}

func (myIt *MyIterator) Next() (interface{}, error) {
	var element interface{}
	var err error
	if myIt.putBackCnt > 0 {
		element = myIt.cache[myIt.ptr]
		myIt.putBackCnt--
	} else {
		element, err = myIt.srcNext()
		myIt.cache[myIt.ptr] = element
		myIt.srcReadCnt++
	}
	myIt.ptr++
	if myIt.ptr == CACHE_SIZE {
		myIt.ptr = 0
	}
	return element, err
}

func (myIt *MyIterator) PutBack(n int) {
	if n > myIt.srcReadCnt || n > CACHE_SIZE {
		panic(ErrOverPutBack)
	}
	myIt.putBackCnt += n
	myIt.ptr = (myIt.ptr - n + CACHE_SIZE) % CACHE_SIZE
}

func NewRuneIterator(reader *bufio.Reader) *MyIterator {
	return CreateMyIterator(func() (interface{}, error) {
		return nextRune(reader)
	})
}

func nextRune(srcReader *bufio.Reader) (string, error) {
	r, _, err := srcReader.ReadRune()
	return string(r), err
}

func NewTokenIterator(tokens []*lexer.Token) *MyIterator {
	return CreateMyIterator(nextTokenGenerator(tokens))
}

func nextTokenGenerator(tokens []*lexer.Token) func() (interface{}, error) {
	i := 0
	return func() (interface{}, error) {
		if i >= len(tokens) {
			return nil, ErrEndOfTokens
		}
		token := tokens[i]
		i++
		return token, nil
	}
}
