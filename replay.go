package dtls13

// replayWindow implements the sliding anti-replay window required by RFC 9147
// section 4.5.1. check does not mutate state; accept commits a sequence number.
type replayWindow struct {
	latest      uint64
	bitmap      uint64
	initialized bool
	size        uint
}

func newReplayWindow(size int) replayWindow { return replayWindow{size: uint(size)} }
func (w *replayWindow) nextExpected() uint64 {
	if !w.initialized {
		return 0
	}
	return w.latest + 1
}
func (w *replayWindow) check(seq uint64) bool {
	if !w.initialized || seq > w.latest {
		return true
	}
	delta := w.latest - seq
	return delta < uint64(w.size) && w.bitmap&(uint64(1)<<delta) == 0
}
func (w *replayWindow) accept(seq uint64) {
	if !w.initialized {
		w.initialized = true
		w.latest = seq
		w.bitmap = 1
		return
	}
	if seq > w.latest {
		delta := seq - w.latest
		if delta >= 64 {
			w.bitmap = 0
		} else {
			w.bitmap <<= delta
		}
		w.latest = seq
		w.bitmap |= 1
		return
	}
	delta := w.latest - seq
	if delta < 64 {
		w.bitmap |= uint64(1) << delta
	}
}
