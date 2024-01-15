package crypto

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"
)

func CheckCertificate(listCertCA, listCertAutosigned string) func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
	// Créer un pool de certificats et ajouter le certificat de confiance
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM([]byte(listCertCA))

	// Load certificate autosigned
	var blocks []byte
	rest := []byte(listCertAutosigned)
	for {
		var block *pem.Block
		block, rest = pem.Decode(rest)
		if block == nil {
			break
		}
		blocks = append(blocks, block.Bytes...)
		if len(rest) == 0 {
			break
		}
	}

	certificatsAutosigned, err := x509.ParseCertificates(blocks)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}

	return checkCertificate(caCertPool, certificatsAutosigned)
}
func checkCertificate(caCertPool *x509.CertPool, certificatsAutosigned []*x509.Certificate) func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
	return func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {

		// Check exist
		if l := len(rawCerts); l == 0 {
			return fmt.Errorf("got len(rawCerts) = %d, wanted >0", l)
		}

		// Check standard FROM GOLANG ( handshakeclient: verifyServerCertificate)
		opts := x509.VerifyOptions{
			Roots:       caCertPool,
			CurrentTime: time.Now(),
			//DNSName:       c.config.ServerName,
			Intermediates: x509.NewCertPool(),
		}

		for _, certRaw := range rawCerts[1:] {
			cert, err := x509.ParseCertificate(certRaw)
			if err == nil {
				opts.Intermediates.AddCert(cert)
			}
		}

		// Check certificate
		cert, err := x509.ParseCertificate(rawCerts[0])
		if err != nil {
			return fmt.Errorf("Server certificat is not here.")
		}
		opts.Intermediates.AddCert(cert)
		_, err = cert.Verify(opts)
		if err == nil {
			// Check standard is OK
			return nil
		}

		// Check if certificate autoself without IP
		for _, certificatAutosigned := range certificatsAutosigned {
			if certificatAutosigned.Equal(cert) {
				return nil
			}
		}

		return fmt.Errorf("Error x509: not trust certificate and not autosigned")
	}
}
