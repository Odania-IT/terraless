package templates

import (
	"bytes"
	"github.com/Odania-IT/terraless/schema"
)

func processFunctions(terralessData *schema.TerralessData, buffer bytes.Buffer) bytes.Buffer {
	consolidatedFunctionEvents := consolidateEventData(*terralessData)

	for _, terralessProvider := range terralessData.TerralessProviders {
		for resourceType, functionEvents := range consolidatedFunctionEvents {
			buffer.WriteString(terralessProvider.RenderFunctionTemplates(resourceType, functionEvents, terralessData))
		}
	}

	return buffer
}
