package schema

import (
	"github.com/ONSdigital/go-ns/avro"
)

// TODO: remove or replace hello called structure and model with app specific
var helloCalledEvent = `{
  "type": "record",
  "name": "hello-called",
  "fields": [
    {"name": "recipient_name", "type": "string", "default": ""}
  ]
}`

// HelloCalledEvent is the Avro schema for Hello Called messages.
var HelloCalledEvent = &avro.Schema{
	Definition: helloCalledEvent,
}
