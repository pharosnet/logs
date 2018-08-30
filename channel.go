package logs

import (
	"context"
	"fmt"
	"github.com/pharosnet/flyline"
	"os"
	"reflect"
)

const (
	flybuffer_cap = 1024 * 16
)

type Channel interface {
	Send(pac *Packet) error
	Recv() (*Packet, bool)
	Close(ctx context.Context) error
}

func NewFlyChannel(sink Sink) Channel {
	ch := &flyChannel{
		buffer:flyline.NewArrayBuffer(flybuffer_cap),
		sink:sink,
	}
	go func(ch *flyChannel) {
		for {
			pac, ok := ch.Recv()
			if !ok {
				break
			}
			ch.sink.FlowInto(pac)
		}
	}(ch)
	return ch
}

type flyChannel struct {
	buffer flyline.Buffer
	sink Sink
}

func (c *flyChannel) Send(pac *Packet) error {
	c.buffer.Send(pac)
	return nil
}

func (c *flyChannel) Recv() (*Packet, bool) {
	v, ok := c.buffer.Recv()
	if !ok {
		return nil, ok
	}
	pac, isPac := v.(*Packet)
	if !isPac {
		panic(fmt.Sprintf("the type of received from channel is not *Packet, received's type is %s", reflect.TypeOf(v)))
		os.Exit(1)
	}
	return pac, ok
}

func (c *flyChannel) Close(ctx context.Context) error {
	bufCloseErr := c.buffer.Close()
	bufSyncErr := c.buffer.Sync(ctx)
	sinkCloseErr := c.sink.Close(ctx)
	if bufCloseErr != nil {
		return bufCloseErr
	}
	if bufSyncErr != nil {
		return bufSyncErr
	}
	if sinkCloseErr != nil {
		return sinkCloseErr
	}
	return nil
}
