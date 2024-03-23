package conversor

import (
	"encoding/json"
	"fmt"
	"github.com/gerardforcada/structera/interfaces"
	"reflect"
)

func ToEra(target any, hub interfaces.Hub) error {
	targetVal := reflect.ValueOf(target)
	if targetVal.Kind() != reflect.Ptr || targetVal.IsNil() {
		return fmt.Errorf("target must be a non-nil pointer")
	}

	hubJSON, err := json.Marshal(hub)
	if err != nil {
		return fmt.Errorf("error marshaling hub: %v", err)
	}

	if _, ok := target.(*interfaces.Era); ok {
		eraType := reflect.TypeOf(target).Elem()
		if eraType.Kind() == reflect.Interface {
			eraType = reflect.ValueOf(target).Elem().Elem().Type()
		}

		newInstance := reflect.New(eraType).Interface()
		err = json.Unmarshal(hubJSON, newInstance)
		if err != nil {
			return fmt.Errorf("error unmarshaling into era: %v", err)
		}

		reflect.ValueOf(target).Elem().Set(reflect.ValueOf(newInstance).Elem())
		return nil
	}

	err = json.Unmarshal(hubJSON, target)
	if err != nil {
		return fmt.Errorf("error unmarshaling into target: %v", err)
	}

	return nil
}
