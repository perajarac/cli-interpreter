package memory

type Memory struct {
	words []string
	index int
}

func New() *Memory {
	return &Memory{
		words: make([]string, 0),
		index: -1,
	}
}

func (m *Memory) Push(word string) {
	m.words = append(m.words, word)
	m.index = len(m.words) - 1
}

func (m *Memory) Up() string {
	if len(m.words) == 0 {
		return ""
	}
	m.index = (m.index - 1 + len(m.words)) % len(m.words)
	return m.words[m.index]
}

func (m *Memory) Down() string {
	if len(m.words) == 0 {
		return ""
	}
	m.index = (m.index + 1) % len(m.words)
	return m.words[m.index]
}

func (m *Memory) Clear() {
	m.index = len(m.words) - 1
}
