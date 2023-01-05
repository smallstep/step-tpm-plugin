package storage

type AK struct {
	Name string
	Data []byte
}

type Key struct {
	Name string
	Data []byte
	// TODO: add properties to identify the AK that attested this key (if it was attested)? Created?
}
