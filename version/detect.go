package version

import "reflect"

func DetectBestMatch[V any, T Entity[V]](entity T) V {
	var bestVersion V
	base := entity.GetBaseStruct()
	baseValue := reflect.ValueOf(base)
	baseType := baseValue.Type()

	highestScore := 0

	for _, v := range entity.GetVersionStructs() {
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
			bestVersion = v.GetVersion()
		}
	}

	return bestVersion
}
