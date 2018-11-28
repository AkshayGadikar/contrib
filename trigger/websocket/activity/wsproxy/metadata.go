package wsproxy

type Settings struct {
	WSconnection   string   `md:"wsconnection,required"`
	Uri            string   `md:"uri,required"`
	maxConnections int	`md:"maxconnections"`
}

