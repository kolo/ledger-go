package ledger

type FilteredIterator struct {
	iter    RecordIterator
	filters []RecordFilter
}

func NewFilteredIterator(iter RecordIterator, filters ...RecordFilter) *FilteredIterator {
	return &FilteredIterator{
		iter:    iter,
		filters: filters,
	}
}

func (iter FilteredIterator) Next() *Record {
	for {
		r := iter.iter.Next()
		if r == nil {
			return nil
		}

		for _, filter := range iter.filters {
			r = filter(r)
			if r == nil {
				break
			}
		}

		if r == nil {
			continue
		}

		return r
	}
}
