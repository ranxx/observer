package observer

import (
	"fmt"
	"reflect"
)

func checkObserver(observer interface{}, fieldName string) (reflect.Type, string, reflect.Value) {
	t, topic, fc := checkObserverForInterface(observer)

	if len(topic) > 0 && !fc.IsZero() {
		return t, topic, fc
	}

	// 字段名
	field, ok := t.FieldByName(fieldName)
	if !ok {
		panic(fmt.Sprintf("%s is no %s field", t.String(), fieldName))
	}

	if len(topic) <= 0 {
		topic, ok = field.Tag.Lookup("topic")
	}

	if !ok {
		panic(fmt.Sprintf("%s.%s field is no topic in the tag", t.String(), fieldName))
	}

	if !fc.IsZero() {
		return t, topic, fc
	}

	function, ok := field.Tag.Lookup("notice")
	if !ok {
		panic(fmt.Sprintf("%s.%s field is no notice in the tag", t.String(), fieldName))
	}

	fc = reflect.ValueOf(observer).MethodByName(function)
	if fc.IsZero() {
		panic(fmt.Sprintf("%s does not exist in %s", function, t.String()))
	}

	return t, topic, fc
}

func checkObserverForInterface(observer interface{}) (reflect.Type, string, reflect.Value) {
	t := reflect.TypeOf(observer)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		panic(fmt.Sprintf("%s is not reflect.Struct", t.String()))
	}

	var topic, function string

	topicf, ok := observer.(Topic)
	if ok {
		topic = topicf.Topic()
	}

	functionf, ok := observer.(Function)
	if ok {
		function = functionf.Function()
	}

	if len(function) <= 0 {
		return t, topic, reflect.Zero(reflect.TypeOf(observer))
	}

	fc := reflect.ValueOf(observer).MethodByName(function)
	if fc.IsZero() {
		panic(fmt.Sprintf("%s does not exist in %s", function, t.String()))
	}

	return t, topic, fc
}

func checkFunc(fc interface{}) (reflect.Type, reflect.Value) {
	t := reflect.TypeOf(fc)
	if t.Kind() != reflect.Func {
		panic(fmt.Sprintf("%s is not reflect.Func", t.String()))
	}
	return t, reflect.ValueOf(fc)
}
