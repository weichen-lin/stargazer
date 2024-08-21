package domain

type AggregateRoot struct {
	version int64
}

func (a *AggregateRoot) Version() int64 {
	if a == nil {
		return 0
	}
	return a.version
}

func (a *AggregateRoot) UpdateVersion() {
	a.version++
}

func NewAggregateRoot() *AggregateRoot {
	return &AggregateRoot{
		version: 1,
	}
}
