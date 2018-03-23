package service

import (
	"github.com/containerum/chkit/pkg/model"
	yaml "gopkg.in/yaml.v2"
)

var (
	_ model.YAMLrenderer = Service{}
	_ yaml.Marshaler     = Service{}
)

func (serv Service) RenderYAML() (string, error) {
	data, err := yaml.Marshal(serv)
	return string(data), err
}

func (serv Service) MarshalYAML() (interface{}, error) {
	return serv.origin, nil
}
