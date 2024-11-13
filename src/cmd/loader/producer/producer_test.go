package producer

import (
	"fmt"
	"io"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRows_ToValues(t *testing.T) {
	rq := require.New(t)
	rows := &Rows{
		"a,b,c",
		"a,b,c",
		"a,b,c",
	}
	rq.Equal("(a,b,c),(a,b,c),(a,b,c)", rows.ToValues())
}

func TestBuffer(t *testing.T) {
	rq := require.New(t)
	buf := NewBuffer(5)

	rq.Equal(0, len(buf.Raw()))
	rq.Equal(5, len(buf.Buffer()))
	rq.True(buf.Empty())

	buf.Peek(2)
	rq.Equal(2, len(buf.Raw()))
	rq.Equal(3, len(buf.Buffer()))
	rq.False(buf.Empty())

	buf.Peek(1)
	rq.Equal(3, len(buf.Raw()))
	rq.Equal(2, len(buf.Buffer()))
	rq.False(buf.Empty())

	buf.Reduce(2)
	rq.Equal(1, len(buf.Raw()))
	rq.Equal(4, len(buf.Buffer()))
	rq.False(buf.Empty())

	buf.Reduce(1)
	rq.Equal(0, len(buf.Raw()))
	rq.Equal(5, len(buf.Buffer()))
	rq.True(buf.Empty())
}

func TestProducer_Produce(t *testing.T) {
	rq := require.New(t)
	f, err := os.Open("test.csv")
	if err != nil {
		fmt.Printf("fail to open file: %v\n", err)
		return
	}
	defer func() { _ = f.Close() }()
	pool := &sync.Pool{
		New: func() interface{} {
			return make(Rows, 0, 5)
		},
	}

	p := NewProducer(f, 2, pool, 7)

	rows1, err := p.Produce()
	rq.NoError(err)
	rq.Equal(2, len(rows1))
	rq.Equal("1,1,1", rows1[0])
	rq.Equal("2,2,2", rows1[1])

	rows2, err := p.Produce()
	rq.EqualError(err, io.EOF.Error())
	rq.Equal(1, len(rows2))
	rq.Equal("3,3,3", rows2[0])
}
