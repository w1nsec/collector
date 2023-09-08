package memstorage

import "fmt"

type Storage interface {
	UpdateCounters(name string, value int64)
	UpdateGauges(name string, value float64)
	String() string
}

type MemStorage struct {
	dataCounters map[string]int64
	dataGauges   map[string]float64
}

func (ms *MemStorage) String() string {
	//TODO implement me
	//panic("implement me")
	var s1 = "[counters]\n"
	var s2 = "[gauges]\n"
	count := 0
	numInLine := 5
	for key, val := range ms.dataCounters {
		s1 += fmt.Sprintf("%s:%d ", key, val)
		count++
		if count == numInLine {
			count = 0
			s1 += "\n"
		}
	}
	s1 += "\n"
	count = 0
	for key, val := range ms.dataGauges {
		s2 += fmt.Sprintf("%s:%f ", key, val)
		count++
		if count == numInLine {
			count = 0
			s2 += "\n"
		}
	}
	s2 += "\n\n"

	return s1 + s2

}

func (ms *MemStorage) UpdateCounters(name string, value int64) {
	ms.dataCounters[name] += value
}

func (ms *MemStorage) UpdateGauges(name string, value float64) {
	ms.dataGauges[name] = value
}

func NewMemStorage() *MemStorage {
	ms := new(MemStorage)
	ms.dataCounters = make(map[string]int64)
	ms.dataGauges = make(map[string]float64)
	return ms
}
