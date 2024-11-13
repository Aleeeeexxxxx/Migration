package producer

import (
	"fmt"
	"io"
	"os"
	"sync"
)

type Rows []string

func (data Rows) ToValues() string {
	values := ""
	for _, r := range data {
		values += fmt.Sprintf("(%s),", r)
	}
	return values[:len(values)-1] // remove last ","
}

type Producer struct {
	buf       *FixedBuffer
	pool      *sync.Pool
	batchSize int

	f         *os.File
	loadedAll bool
}

func NewProducer(f *os.File, batchSize int, pool *sync.Pool, bufSize int) *Producer {
	return &Producer{
		buf:       NewBuffer(bufSize),
		pool:      pool,
		batchSize: batchSize,
		f:         f,
		loadedAll: false,
	}
}

func (p *Producer) Produce() (Rows, error) {
	if p.loadedAll {
		return nil, io.EOF
	}

	x := p.pool.Get()
	rows := x.(Rows) // should be reset already

	for {
		rows = p.parse(rows)

		if len(rows) >= p.batchSize {
			if p.loadedAll && p.buf.Empty() {
				return rows, io.EOF
			}
			return rows, nil
		} else if p.loadedAll {
			return rows, io.EOF
		}

		if err := p.load(); err != nil {
			return nil, err
		}
	}

}

func (p *Producer) parse(rows Rows) Rows {
	next := 0
	data := string(p.buf.Raw())

LOOP:
	for i := 0; i < len(data); i++ {
		for ; i < len(data); i++ {
			if data[i] == '\n' {
				rows = append(rows, data[next:i])
				next = i + 1

				if len(rows) == p.batchSize {
					break LOOP
				}
			}
		}
	}

	if len(rows) != p.batchSize && p.loadedAll && next != len(data) {
		// last one may not end with '\n'
		rows = append(rows, data[next:])
		next = len(rows)
	}

	p.buf.Reduce(next)
	return rows
}

func (p *Producer) load() error {
	n, err := p.f.Read(p.buf.Buffer())
	if err != nil {
		if err == io.EOF {
			p.loadedAll = true
		} else {
			return err
		}
	}
	p.buf.Peek(n)
	return nil
}
