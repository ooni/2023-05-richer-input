package undsl

import "github.com/ooni/2023-05-richer-input/pkg/unlang/uncompiler"

func templateName[T uncompiler.FuncTemplate]() string {
	return (*new(T)).TemplateName()
}
