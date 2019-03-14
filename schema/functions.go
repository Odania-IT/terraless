package schema

type FunctionEvent struct {
	FunctionName  string
	FunctionEvent TerralessFunctionEvent
}

type FunctionEvents struct {
	Events map[string][]FunctionEvent
}
