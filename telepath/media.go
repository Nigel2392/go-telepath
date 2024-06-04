package telepath

import "html/template"

type nullMedia struct{}

func (m *nullMedia) Merge(other Media) Media {
	return other
}

func (m *nullMedia) JS() []template.HTML {
	return []template.HTML{}
}

func (m *nullMedia) CSS() []template.HTML {
	return []template.HTML{}
}
