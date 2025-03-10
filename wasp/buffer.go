package wasp

// SliceBuffer keeps Capacity of type T, after len => cap overrides old data
type SliceBuffer[T any] struct {
	Idx      int
	Capacity int
	Data     []T
}

// NewSliceBuffer creates a new SliceBuffer with the specified capacity.
// It provides an efficient way to store and manage a fixed number of elements,
// enabling optimized access and manipulation in concurrent and decentralized applications.
func NewSliceBuffer[T any](capacity int) *SliceBuffer[T] {
	return &SliceBuffer[T]{Capacity: capacity, Data: make([]T, 0)}
}

// Append adds an element to the SliceBuffer. When the buffer reaches its capacity, it overwrites the oldest item.
// This function is useful for maintaining a fixed-size, circular collection of elements.
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
