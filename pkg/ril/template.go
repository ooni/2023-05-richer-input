package ril

import "github.com/ooni/2023-05-richer-input/pkg/ric"

func templateName[T ric.FuncTemplate]() string {
	return (*new(T)).TemplateName()
}
