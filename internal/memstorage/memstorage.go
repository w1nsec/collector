package memstorage

type MemStorage struct {
	dataCounters map[string]int64
	dataGauges   map[string]float64
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
