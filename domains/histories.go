package domains

type Histories []History

func (x Histories) Filter(filter func(History) bool) Histories {
	var histories Histories
	for _, v := range x {
		if filter(v) {
			histories = append(histories, v)
		}
	}
	return histories
}

func (x Histories) Len() int { return len(x) }
