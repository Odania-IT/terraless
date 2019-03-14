package templates

import (
	"github.com/Odania-IT/terraless/schema"
	"github.com/sirupsen/logrus"
)

func consolidateEventData(terralessData schema.TerralessData) map[string]schema.FunctionEvents {
	logrus.Debugf("Consolidating function data")
	functionEvents := map[string]schema.FunctionEvents{}
	for functionName, functionConfig := range terralessData.Config.Functions {
		for _, event := range functionConfig.Events {
			events := functionEvents[functionConfig.Type].Events
			functionEvent := schema.FunctionEvent{
				FunctionName:  functionName,
				FunctionEvent: event,
			}

			if events == nil {
				functionEvents[functionConfig.Type] = schema.FunctionEvents{
					Events: map[string][]schema.FunctionEvent{},
				}
			}

			functionEvents[functionConfig.Type].Events[event.Type] = append(functionEvents[functionConfig.Type].Events[event.Type], functionEvent)
		}
	}

	return functionEvents
}
