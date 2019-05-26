package disruptor

import "runtime"

type Writer struct {
	written  *Cursor // highest committed sequence number
	upstream Barrier
	capacity int64 // number of values that can be stored
	previous int64 // highest sequence number handed out by Reserve
	gate     int64 // highest sequence number read by all consumers
}

/*
Think of the writable slots in the buffer as a sliding window of size capacity on the
infinite list of sequence numbers.

gate is the number just before the first value in the window. All consumers have read
values up to and including gate.

The highest value we can give out is gate + capacity.


0   1   2   3   4   5   6   7   8   9   10
            _____________
gate ---^
highest reservable -----^

*/

func NewWriter(written *Cursor, upstream Barrier, capacity int64) *Writer {
	assertPowerOfTwo(capacity)

	return &Writer{
		upstream: upstream,
		written:  written,
		capacity: capacity,
		previous: InitialSequenceValue,
		gate:     InitialSequenceValue,
	}
}

func assertPowerOfTwo(value int64) {
	if value > 0 && (value&(value-1)) != 0 {
		// Wikipedia entry: http://bit.ly/1krhaSB
		panic("The ring capacity must be a power of two, e.g. 2, 4, 8, 16, 32, 64, etc.")
	}
}

// Reserve allocates count sequence numbers to the caller.
// It returns the highest allocated value.
// E.g. after
//   s := w.Reserve(3)
// the caller can use sequence numbers s-2, s-1 and s.
//
// If count is greater than the capacity of the buffer, this will loop forever.
func (this *Writer) Reserve(count int64) int64 {
	this.previous += count

	// Wait until the highest reserved sequence number to be returned (previous) has
	// a place in the buffer (whose size is capacity). gate is the highest consumed
	// sequence number.
	for spin := int64(0); this.previous-this.capacity > this.gate; spin++ {
		// Occasionally call the scheduler to release the CPU. If SpinMask is 2^n-1,
		// this will happen every 2^n iterations.
		if spin&SpinMask == 0 {
			runtime.Gosched() // LockSupport.parkNanos(1L); http://bit.ly/1xiDINZ
		}

		this.gate = this.upstream.Read(0)
	}

	return this.previous
}

func (this *Writer) Await(next int64) {
	for next-this.capacity > this.gate {
		this.gate = this.upstream.Read(0)
	}
}

const SpinMask = 1024*16 - 1 // arbitrary; we'll want to experiment with different values
