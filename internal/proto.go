package internal

import (
	"encoding/json"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v2"
)

func ProtosFromYaml(s []byte, m proto.Message) error {
	var obj interface{}
	if err := yaml.Unmarshal(s, &obj); err != nil {
		return err
	}

	obj = convert(obj)

	rawJSON, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	err = protojson.Unmarshal(rawJSON, m)
	if err != nil {
		return err
	}
	return nil
}

func convert(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = convert(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = convert(v)
		}
	}
	return i
}
