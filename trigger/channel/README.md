<!--
title: Channel
weight: 4706
-->
# Channel Trigger
This trigger provides your flogo application the ability to start an action via a named engine channel

## Installation

```bash
flogo install github.com/project-flogo/contrib/trigger/channel
```

## Schema
Settings, Outputs and Endpoint:

```json
{
  "settings": [
  ],
  "output": [
    {
      "name": "data",
      "type": "any"
    }
  ],
  "handler": {
    "settings": [
      {
        "name": "channel",
        "type": "string",
        "required" : true
      }
    ]
  }
}
```
## Settings      
### Handler:
| Setting     | Description    |
|:------------|:---------------|
| channel      | The internal engine channel |         


## Example Configurations

Triggers are configured via the triggers.json of your application. The following are some example configuration of the CHANNEL Trigger.

### Run Flow
Configure the Trigger to handle an event recieved on the 'test' channel

```json
{
  "triggers": [
    {
      "id": "flogo-channel",
      "ref": "github.com/project-flogo/contrib/trigger/channel",
      "handlers": [
        {
          "settings": {
            "channel": "test"
          },
          "action": {
            "ref": "github.com/TIBCOSoftware/flogo-contrib/action/flow",
            "settings": {
                "flowURI": "res://flow:testflow"
            }       
          }
        }
      ]
    }
  ]
}
```
