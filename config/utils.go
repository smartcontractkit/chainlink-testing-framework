package config

import (
	"reflect"

	"github.com/pelletier/go-toml/v2"
	"github.com/rs/zerolog"
)

func BytesToAnyTomlStruct(logger zerolog.Logger, filename, configurationName string, target any, content []byte) error {
	// Create a new empty struct of the same type as the target
	newStruct := reflect.New(reflect.TypeOf(target).Elem()).Interface()

	// Unmarshal the content into the new struct
	err := toml.Unmarshal(content, newStruct)
	if err != nil {
		return err
	}

	// Unmarshal the content into the target struct
	err = toml.Unmarshal(content, target)
	if err != nil {
		return err
	}

	// Replace maps and slices in target with those from newStruct (since toml.Unmarshal does not replace maps and slices, but merges them instead)
	replaceMapsAndSlices(reflect.ValueOf(target).Elem(), reflect.ValueOf(newStruct).Elem())

	logger.Debug().Msgf("Successfully unmarshalled %s config file", filename)

	var someToml map[string]interface{}
	err = toml.Unmarshal(content, &someToml)
	if err != nil {
		return err
	}

	if configurationName == "" {
		logger.Debug().Msgf("No configuration name provided, will read only default configuration.")
		return nil
	}

	if _, ok := someToml[configurationName]; !ok {
		logger.Debug().Msgf("Config file %s does not contain configuration named '%s', will read only default configuration.", filename, configurationName)
		return nil
	}

	marshalled, err := toml.Marshal(someToml[configurationName])
	if err != nil {
		return err
	}

	newStruct = reflect.New(reflect.TypeOf(target).Elem()).Interface()
	err = toml.Unmarshal(marshalled, newStruct)
	if err != nil {
		return err
	}

	err = toml.Unmarshal(marshalled, target)
	if err != nil {
		return err
	}

	replaceMapsAndSlices(reflect.ValueOf(target).Elem(), reflect.ValueOf(newStruct).Elem())
	logger.Debug().Msgf("Configuration named '%s' read successfully.", configurationName)

	return nil
}

func replaceMapsAndSlices(target, newStruct reflect.Value) {
	target = reflect.Indirect(target)
	newStruct = reflect.Indirect(newStruct)

	for i := 0; i < target.NumField(); i++ {
		targetField := target.Field(i)
		newStructField := newStruct.Field(i)
		structField := target.Type().Field(i)

		// Check if the field is exported
		if structField.PkgPath != "" || !targetField.CanSet() {
			continue
		}

		switch targetField.Kind() {
		case reflect.Map, reflect.Slice:
			targetField.Set(newStructField)
		case reflect.Ptr:
			if newStructField.IsNil() {
				continue
			}
			if targetField.Elem().Kind() == reflect.Map || targetField.Elem().Kind() == reflect.Slice {
				if !newStructField.IsNil() {
					targetField.Set(newStructField)
				}
			} else if targetField.Elem().Kind() == reflect.Struct {
				replaceMapsAndSlices(targetField.Elem(), newStructField.Elem())
			}
		case reflect.Struct:
			replaceMapsAndSlices(targetField, newStructField)
		default:
			continue
		}
	}
}

//func overrideSlices(target any, newContent any) {
//	targetValue := reflect.ValueOf(target).Elem()
//	newContentValue := reflect.ValueOf(newContent).Elem()
//
//	traverseAndOverrideSlices(targetValue, newContentValue)
//}
//
//func traverseAndOverrideSlices(targetValue, newContentValue reflect.Value) {
//	for i := 0; i < targetValue.NumField(); i++ {
//		targetField := targetValue.Field(i)
//		newContentField := newContentValue.Field(i)
//		structField := targetValue.Type().Field(i)
//
//		fmt.Printf("targetField before: %v\n", targetField)
//
//		// Check if the field is exported
//		if structField.PkgPath != "" || !targetField.CanSet() {
//			fmt.Printf("targetField after: %v\n", targetValue.Field(i))
//			continue
//		}
//
//		if !newContentField.IsValid() {
//			fmt.Printf("targetField after: %v\n", targetValue.Field(i))
//			continue
//		}
//
//		switch targetField.Kind() {
//		case reflect.Struct:
//			traverseAndOverrideSlices(targetField, newContentField)
//		case reflect.Ptr:
//			if newContentField.IsNil() {
//				fmt.Printf("targetField after: %v\n", targetValue.Field(i))
//				continue // retain the original value
//			}
//			if targetField.Elem().Kind() == reflect.Struct {
//				traverseAndOverrideSlices(targetField.Elem(), newContentField.Elem())
//			} else if targetField.Elem().Kind() == reflect.Slice || targetField.Elem().Kind() == reflect.Map {
//				targetField.Elem().Set(newContentField.Elem())
//			}
//		case reflect.Slice, reflect.Map:
//			targetField.Set(newContentField)
//		default:
//			fmt.Printf("targetField after: %v\n", targetValue.Field(i))
//			continue // retain the original value
//		}
//
//		fmt.Printf("targetField after: %v\n", targetValue.Field(i))
//	}
//}

//func traverseAndOverrideSlices(targetValue, newContentValue reflect.Value) {
//	for i := 0; i < targetValue.NumField(); i++ {
//		targetField := targetValue.Field(i)
//		newContentField := newContentValue.Field(i)
//
//		// Skip invalid fields (fields that are not present in newContent)
//		if !newContentField.IsValid() || (newContentField.Kind() == reflect.Ptr && newContentField.IsNil()) {
//			continue
//		}
//
//		switch targetField.Kind() {
//		case reflect.Struct:
//			traverseAndOverrideSlices(targetField, newContentField)
//		case reflect.Ptr:
//			if newContentField.IsNil() {
//				continue // retain the original value
//			}
//			if targetField.IsNil() {
//				targetField.Set(reflect.New(targetField.Type().Elem()))
//			}
//			if targetField.Elem().Kind() == reflect.Struct {
//				traverseAndOverrideSlices(targetField.Elem(), newContentField.Elem())
//			} else if targetField.Elem().Kind() == reflect.Slice || targetField.Elem().Kind() == reflect.Ptr {
//				targetField.Elem().Set(newContentField.Elem())
//			} else {
//				targetField.Set(newContentField)
//			}
//		case reflect.Slice:
//			targetField.Set(newContentField)
//		default:
//			if targetField.CanSet() {
//				targetField.Set(newContentField)
//			}
//		}
//	}
//}
