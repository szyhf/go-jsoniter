package jsoniter

import (
	"io"
)

// IteratorPool a thread safe pool of iterators with same configuration
type IteratorPool interface {
	BorrowIterator(data []byte) *Iterator
	ReturnIterator(iter *Iterator)
}

// StreamPool a thread safe pool of streams with same configuration
type StreamPool interface {
	BorrowStream(writer io.Writer) *Stream
	ReturnStream(stream *Stream)
}

func (cfg *frozenConfig) BorrowStream(writer io.Writer) *Stream {
	stream := cfg.streamPool.Get().(*Stream)
	stream.Reset(writer)
	return stream
}

func (cfg *frozenConfig) ReturnStream(stream *Stream) {
	stream.out = nil
	stream.Error = nil
	stream.Attachment = nil
	// 参考 fmt 包中关于 buf 的使用
	// fmt/print.go

	//
	// Proper usage of a sync.Pool requires each entry to have approximately
	// the same memory cost. To obtain this property when the stored type
	// contains a variably-sized buffer, we add a hard limit on the maximum
	// buffer to place back in the pool. If the buffer is larger than the
	// limit, we drop the buffer and recycle just the printer.
	//
	// See https://golang.org/issue/23199.

	// 设置一个强限制避免对象池持有缓冲区过大的对象
	if cap(stream.buf) > 64*1024 {
		return
	}
	cfg.streamPool.Put(stream)
}

func (cfg *frozenConfig) BorrowIterator(data []byte) *Iterator {
	iter := cfg.iteratorPool.Get().(*Iterator)
	iter.ResetBytes(data)
	return iter
}

func (cfg *frozenConfig) ReturnIterator(iter *Iterator) {
	iter.Error = nil
	iter.Attachment = nil
	cfg.iteratorPool.Put(iter)
}
