package main

import (
	"testing"

	"github.com/linkedin/goavro/v2"
	"github.com/stretchr/testify/require"
)

var testCases = []struct {
	Payload string
}{
	{
		Payload: `{
			"id": "1r43ooyCVO683LCeSIABzZ0BV9Q",
			"title": "Goal",
			"description": "",
			"competitionId": 123,
			"tags": ["tag"],
			"duration": 20000.48,
			"additionalProps": {
				"prop": 100,
				"prop2": true,
				"prop3": 10.01,
				"prop4": "string"
			},
			"events": [
				{
					"id": "1r43ooyCVO683LCeSIABzZ0BV9Q",
					"type": "Goal",
					"timestamp": "2020-03-07T20:20:00",
					"gameTime": "58:49",
					"players": [
						{
							"id": 123,
							"name": "John Doe",
							"captain": true,
							"team": {
								"id": 123,
								"name": "Rock Stars",
								"abbr": "RCF"
							}
						}
					],
					"match": {
						"timestamp": "2020-03-07T19:00:00Z",
						"teams": {
							"home": {
								"id": 123,
								"name": "Rock Stars",
								"abbr": "RCF"
							},
							"away": {
								"id": 123,
								"name": "Pop Stars",
								"abbr": "PCF"
							}
						}
					}
				}
			]
		}`,
	},
}

func Test_JSON2AVRO_Parse(t *testing.T) {
	for _, tc := range testCases {
		res, err := Parse("test", []byte(tc.Payload), false)
		require.NoError(t, err)

		codec, err := goavro.NewCodec(string(res))
		require.NoError(t, err)

		// Convert textual Avro data (in Avro JSON format) to native Go form
		native, _, err := codec.NativeFromTextual([]byte(tc.Payload))
		require.NoError(t, err)

		// Convert native Go form to binary Avro data
		binary, err := codec.BinaryFromNative(nil, native)
		require.NoError(t, err)

		// Convert binary Avro data back to native Go form
		native, _, err = codec.NativeFromBinary(binary)
		require.NoError(t, err)

		// Convert native Go form to textual Avro data
		_, err = codec.TextualFromNative(nil, native)
		require.NoError(t, err)
	}
}
