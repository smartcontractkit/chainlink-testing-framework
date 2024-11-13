package wasp

// SliceBuffer keeps Capacity of type T, after len => cap overrides old data
type SliceBuffer[T any] struct {
	Idx      int
	Capacity int
	Data     []T
}

// NewSliceBuffer creates a new SliceBuffer with the specified capacity.
// It initializes the buffer to hold elements of any type T, starting with an empty slice.
// The function returns a pointer to the newly created SliceBuffer.
func NewSliceBuffer[T any](cap int) *SliceBuffer[T] {
	return &SliceBuffer[T]{Capacity: cap, Data: make([]T, 0)}
}

// Append adds an element of type T to the SliceBuffer. If the buffer's capacity is reached,
// it overwrites the oldest element, maintaining a circular buffer. The index is incremented
// after each append operation, wrapping around to the start if necessary.
func (m *SliceBuffer[T]) Append(s T) {
	if m.Idx >= m.Capacity {
		m.Idx = 0
	}
	if len(m.Data) <= m.Capacity {
		m.Data = append(m.Data, s)
	} else {
		m.Data[m.Idx] = s
	}
	m.Idx++
}
