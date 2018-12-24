package main

import "fmt"

type Property struct {
	Name string
	Type string
	Example string
	Required bool
}

func(p Property) String() string {
	return fmt.Sprintf("{\n\tPropertyName: %s\n\tPropertyType: %s\n\tExample: %s\n\tRequired: %v\n}",
		p.Name, p.Type, p.Example, p.Required)
}

type Component struct {
	Name string
	Type string
	Properties []Property
}

func (c Component) String() string {
	return fmt.Sprintf("{\n\tComponentName: %s\n\tComponentType: %s\n\tProperties: %v\n}",
		c.Name, c.Type, c.Properties)
}