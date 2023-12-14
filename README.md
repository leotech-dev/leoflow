# Leotech Numaflow plugins

This repository contains a set of Numaflow plugins, including custom sources and sinks, map and reduce UDFs.

## Plugin list

- Map
  - [jq](#jq)
<!-- - Reduce
  - batcher
- Sink
  - Salesforce
  - Clickhouse
  - Snowflake -->

## Build

Clone the repository and build it with

`docker build . -t leoflow`

## Configuration

Each plugin defines its own configuration options that can be specified using environment variables, as outlined in Numaflow documentation.

## Map UDFs

### jq

The plugin allows you to run [jq](https://jqlang.github.io/jq/) expressions on the input JSON data.

> jq is like sed for JSON data - you can use it to slice and filter and map and transform JSON with the same ease that sed, awk, grep and friends let you play with text.

This UDF applies a specified jq expression to the input data and sends the result in an output message.

If the expression produces multiple result sets, each result is published as a separate message.

#### Usage

Command for the UDF container:

```sh
/leoflow map jq
```

Example pipeline config:

```yaml
apiVersion: numaflow.numaproj.io/v1alpha1
kind: Pipeline
metadata:
  name: my-pipeline
spec:
  vertices:
    # ...
    - name: jq
      udf:
        container:
          image: leoflow@latest
          command: [ '/leoflow' ]
          args: [ 'map', 'jq' ]
          env:
            - name: JQ_DEBUG
              value: 'true' # booleans must be strings in env var declaration
            - name: JQ_MODE
              value: 'map' # can be "map" or "tag", see below
            - name: JQ_EXPRESSION
              value: |
                { 
                  "external_id": (.Data.id | tostring), 
                  "new_field": "This is a new field"
                }

# ...

```

#### Configuration

| Name | Type | Values | Description |
|-|-|-|-|
| `JQ_DEBUG` | string | `true`, `false` | Debug mode |
| `JQ_MODE` | string | `map`, `tag`| Operation mode (see description below). Default is `map` |
| `JQ_EXPRESSION` | string | | jq expression to execute |
| `JQ_TIMEOUT` | string | `"5s", "1m"` | Maximum execution time. Default is `1s` |

#### Mode

This UDF supports two modes:

1. `map`: in this mode the UDF will execute the expression and send its result to the output. Just like regular `jq` behavior.
2. `tag`: in this mode, the expression must return a string or an array of strings which are then used to tag the input message before sending it to the output.


##### Examples:

###### Mode: `map`, single result set

Input message:

```json
{
    "field1": "value1",
    "field2": "value2"
}
```

Expression to update a field value:

```
.field1 = "new_value"
```

Output:

```json
{
    "field1": "new_value",
    "field2": "value2"
}
```

---

###### Mode: `map`, multiple result sets

Input message:

```json
[
    {
        "field": "value1"
    },
    {
        "field": "value2"
    }
]
```

Expression to split the input array into separate elements:

```
.[]
```

Output:

Message 1:
```json
{
    "field": "value1",
}
```

Message 2:
```json
{
    "field": "value2",
}
```

---

###### Mode: `tag`

Input message:

```json
{
    "data": 1,
}
```

Expression returning a tag or set of tags:

```
if .data % 2 == 0 then "even" else "odd" end
```

Output:

The output message will be tagged with "odd".

```json
{
    "data": 1,
}
```

---

#### jq extensions

The plugin extends jq with several handy functions.

##### `fetch(url; opts)`

> [!IMPORTANT]
> Note that function params in jq are separated with semicolon.

The function allows you to perform an HTTP(s) call to the specified URL and return the result in `jq` expressions.

It is useful for scenarios when you need to enrich the message with the data from the API, to make a decision based on the API response, and in many other cases.

Parameters:

| Name | Description |
|-|-|
| `url` | The full request URL |
| `opts` | Request options, see below. |

###### Request options

The `opts` parameter is a JSON object with the following fields:

| Name | Type | Description |
|-|-|-|
| `method` | string | HTTP verb. Examples: `GET`, `POST`, `DELETE`, etc. Default is `GET`.|
| `headers` | JSON object | Request headers. Keys are header names, values can be a string or an array of strings. If the value is an array of strings, multiple headers with the same name are added, one per value. |
| `body` | string | Request body. |

_**Note**: to send piped input as request body, do:_

```js
{
    "method": "POST",
    "body": ( . | tostring) 
}
```



###### Example

Adding a new field called "api_response" with the data returned by an API call to the input:

```
.api_response = fetch(
    "https://mydomain.example/internal-api";
    {
        "method": "POST",
        "headers": {
            "Authorization": "Bearer Njg1M2ViMGQtN2NlNC00MmIwLWExMzEtN2RhMDQ1NTQ3YjE2",
            
            "X-Custom-Header": [
                "header value 1",
                "header value 2"
            ]
        },

        "body": (. | tostring)
    }
)
```


###### Return value

The function returns the following JSON object:

```json
{
    "status_code": 200,
    "headers": {
        "X-Response-Header": [
            "header value 1",
            "header value 2",
        ],
    },
    "body": "response body",
}
```


###### Multiple result sets

An expression like this will split the input JSON array into separate elements.
```
.[]
```

In this case, each element returned by the expression will be passed to output as a separate message.
