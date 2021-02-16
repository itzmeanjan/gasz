package data

// WSTraffic - Keeping track of how many read & write message ops happened
// in lifetime of one single websocket connection
type WSTraffic struct {
	Read  uint64
	Write uint64
}
