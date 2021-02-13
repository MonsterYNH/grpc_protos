package route

// Info route info
type Info struct {
	ServiceName  string
	Method       string
	Pattern      string
	RequestType  string
	ResponseType string
}

var routes []Info

// RegisterInfo regist info
func RegisterInfo(info Info) {
	routes = append(routes, info)
}

// GetInfos get infos
func GetInfos() []Info {
	return routes
}
