// SPDX-FileCopyrightText: 2020 The tls-interop-runner Authors
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: 2009 The Go Authors
// SPDX-License-Identifier: BSD-3-Clause

// This file is based on code found in
// https://boringssl.googlesource.com/boringssl/+/refs/heads/master/ssl/test/runner/

package main

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	_ "crypto/sha256"
	_ "crypto/sha512"
	"fmt"
	"io"
)

// Signer represents an structure holding the signing information
type Signer struct {
	bugs  *CertificateBugs
	curve elliptic.Curve
	hash  crypto.Hash
	priv  crypto.PrivateKey
	rand  io.Reader
	ecdsa bool
}

// GenerateKey generates a public and private key pair.
// TODO: as this is used beyond DCs, it needs to support all the other algos.
func (e *Signer) GenerateKey() (crypto.PrivateKey, crypto.PublicKey, error) {
	var privK crypto.PrivateKey
	var pubK crypto.PublicKey
	var err error

	if e.ecdsa == true {
		privK, err = ecdsa.GenerateKey(e.curve, e.rand)
		if err != nil {
			return nil, nil, err
		}
		pubK = privK.(*ecdsa.PrivateKey).Public()
	} else {
		pubK, privK, err = ed25519.GenerateKey(e.rand)
		if err != nil {
			return nil, nil, err
		}
	}

	e.priv = privK
	return privK, pubK, err
}

var directSigning crypto.Hash = 0

// SignWithKey sings a message with the appropriate key. It is only used by
// delegated credentials, and only supports algorithms allowed for them.
func (e *Signer) SignWithKey(key crypto.PrivateKey, msg []byte) ([]byte, error) {
	var digest []byte
	if e.hash != directSigning {
		h := e.hash.New()
		h.Write(msg)
		digest = h.Sum(nil)
	}

	var sig []byte
	var err error
	switch sk := key.(type) {
	case *ecdsa.PrivateKey:
		opts := crypto.SignerOpts(e.hash)
		sig, err = sk.Sign(e.rand, digest, opts)
		if err != nil {
			fatalIfErr(err, "failed to sign parameters")
			return nil, err
		}
	case ed25519.PrivateKey:
		opts := crypto.SignerOpts(e.hash)
		sig, err = sk.Sign(e.rand, msg, opts)
		if err != nil {
			fatalIfErr(err, "failed to sign parameters")
			return nil, err
		}
	default:
		return nil, fmt.Errorf("tls: unsupported key type")
	}

	return sig, nil
}

// TODO: as this is used beyond DCs, it needs to support all the other algos.
func getSigner(bugs *CertificateBugs, rand io.Reader, sigAlg signatureAlgorithm) (*Signer, error) {
	switch sigAlg {
	case signatureECDSAWithSHA1:
		return &Signer{bugs, nil, crypto.SHA1, nil, rand, true}, nil
	case signatureECDSAWithP256AndSHA256:
		return &Signer{bugs, elliptic.P256(), crypto.SHA256, nil, rand, true}, nil
	case signatureECDSAWithP384AndSHA384:
		return &Signer{bugs, elliptic.P384(), crypto.SHA384, nil, rand, true}, nil
	case signatureECDSAWithP521AndSHA512:
		return &Signer{bugs, elliptic.P521(), crypto.SHA512, nil, rand, true}, nil
	case signatureEd25519:
		return &Signer{bugs, nil, directSigning, nil, rand, false}, nil
	}

	return nil, fmt.Errorf("unsupported signature algorithm %04x", sigAlg)
}
