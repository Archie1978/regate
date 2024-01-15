package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestX509(t *testing.T) {

	certsSystem, err := x509.SystemCertPool()
	if err != nil {
		t.Fatal(err)
	}

	// Get certificate www.google.com
	conn, err := tls.Dial("tcp", "www.google.com:443", &tls.Config{
		VerifyPeerCertificate: checkCertificate(certsSystem, nil),
		InsecureSkipVerify:    true,
	})
	if err != nil {
		t.Fatal("Test Error: ", err)
	}
	conn.Close()
}

func TestX509FAiled(t *testing.T) {

	hosts := []string{
		"expired.badssl.com:443",
		"self-signed.badssl.com:443",
		"untrusted-root.badssl.com:443",
	}
	//    https:///
	//    https://untrusted-root.badssl.com/
	errors := []string{
		"Error x509: not trust certificate and not autosigned",
		"Error x509: not trust certificate and not autosigned",
		"Error x509: not trust certificate and not autosigned",
	}

	certsSystem, err := x509.SystemCertPool()
	if err != nil {
		t.Fatal(err)
	}

	// Dialer for change timeout
	dialer := &net.Dialer{
		Timeout:   5 * time.Minute, // Too Too long
		DualStack: false,           // Activer le support IPv6 si disponible
	}

	// Get certificate www.google.com
	for i, host := range hosts {
		conn, err := tls.DialWithDialer(dialer, "tcp", host, &tls.Config{
			VerifyPeerCertificate: checkCertificate(certsSystem, nil),
			InsecureSkipVerify:    true,
		})
		if err != nil {
			if err.Error() != errors[i] {
				t.Fatal("host:", host, "  ", err)
			}
		} else {
			t.Fatal("host:", host, "  ", err)
			conn.Close()
		}
	}
}

func TestX509AutoCert(t *testing.T) {

	// Générer une clé privée RSA
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal("erreur lors de la génération de la clé privée : ", err)
	}

	// Créer un certificat auto-signé
	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		t.Fatal("erreur lors de la génération du numéro de série : ", err)
	}

	cert := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"My Organization"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour), // Valide pendant un an
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// Créer le certificat auto-signé
	certDER, err := x509.CreateCertificate(rand.Reader, &cert, &cert, &privateKey.PublicKey, privateKey)
	if err != nil {
		t.Fatal("erreur lors de la création du certificat auto-signé :", err)
	}

	// Encapsuler la clé privée dans un bloc PEM
	privateKeyPEM := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}

	// Start serveur localhost
	certServer, err := tls.X509KeyPair(
		pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER}),
		pem.EncodeToMemory(privateKeyPEM),
	)

	if err != nil {
		t.Fatal(err)
	}
	cfg := &tls.Config{Certificates: []tls.Certificate{certServer}}
	srv := &http.Server{
		Addr:         "127.0.0.1:4329",
		TLSConfig:    cfg,
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
	}
	go func() {
		t.Fatal(srv.ListenAndServeTLS("", ""))
	}()

	// Wait serveur Start
	<-time.After(time.Second)

	// Check connection
	certsSystem, err := x509.SystemCertPool()
	if err != nil {
		t.Fatal(err)
	}

	// check verification fail
	fmt.Printf("1) Test certificate autosign ( connection failed )\n")
	conn, err := tls.Dial("tcp", "127.0.0.1:4329", &tls.Config{
		VerifyPeerCertificate: checkCertificate(certsSystem, nil),
		InsecureSkipVerify:    true,
	})
	if err == nil {
		conn.Close()
		t.Fatal("Test Error: ", err)
	}
	<-time.After(time.Second)

	// check verification OK
	fmt.Printf("2) Test certificate autosign ( Good Connexion )\n")
	certificatePEM := string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER}))
	conn, err = tls.Dial("tcp", "127.0.0.1:4329", &tls.Config{
		VerifyPeerCertificate: CheckCertificate("", certificatePEM+"\n"),
		InsecureSkipVerify:    true,
	})
	if err != nil {
		t.Fatal("Bad test autosigned: ", err)
	}
	conn.Close()

	// Test request https
	fmt.Printf("3) Test certificate autosign by https ( Good Connexion )\n")
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				VerifyPeerCertificate: CheckCertificate("", certificatePEM+"\n"),
				InsecureSkipVerify:    true,
			},
		},
	}
	// Effectuer une requête HTTP GET
	response, err := client.Get("https://127.0.0.1:4329")
	if err != nil {
		fmt.Println("Erreur lors de la requête HTTP:", err)
		return
	}

	// Page not found
	if response.StatusCode != 404 {
		t.Fatal("Response must return 404 not: ", response.StatusCode)
	}
	defer response.Body.Close()

}
