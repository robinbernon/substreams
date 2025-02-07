package store

func (s *BaseStore) Append(ord uint64, key string, value []byte) {
	var newVal []byte
	oldVal, found := s.GetAt(ord, key)
	if !found {
		newVal = make([]byte, len(value))
		copy(newVal[0:], value)
	} else {
		newVal = make([]byte, len(oldVal)+len(value))
		copy(newVal[0:], oldVal)
		copy(newVal[len(oldVal):], value)
	}
	s.set(ord, key, newVal)
}
