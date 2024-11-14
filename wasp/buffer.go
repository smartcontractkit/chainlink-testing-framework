package wasp

// SliceBuffer keeps Capacity of type T, after len => cap overrides old data
type SliceBuffer[T any] struct {
	Idx      int
	Capacity int
	Data     []T
}

// NewSliceBuffer creates and returns a new SliceBuffer for elements of type T with the specified capacity.
func NewSliceBuffer[T any](cap int) *SliceBuffer[T] {
	return &SliceBuffer[T]{Capacity: cap, Data: make([]T, 0)}
}

// Append adds the element s to the SliceBuffer. If the buffer has not reached its capacity, s is appended to the data slice. Once the capacity is exceeded, Append overwrites the oldest element in the buffer. The internal index is incremented to track the next insertion point.
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
