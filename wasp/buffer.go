package wasp

// SliceBuffer keeps Capacity of type T, after len => cap overrides old data
type SliceBuffer[T any] struct {
	Idx      int
	Capacity int
	Data     []T
}

// NewSliceBuffer creates new limited capacity slice
func NewSliceBuffer[T any](cap int) *SliceBuffer[T] {
	return &SliceBuffer[T]{Capacity: cap, Data: make([]T, 0)}
}

// Append appends T if len <= cap, overrides old data otherwise
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
