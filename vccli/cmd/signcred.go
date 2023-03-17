package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/go-bip39"
	didtypes "github.com/medibloc/panacea-core/v2/x/did/types"
	"github.com/medibloc/vc-sdk/pkg/vc"
	"github.com/medibloc/vc-sdk/pkg/vdr"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

func init() {
	rootCmd.AddCommand(signCredentialCmd)
}

var signCredentialCmd = &cobra.Command{
	Use:   "sign-credential",
	Short: "Sign credential",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		panaceaGRPCAddr := args[0]
		mnemonic := args[1]

		privKey, err := privKeyFromMnemonic(mnemonic, 371, 0, 0)
		if err != nil {
			return fmt.Errorf("failed to get private key from mnemonic: %w", err)
		}
		did := deriveDID(privKey)

		didClient, err := vdr.NewDefaultPanaceaDIDClient(panaceaGRPCAddr)
		if err != nil {
			return fmt.Errorf("failed to init DID client with %v: %w", panaceaGRPCAddr, err)
		}
		defer didClient.Close()

		framework, err := vc.NewFramework(vdr.NewPanaceaVDR(didClient))
		if err != nil {
			return fmt.Errorf("failed to init framework: %w", err)
		}

		credential, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("failed to read input from stdin: %w", err)
		}

		vc, err := framework.SignCredential(credential, privKey, &vc.ProofOptions{
			VerificationMethod: fmt.Sprintf("%v#key1", did),
			SignatureType:      "EcdsaSecp256k1Signature2019",
		})
		if err != nil {
			return fmt.Errorf("failed to sign credential: %w", err)
		}

		fmt.Print(string(vc))

		return nil
	},
}

func privKeyFromMnemonic(mnemonic string, coinType, accNum, index uint32) (secp256k1.PrivKey, error) {
	hdPath := hd.NewFundraiserParams(accNum, coinType, index).String()
	master, ch := hd.ComputeMastersFromSeed(bip39.NewSeed(mnemonic, ""))
	return hd.DerivePrivateKeyForPath(master, ch, hdPath)
}

func deriveDID(privKey secp256k1.PrivKey) string {
	pubKey := privKey.PubKey().(secp256k1.PubKey)[:]
	return didtypes.NewDID(pubKey)
}
