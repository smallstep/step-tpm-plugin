module github.com/smallstep/step-tpm-plugin

go 1.19

require (
	github.com/google/go-attestation v0.4.4-0.20220404204839-8820d49b18d9
	github.com/jedib0t/go-pretty v4.3.0+incompatible
	github.com/spf13/cobra v1.6.1
	github.com/spf13/pflag v1.0.5
	go.step.sm/crypto v0.0.0-00010101000000-000000000000
)

require (
	github.com/asaskevich/govalidator v0.0.0-20200907205600-7a23bdc65eef // indirect
	github.com/go-openapi/errors v0.20.2 // indirect
	github.com/go-openapi/strfmt v0.21.3 // indirect
	github.com/google/btree v1.0.1 // indirect
	github.com/google/certificate-transparency-go v1.1.2 // indirect
	github.com/google/go-tpm v0.3.3 // indirect
	github.com/google/go-tspi v0.2.1-0.20190423175329-115dea689aad // indirect
	github.com/inconshreveable/mousetrap v1.0.1 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/mitchellh/mapstructure v1.3.3 // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/peterbourgon/diskv/v3 v3.0.1 // indirect
	github.com/rogpeppe/go-internal v1.9.0 // indirect
	github.com/ryboe/q v1.0.19 // indirect
	github.com/schollz/jsonstore v1.1.0 // indirect
	go.mongodb.org/mongo-driver v1.11.1 // indirect
	golang.org/x/crypto v0.6.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
)

replace github.com/google/go-attestation v0.4.4-0.20220404204839-8820d49b18d9 => github.com/smallstep/go-attestation v0.4.4-0.20230113130042-0ad94dd6a52e

replace go.step.sm/crypto => ./../crypto
