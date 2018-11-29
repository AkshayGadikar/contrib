package wsproxy

type Settings struct {
	WSconnection   interface{}   `md:"wsconnection,required"`
	Uri            string        `md:"uri,required"`
	maxConnections string	     `md:"maxconnections"`
}

type Input struct {
}

type Output struct {
}