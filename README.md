# CRD validation

The code here is just a wrapper around [apiextensions-apiserver(crd validation)](https://github.com/kubernetes/apiextensions-apiserver/blob/master/pkg/apis/apiextensions/validation/validation.go). We are handling the v1beta1 version of CRDs.

## How to use

```
go run main.go --def crddefinition.json
```

Or

```
go build .
./validation --def crddefinition.json
```

## Testing

Try adding invalid fields based on [(crd validation)](https://github.com/kubernetes/apiextensions-apiserver/blob/master/pkg/apis/apiextensions/validation/validation.go).

Like:
```
        "openAPIV3Schema": {
          "description": "PetstoreAPI is the Schema for the PetstoreAPIs API",
          "additionalProperties": "this will give error",
```

Ouput: ` detect the invalid fields`
```
spec.validation.openAPIV3Schema.additionalProperties: Forbidden: must not be used at the root
spec.validation.openAPIV3Schema.type: Required value: must not be empty at the root
``

## Testing with craft

Firstly build the operator. Then, convert `config/crd/bases/*.yaml` in created operator to json and validate that json.
```
go build .
./validation --def <jsonpath>
```