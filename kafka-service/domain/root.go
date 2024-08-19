package domain

type AggregateRoot struct {
	version int
}

func (a *AggregateRoot) Version() int {
	if a == nil {
		return 0
	}
	return a.version
}

func (a *AggregateRoot) SetVersion() {
	if a == nil {
		a.version = 1
		return
	}
	a.version++
}

func NewAggregateRoot() *AggregateRoot {
	return &AggregateRoot{
		version: 1,
	}
}
