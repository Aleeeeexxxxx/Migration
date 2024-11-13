package producer

type FixedBuffer struct {
	buf  []byte
	size int

	start int
	next  int
}

func NewBuffer(size int) *FixedBuffer {
	return &FixedBuffer{
		buf:   make([]byte, size),
		size:  size,
		next:  0,
		start: 0,
	}
}

func (buf *FixedBuffer) Buffer() []byte {
	copy(buf.buf, buf.Raw())

	buf.next = buf.next - buf.start
	buf.start = 0

	return buf.buf[buf.next:]
}

func (buf *FixedBuffer) Peek(n int) {
	buf.next += n
}

func (buf *FixedBuffer) Reduce(n int) {
	buf.start += n
}

func (buf *FixedBuffer) Raw() []byte {
	return buf.buf[buf.start:buf.next]
}

func (buf *FixedBuffer) Empty() bool {
	return buf.next == buf.start
}
