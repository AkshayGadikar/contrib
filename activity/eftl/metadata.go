package eftl

type Settings struct {}

type Input struct {
	URL            string `md:"url"`
	Id             string `md:"id"`
	User           string `md:"user"`
	Password       string `md:"password"`
	CA             string `md:"ca"`
	Dest           string `md:"dest"`
	Content        string `md:"content"`
}


/*type Output struct {
	PathParams  map[string]string `md:"pathParams"`
	QueryParams map[string]string `md:"queryParams"`
	Params      map[string]string `md:"params"`
	Content     interface{}       `md:"content"`
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"pathParams":  o.PathParams,
		"queryParams": o.QueryParams,
		"params":      o.Params,
		"content":     o.Content,
	}
}

func (o *Output) FromMap(values map[string]interface{}) error {
	o.PathParams = values["pathParams"].(map[string]string)
	o.QueryParams = values["queryParams"].(map[string]string)
	o.Params = values["params"].(map[string]string)
	o.Content = values["content"].(interface{})
	return nil
}*/