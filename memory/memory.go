package memory

type Memory struct {
	words  []string
	index  int
	length int
}

func New() *Memory {
	return &Memory{
		words:  make([]string, 0),
		index:  -1,
		length: 0,
	}
}

func (m *Memory) Push(word string) {
	m.words = append(m.words, word)
	m.index = len(m.words) - 1
	m.length = len(m.words) + 1
}

func (m *Memory) Up() string {
	if m.length == 0 {
		return ""
	}
	m.index = (m.index - 1) % m.length
	return m.words[m.index]
}

func (m *Memory) Down() string {
	if m.length == 0 {
		return ""
	}
	m.index = (m.index + 1) % m.length
	return m.words[m.index]
}

func (m *Memory) Clear() {
	m.index = len(m.words) - 1
}
