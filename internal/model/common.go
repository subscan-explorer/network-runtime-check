package model

type Metadata struct {
	Name   string          `json:"name"`
	Prefix string          `json:"prefix"`
	Events []MetadataEvent `json:"events"`
	Calls  []MetadataCalls `json:"calls"`
}

type MetadataCalls struct {
	Lookup string `json:"lookup"`
	Name   string `json:"name"`
	//Docs   []string                     `json:"docs"`
	Args []MetadataModuleCallArgument `json:"args"`
}

type MetadataModuleCallArgument struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	TypeName string `json:"type_name"`
}

type MetadataEvent struct {
	Lookup       string   `json:"lookup"`
	Name         string   `json:"name"`
	Args         []string `json:"args"`
	ArgsTypeName []string `json:"args_type_name"`
}

type NetworkData[T any] struct {
	Network string
	Data    []T
	Err     error
}
