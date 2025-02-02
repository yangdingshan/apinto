package response_rewrite

import (
	"github.com/eolinker/apinto/utils"
	"github.com/eolinker/eosc"
	"reflect"
)

type Driver struct {
	profession string
	name       string
	label      string
	desc       string
	configType reflect.Type
}

func (d *Driver) Check(v interface{}, workers map[eosc.RequireId]eosc.IWorker) error {
	_, err := d.check(v)
	if err != nil {
		return err
	}
	return nil
}

func (d *Driver) check(v interface{}) (*Config, error) {
	conf, ok := v.(*Config)
	if !ok {
		return nil, eosc.ErrorConfigFieldUnknown
	}

	err := conf.doCheck()
	if err != nil {
		return nil, err
	}

	return conf, nil
}

func (d *Driver) ConfigType() reflect.Type {
	return d.configType
}

func (d *Driver) Create(id, name string, v interface{}, workers map[eosc.RequireId]eosc.IWorker) (eosc.IWorker, error) {
	conf, err := d.check(v)
	if err != nil {
		return nil, err
	}

	//若body非空且需要base64转码
	if conf.Body != "" && conf.BodyBase64 {
		conf.Body, err = utils.B64Decode(conf.Body)
		if err != nil {
			return nil, err
		}
	}

	r := &ResponseRewrite{
		Driver:     d,
		id:         id,
		statusCode: conf.StatusCode,
		body:       conf.Body,
		headers:    conf.Headers,
		match:      conf.Match,
	}

	return r, nil
}
