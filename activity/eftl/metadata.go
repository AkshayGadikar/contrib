package eftl

type Settings struct{}

type Input struct {
	URL      string `md:"url"`
	Id       string `md:"id"`
	User     string `md:"user"`
	Password string `md:"password"`
	CA       string `md:"ca"`
	Dest     string `md:"dest"`
	Content  map[string]interface{} `md:"content"`
}


func (o *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"url":o.URL,
		"id": o.Id,
		"user": o.User,
		"password": o.Password,
		"ca": o.CA,
		"dest": o.Dest,
		"content": o.Content,
	}
}

func (o *Input) FromMap(values map[string]interface{}) error {
	o.URL = values["url"].(string)
	o.Id = values["id"].(string)
	o.User = values["user"].(string)
	o.Password = values["password"].(string)
	o.CA = values["ca"].(string)
	o.Dest = values["dest"].(string)
	o.Content = values["content"].(map[string]interface{})
	return nil
}
