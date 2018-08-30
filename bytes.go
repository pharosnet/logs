package logs

import (
	"io"
	"sort"
	"sync"
	"sync/atomic"
)

const (
	minBitSize = 6
	steps      = 20
	minSize = 1 << minBitSize
	calibrateCallsThreshold = 42000
	maxPercentile           = 0.95
)

var defaultPool Pool

func getBuffer() *ByteBuffer { return defaultPool.Get() }

func putBuffer(b *ByteBuffer) { defaultPool.Put(b) }

type ByteBuffer struct {
	B []byte
}

func (b *ByteBuffer) Len() int {
	return len(b.B)
}

func (b *ByteBuffer) ReadFrom(r io.Reader) (int64, error) {
	p := b.B
	nStart := int64(len(p))
	nMax := int64(cap(p))
	n := nStart
	if nMax == 0 {
		nMax = 64
		p = make([]byte, nMax)
	} else {
		p = p[:nMax]
	}
	for {
		if n == nMax {
			nMax *= 2
			bNew := make([]byte, nMax)
			copy(bNew, p)
			p = bNew
		}
		nn, err := r.Read(p[n:])
		n += int64(nn)
		if err != nil {
			b.B = p[:n]
			n -= nStart
			if err == io.EOF {
				return n, nil
			}
			return n, err
		}
	}
}

func (b *ByteBuffer) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(b.B)
	return int64(n), err
}

func (b *ByteBuffer) Bytes() []byte {
	return b.B
}

func (b *ByteBuffer) Write(p []byte) (int, error) {
	b.B = append(b.B, p...)
	return len(p), nil
}

func (b *ByteBuffer) WriteByte(c byte) error {
	b.B = append(b.B, c)
	return nil
}

func (b *ByteBuffer) WriteString(s string) (int, error) {
	b.B = append(b.B, s...)
	return len(s), nil
}

func (b *ByteBuffer) Set(p []byte) {
	b.B = append(b.B[:0], p...)
}

func (b *ByteBuffer) SetString(s string) {
	b.B = append(b.B[:0], s...)
}

func (b *ByteBuffer) String() string {
	return string(b.B)
}

func (b *ByteBuffer) Reset() {
	b.B = b.B[:0]
}

type Pool struct {
	calls       [steps]uint64
	calibrating uint64

	defaultSize uint64
	maxSize     uint64

	pool sync.Pool
}

func (p *Pool) Get() *ByteBuffer {
	v := p.pool.Get()
	if v != nil {
		return v.(*ByteBuffer)
	}
	return &ByteBuffer{
		B: make([]byte, 0, atomic.LoadUint64(&p.defaultSize)),
	}
}

func (p *Pool) Put(b *ByteBuffer) {
	idx := index(len(b.B))

	if atomic.AddUint64(&p.calls[idx], 1) > calibrateCallsThreshold {
		p.calibrate()
	}

	maxSize := int(atomic.LoadUint64(&p.maxSize))
	if maxSize == 0 || cap(b.B) <= maxSize {
		b.Reset()
		p.pool.Put(b)
	}
}

func (p *Pool) calibrate() {
	if !atomic.CompareAndSwapUint64(&p.calibrating, 0, 1) {
		return
	}

	a := make(callSizes, 0, steps)
	var callsSum uint64
	for i := uint64(0); i < steps; i++ {
		calls := atomic.SwapUint64(&p.calls[i], 0)
		callsSum += calls
		a = append(a, callSize{
			calls: calls,
			size:  minSize << i,
		})
	}
	sort.Sort(a)

	defaultSize := a[0].size
	maxSize := defaultSize

	maxSum := uint64(float64(callsSum) * maxPercentile)
	callsSum = 0
	for i := 0; i < steps; i++ {
		if callsSum > maxSum {
			break
		}
		callsSum += a[i].calls
		size := a[i].size
		if size > maxSize {
			maxSize = size
		}
	}

	atomic.StoreUint64(&p.defaultSize, defaultSize)
	atomic.StoreUint64(&p.maxSize, maxSize)

	atomic.StoreUint64(&p.calibrating, 0)
}

type callSize struct {
	calls uint64
	size  uint64
}

type callSizes []callSize

func (ci callSizes) Len() int {
	return len(ci)
}

func (ci callSizes) Less(i, j int) bool {
	return ci[i].calls > ci[j].calls
}

func (ci callSizes) Swap(i, j int) {
	ci[i], ci[j] = ci[j], ci[i]
}

func index(n int) int {
	n--
	n >>= minBitSize
	idx := 0
	for n > 0 {
		n >>= 1
		idx++
	}
	if idx >= steps {
		idx = steps - 1
	}
	return idx
}
