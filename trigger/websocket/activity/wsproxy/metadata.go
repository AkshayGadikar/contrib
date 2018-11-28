package wsproxy

type Settings struct {
	WSconnection   interface{}   `md:"wsconnection,required"`
	Uri            string        `md:"uri,required"`
	maxConnections int	     `md:"maxconnections"`
}

