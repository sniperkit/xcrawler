package scraper

var TargetTypes = []string{"title", "desc", "image", "price", "stock", "count", "url", "tag", "extra", "cat"}
var TransportTypes = []string{"http", "https", "grpc", "tcp", "udp", "udp", "udp", "inproc", "ipc", "tlstcp", "ws", "wss"}
var MethodTypes = []string{"GET", "POST"}
var SelectorEngines = []string{"css", "xpath", "json", "xml", "csv"}
