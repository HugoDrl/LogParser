package parser

func filterLog(log *Log, filters *ParseSettings) bool {
	if !filters.StartDate.IsZero() && log.Time.Compare(filters.StartDate) < 0 {
		return false
	}
	if !filters.EndDate.IsZero() && log.Time.Compare(filters.EndDate) > 0 {
		return false
	}
	if filters.Level != "" && log.Level != filters.Level {
		return false
	}
	if filters.Service != "" && log.Service != filters.Service {
		return false
	}
	return true
}
