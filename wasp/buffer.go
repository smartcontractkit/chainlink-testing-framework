package wasp

// SliceBuffer keeps Capacity of type T, after len => cap overrides old data
type SliceBuffer[T any] struct {
	Idx      int
	Capacity int
	Data     []T
}

// NewSliceBuffer initializes a new SliceBuffer with the specified capacity. 
// It returns a pointer to the newly created SliceBuffer, which contains an empty slice of the specified type. 
// This function is useful for creating a buffer that can hold elements of a generic type T, allowing for dynamic data storage.
func NewSliceBuffer[T any](cap int) *SliceBuffer[T] {
	return &SliceBuffer[T]{Capacity: cap, Data: make([]T, 0)}
}

// Append adds a new element to the SliceBuffer. If the buffer has reached its capacity, 
// it will overwrite the oldest element. The function maintains the current index for 
// the next insertion, wrapping around when the end of the buffer is reached. 
// This allows for efficient storage of a fixed-size collection of elements.
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
