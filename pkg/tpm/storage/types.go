package storage

import "fmt"

type AK struct {
	Name string
	Data []byte
}

type Key struct {
	Name string
	Data []byte
	// TODO: add properties to identify the AK that attested this key (if it was attested)? Created?
}

const (
	akPrefix  = "ak-"
	keyPrefix = "key-"
)

type tpmObjectType string

const (
	typeAK  tpmObjectType = "AK"
	typeKey tpmObjectType = "KEY"
)

type serializedAK struct {
	Name string
	Type tpmObjectType
	Data []byte
}

type serializedKey struct {
	Name string
	Type tpmObjectType
	Data []byte
}

func keyForKey(name string) string {
	return fmt.Sprintf("%s%s", keyPrefix, name)
}

func keyForAK(name string) string {
	return fmt.Sprintf("%s%s", akPrefix, name)
}
