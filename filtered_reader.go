package main

type filteredReader struct {
	reader  recordReader
	filters []filterFunc
}

func (rd *filteredReader) Next() *record {
	for {
		r := rd.reader.Next()
		if r == nil {
			return nil
		}

		for _, filter := range rd.filters {
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

func newFilteredReader(rd recordReader, filters ...filterFunc) *filteredReader {
	return &filteredReader{
		reader:  rd,
		filters: filters,
	}
}
