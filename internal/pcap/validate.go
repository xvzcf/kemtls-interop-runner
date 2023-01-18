// SPDX-FileCopyrightText: 2020 The tls-interop-runner Authors
// SPDX-License-Identifier: MIT

package pcap

import (
	"errors"
)

// Validate takes a transcript of a TLS handshake and a testcase, and checks
// whether the transcript conforms with the transcript expected for the
// testcase.
func Validate(transcript TLSTranscript, testCase string) error {
	switch testCase {
	case "dc":
		if transcript.ClientHello.version != 0x0303 {
			return errors.New("ClientHello: legacy_version is not TLS 1.2")
		}
		if !transcript.ClientHello.supportsDC {
			return errors.New("ClientHello: support for delegated credentials not indicated")
		}
		if transcript.ClientHello.serverName != "example.com" {
			return errors.New("ClientHello: SNI should specify example.com")
		}
		for _, v := range transcript.ClientHello.supportedVersions {
			if v == 0x0304 {
				return nil
			}
		}
		return errors.New("ClientHello: supported_versions does not include TLS 1.3")
	case "ech-accept":
		if transcript.ClientHello.version != 0x0303 {
			return errors.New("ClientHello: legacy_version is not TLS 1.2")
		}
		for _, v := range transcript.ClientHello.supportedVersions {
			if v == 0x0304 {
				return nil
			}
		}
		return errors.New("ClientHello: supported_versions does not include TLS 1.3")
	case "ech-reject":
		// Check sni in clienthello outer. 
		// SHOULD NOT be the backend server sni backend.com
		if transcript.ClientHello.serverName != "client-facing.com" {
			return errors.New("ClientHello: SNI should specify client-facing.com")
		}
		return nil
	}
	return nil
}
