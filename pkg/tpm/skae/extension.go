package skae

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"math/big"

	"github.com/google/go-attestation/attest"
)

var oidSubjectKeyAttestationEvidence = asn1.ObjectIdentifier{2, 23, 133, 6, 1, 1} // SKAE (Subject Key Attestation Evidence) OID: 2.23.133.6.1.1

// func makeObjectIdentifier(oid []int) (e encoder, err error) {
// 	if len(oid) < 2 || oid[0] > 2 || (oid[0] < 2 && oid[1] >= 40) {
// 		return nil, StructuralError{"invalid object identifier"}
// 	}

// 	return oidEncoder(oid), nil
// }

// [2 5 4 3]
// [1 3 6 1 5 5 7 48 1]
// [1 3 6 1 5 5 7 48 2]
// [] TODO: this empty OID triggers an invalid OID error. Haven't found what makes it empty yet, though.

var oidAuthorityInfoAccessOcsp = asn1.ObjectIdentifier{1, 3, 6, 1, 5, 5, 7, 48, 1}
var oidAuthorityInfoAccessIssuers = asn1.ObjectIdentifier{1, 3, 6, 1, 5, 5, 7, 48, 2}

func CreateSubjectKeyAttestationEvidenceExtension(akCert *x509.Certificate, params attest.CertificationParameters) (pkix.Extension, error) {

	// asn1Issuer, err := asn1.Marshal(pkix.Name{CommonName: "Test Issuer"}.ToRDNSequence())
	// if err != nil {
	// 	return pkix.Extension{}, fmt.Errorf("error marshaling issuer: %w", err)
	// }

	// skaeExtension := asn1SKAE{
	// 	TCGSpecVersion: asn1TCGSpecVersion{Major: 2, Minor: 0},
	// 	KeyAttestationEvidence: asn1KeyAttestationEvidence{ // TODO: this requires a choice to be encoded; normal vs encrypted
	// 		AttestEvidence: asn1AttestationEvidence{
	// 			TPMCertifyInfo: asn1TPMCertifyInfo{
	// 				CertifyInfo: asn1.BitString{ // TODO: check if setting the values like this is correct
	// 					Bytes:     params.CreateAttestation,
	// 					BitLength: len(params.CreateAttestation) * 8,
	// 				},
	// 				Signature: asn1.BitString{
	// 					Bytes:     params.CreateSignature,
	// 					BitLength: len(params.CreateSignature) * 8,
	// 				},
	// 			},
	// 			TPMIdentityCredAccessInfo: asn1TPMIdentityCredentialAccessInfo{ // TODO: this should contain information from the AIK/AK cert. See x509.go on how to handle these values.
	// 				AuthorityInfoAccess: []asn1AuthorityInfoAccessSyntax{
	// 					{
	// 						Method:   oidAuthorityInfoAccessOcsp,
	// 						Location: asn1.RawValue{Tag: asn1.TagOID, Class: asn1.ClassContextSpecific, Bytes: []byte("https://www.example.com/ocsp/cert1")},
	// 					},
	// 					{
	// 						Method:   oidAuthorityInfoAccessIssuers,
	// 						Location: asn1.RawValue{Tag: asn1.TagOID, Class: asn1.ClassContextSpecific, Bytes: []byte("https://www.example.com/issuing/cert1")},
	// 					},
	// 				},
	// 				IssuerSerial: asn1IssuerSerial{ // TODO: should come from the AK cert
	// 					Issuer: asn1.RawValue{FullBytes: asn1Issuer},
	// 					Serial: big.NewInt(1337),
	// 				},
	// 			},
	// 		},
	// 	},
	// }
	// skaeExtensionBytes, err := asn1.Marshal(skaeExtension)
	// if err != nil {
	// 	return pkix.Extension{}, fmt.Errorf("creating SKAE extension failed: %w", err)
	// }

	skaeExtensionBytes := []byte{}

	return pkix.Extension{
		Id:       oidSubjectKeyAttestationEvidence,
		Critical: false, // non standard extension; let's not break clients
		Value:    skaeExtensionBytes,
	}, nil
}

type asn1SKAE struct {
	TCGSpecVersion         asn1TCGSpecVersion
	KeyAttestationEvidence asn1KeyAttestationEvidence
}

type asn1TCGSpecVersion struct {
	Major int
	Minor int
}

type asn1KeyAttestationEvidence struct {
	AttestEvidence          asn1AttestationEvidence // TODO: ASN1 CHOICE between those two
	EnvelopedAttestEvidence asn1EnvelopedAttestationEvidence
}

type asn1AttestationEvidence struct {
	TPMCertifyInfo            asn1TPMCertifyInfo
	TPMIdentityCredAccessInfo asn1TPMIdentityCredentialAccessInfo
}

type asn1TPMCertifyInfo struct {
	CertifyInfo asn1.BitString
	Signature   asn1.BitString
}

type asn1TPMIdentityCredentialAccessInfo struct {
	AuthorityInfoAccess []asn1AuthorityInfoAccessSyntax
	IssuerSerial        asn1IssuerSerial
}

type asn1AuthorityInfoAccessSyntax struct {
	Method   asn1.ObjectIdentifier
	Location asn1.RawValue
}

type asn1IssuerSerial struct {
	Issuer asn1.RawValue
	Serial *big.Int
}

type asn1EnvelopedAttestationEvidence struct {
	//RecipientInfos
	EncryptedAttestInfo asn1EncryptedAttestationInfo
}

type asn1EncryptedAttestationInfo struct {
	EncryptionAlgorithm     pkix.AlgorithmIdentifier
	EncryptedAttestEvidence []byte
}
