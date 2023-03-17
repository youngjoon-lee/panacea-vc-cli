package cmd

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/medibloc/vc-sdk/pkg/vc"
	"github.com/medibloc/vc-sdk/pkg/vdr"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(signPresentationCmd)
}

var signPresentationCmd = &cobra.Command{
	Use:   "sign-presentation",
	Short: "Sign presentation",
	Args:  cobra.ExactArgs(4),
	RunE: func(cmd *cobra.Command, args []string) error {
		panaceaGRPCAddr := args[0]
		mnemonic := args[1]
		domain := args[2]
		challenge := args[3]

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

		verifiableCredential, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("failed to read input from stdin: %w", err)
		}

		presentation := []byte(fmt.Sprintf(`{
			"@context": ["https://www.w3.org/2018/credentials/v1"],
			"id": "%v",
			"type": ["VerifiablePresentation"],
			"verifiableCredential": [%s]
		}`, uuid.New().String(), string(verifiableCredential)))

		vp, err := framework.SignPresentation(presentation, privKey, &vc.ProofOptions{
			VerificationMethod: fmt.Sprintf("%v#key1", did),
			SignatureType:      "EcdsaSecp256k1Signature2019",
			Domain:             domain,
			Challenge:          challenge,
			Created:            time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		})
		if err != nil {
			return fmt.Errorf("failed to sign presentation: %w", err)
		}

		fmt.Print(string(vp))

		return nil
	},
}
