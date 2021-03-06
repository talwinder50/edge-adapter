/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package crypto

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/hyperledger/aries-framework-go/pkg/doc/verifiable"
	cryptomock "github.com/hyperledger/aries-framework-go/pkg/mock/crypto"
	"github.com/hyperledger/aries-framework-go/pkg/mock/kms"
	vdrimock "github.com/hyperledger/aries-framework-go/pkg/mock/vdri"
	"github.com/stretchr/testify/require"

	"github.com/trustbloc/edge-adapter/pkg/internal/mock/diddoc"
)

func TestSignCredential(t *testing.T) {
	t.Run("test sign vc - success", func(t *testing.T) {
		didDoc := diddoc.GetMockDIDDoc("did:example:abc789")

		c := New(&kms.KeyManager{}, &cryptomock.Crypto{}, &vdrimock.MockVDRIRegistry{ResolveValue: didDoc})

		vc := &verifiable.Credential{ID: uuid.New().URN()}
		signingKeyID := didDoc.AssertionMethod[0].PublicKey.ID

		signedVC, err := c.SignCredential(vc, signingKeyID)
		require.NoError(t, err)
		require.Equal(t, vc.ID, signedVC.ID)
		require.Equal(t, AssertionMethod, signedVC.Proofs[0]["proofPurpose"])
		require.Equal(t, didDoc.AssertionMethod[0].PublicKey.ID, signedVC.Proofs[0]["verificationMethod"])
		require.NotEmpty(t, signedVC.Proofs[0]["created"])
	})

	t.Run("test sign vc - error", func(t *testing.T) {
		// invalid signing key value
		didDoc := diddoc.GetMockDIDDoc("did:example:xyz123")

		c := New(&kms.KeyManager{}, &cryptomock.Crypto{}, &vdrimock.MockVDRIRegistry{ResolveValue: didDoc})

		vc := &verifiable.Credential{ID: uuid.New().URN()}
		signingKeyID := "invalid_key_format"

		signedVC, err := c.SignCredential(vc, signingKeyID)
		require.Error(t, err)
		require.Contains(t, err.Error(), "sign credential : verificationMethod value")
		require.Nil(t, signedVC)

		// signing key not exists
		signingKeyID = "did:example:xyz123#invalidKey"
		signedVC, err = c.SignCredential(vc, signingKeyID)
		require.Error(t, err)
		require.Contains(t, err.Error(),
			"sign credential : validate did doc : unable to find matching assertionMethod key IDs for given "+
				"verification method did:example:xyz123#invalidKey")
		require.Nil(t, signedVC)

		// did resolve error
		c = New(&kms.KeyManager{}, &cryptomock.Crypto{},
			&vdrimock.MockVDRIRegistry{ResolveErr: errors.New("resolve error")})
		signedVC, err = c.SignCredential(vc, signingKeyID)
		require.Error(t, err)
		require.Contains(t, err.Error(), "sign credential : validate did doc : resolve error")
		require.Nil(t, signedVC)
	})
}

func TestSignPresentation(t *testing.T) {
	t.Run("test sign vp - failure", func(t *testing.T) {
		didDoc := diddoc.GetMockDIDDoc("did:example:xyz123")

		c := New(&kms.KeyManager{}, &cryptomock.Crypto{}, &vdrimock.MockVDRIRegistry{ResolveValue: didDoc})

		vp := &verifiable.Presentation{ID: uuid.New().URN()}
		signingKeyID := didDoc.AssertionMethod[0].PublicKey.ID

		signedVP, err := c.SignPresentation(vp, signingKeyID)
		require.NoError(t, err)
		require.Equal(t, Authentication, signedVP.Proofs[0]["proofPurpose"])
		require.Equal(t, didDoc.AssertionMethod[0].PublicKey.ID, signedVP.Proofs[0]["verificationMethod"])
		require.NotEmpty(t, signedVP.Proofs[0]["created"])
	})

	t.Run("test sign vp - error", func(t *testing.T) {
		didDoc := diddoc.GetMockDIDDoc("did:example:xyz123")

		c := New(&kms.KeyManager{}, &cryptomock.Crypto{}, &vdrimock.MockVDRIRegistry{ResolveValue: didDoc})

		vp := &verifiable.Presentation{ID: uuid.New().URN()}
		signingKeyID := "invalid_key_format"

		signedVP, err := c.SignPresentation(vp, signingKeyID)
		require.Error(t, err)
		require.Contains(t, err.Error(), "sign presentation")
		require.Nil(t, signedVP)
	})
}
