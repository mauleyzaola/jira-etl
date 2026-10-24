package domains

type HistoryItems []HistoryItem

func (x HistoryItems) Filter(filter func(HistoryItem) bool) HistoryItems {
	var results HistoryItems
	for _, v := range x {
		if filter(v) {
			results = append(results, v)
		}
	}
	return results
}

func (x HistoryItems) Len() int { return len(x) }
