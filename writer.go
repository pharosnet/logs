package logs

import (
	"io"
	"sync"
)

type ErrorHandle func(int64, error)

type Writer interface {
	Writer(Element) int64
}

// mutex-writer
type ByteBufferWriter struct {
	Out io.Writer
	errorHandle ErrorHandle
	bytePool *Pool
	lock *sync.Mutex
}

func NewByteBufferWriter(out io.Writer, errorHandle ErrorHandle) *ByteBufferWriter {
	w := new(ByteBufferWriter)
	w.Out = out
	w.errorHandle = errorHandle
	w.bytePool = new(Pool)
	w.lock = new(sync.Mutex)
	return w
}

func (w *ByteBufferWriter) Writer(e Element) int64  {
	w.lock.Lock()
	defer w.lock.Unlock()
	buf := w.bytePool.Get()
	defer w.bytePool.Put(buf)
	buf.Write(e.Bytes())
	buf.WriteByte('\n')
	n, err := buf.WriteTo(w.Out)
	if err != nil {
		if w.errorHandle != nil {
			w.errorHandle(n, err)
		}
		return n
	}
	return n
}

// async
type AsyncByteBufferWriter struct {
	Out io.Writer
	errorHandle ErrorHandle
	bytePool *Pool
	ch chan writeEvent
	wg *sync.WaitGroup
}

type writeEvent struct {
	e Element
	ch chan int64
}


func NewAsyncByteBufferWriter(out io.Writer, errorHandle ErrorHandle, bufSize int) *AsyncByteBufferWriter {
	w := new(AsyncByteBufferWriter)
	w.Out = out
	w.errorHandle = errorHandle
	w.bytePool = new(Pool)
	if bufSize <= 0 {
		bufSize = 64
	}
	w.ch = make(chan writeEvent, bufSize)
	w.wg = new(sync.WaitGroup)
	w.listenCh()
	return w
}

func (w *AsyncByteBufferWriter) listenCh() {
	go func(w *AsyncByteBufferWriter) {
		for {
			event, ok := <- w.ch
			if !ok {
				break
			}
			buf := w.bytePool.Get()
			buf.Write(event.e.Bytes())
			buf.WriteByte('\n')
			n, err := buf.WriteTo(w.Out)
			if err != nil {
				if w.errorHandle != nil {
					w.errorHandle(n, err)
				}
			}
			w.bytePool.Put(buf)
			event.ch <- n
		}
	}(w)
}

func (w *AsyncByteBufferWriter) Writer(e Element) int64 {
	w.wg.Add(1)
	defer w.wg.Done()
	result := make(chan int64, 1)
	event := writeEvent{e:e, ch:result}
	w.ch <- event
	n := <- result
	close(result)
	return n
}

func (w *AsyncByteBufferWriter) Flush() {
	w.wg.Wait()
}
