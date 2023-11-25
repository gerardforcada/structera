package detector

import (
	"github.com/gerardforcada/structera/interfaces"
	"reflect"
)

func BestMatchingEra[T interfaces.Hub](hub T) (bestEra int) {
	base := hub.GetBaseStruct()
	baseValue := reflect.ValueOf(base)
	baseType := baseValue.Type()

	highestScore := 0

	for _, v := range hub.GetVersionStructs() {
		score := 0
		candidateValue := reflect.ValueOf(v)
		candidateType := candidateValue.Type()

		for i := 0; i < baseType.NumField(); i++ {
			baseField := baseType.Field(i)
			baseFieldVal := baseValue.Field(i)
			baseFieldType := baseField.Type

			if baseFieldType.Kind() == reflect.Ptr && baseFieldVal.IsNil() {
				continue
			}
			if baseFieldType.Kind() == reflect.Ptr {
				baseFieldType = baseFieldType.Elem()
			}

			for j := 0; j < candidateType.NumField(); j++ {
				candidateField := candidateType.Field(j)
				if baseField.Name == candidateField.Name && baseFieldType == candidateField.Type {
					score++
					break
				}
			}
		}

		if score > highestScore {
			highestScore = score
			bestEra = v.GetVersion()
		}
	}

	return bestEra
}
