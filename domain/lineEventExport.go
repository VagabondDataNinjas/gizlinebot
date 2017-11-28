package domain

// LineEventExport: used for CSV events download
type LineEventExport struct {
	EventType   string
	UserId      string
	DisplayName string
	EventTime   string
}
