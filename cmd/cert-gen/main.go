//nolint:gomnd
package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"path"
	"time"
)

/*
Не стал выносить в отдельные библиотеки, так как вообще считаю что сертификат можно выпустить через openssl.
Для учебных целей, реализовал тут.
*/

func main() {
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(time.Now().UnixNano()),
		Subject: pkix.Name{
			Organization: []string{"Practicum"},
			Country:      []string{"UZ"},
		},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Fatal(err)
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &key.PublicKey, key)
	if err != nil {
		log.Fatal(err)
	}

	var keyPEM bytes.Buffer
	if err = pem.Encode(&keyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}); err != nil {
		return
	}

	var certPEM bytes.Buffer
	if err = pem.Encode(&certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	}); err != nil {
		return
	}

	if _, err = os.Stat("cert"); os.IsNotExist(err) {
		if err = os.Mkdir("cert", 0755); err != nil {
			panic(fmt.Errorf("error creating directory: %w", err))
		}
	}

	if err = os.WriteFile(path.Join("cert", "key.pem"), keyPEM.Bytes(), 0600); err != nil {
		log.Fatalf("Failed to write to file: %v", err)
	}

	if err = os.WriteFile(path.Join("cert", "cert.pem"), certPEM.Bytes(), 0600); err != nil {
		log.Fatalf("Failed to write to file: %v", err)
	}
}
