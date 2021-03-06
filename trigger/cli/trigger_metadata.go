package cli

import (
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
)

var jsonMetadata = `{
  "name": "flogo-cli",
  "type": "flogo:trigger",
  "shim": "main",
  "ref": "github.com/TIBCOSoftware/flogo-contrib/trigger/cli",
  "version": "0.0.1",
  "title": "CLI Trigger",
  "description": "Simple CLI Trigger",
  "homepage": "https://github.com/TIBCOSoftware/flogo-contrib/tree/master/trigger/cli",
  "output": [
    {
      "name": "args",
      "type": "array"
    }
  ],
  "reply": [
    {
      "name": "data",
      "type": "any"
    }
  ],
  "handler": {
    "settings": [
      {
        "name": "command",
        "type": "string"
      },
      {
        "name": "default",
        "type": "boolean"
      }
    ]
  }
}
`

// init create & register trigger factory
func init() {
	md := trigger.NewMetadata(jsonMetadata)
	trigger.RegisterFactory(md.ID, NewFactory(md))
}
