package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/medibloc/vc-sdk/pkg/vc"
	"github.com/medibloc/vc-sdk/pkg/vdr"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(verifyPresentationCmd)
}

var verifyPresentationCmd = &cobra.Command{
	Use:   "verify-presentation",
	Short: "Verify presentation",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		panaceaGRPCAddr := args[0]

		didClient, err := vdr.NewDefaultPanaceaDIDClient(panaceaGRPCAddr)
		if err != nil {
			return fmt.Errorf("failed to init DID client with %v: %w", panaceaGRPCAddr, err)
		}
		defer didClient.Close()

		framework, err := vc.NewFramework(vdr.NewPanaceaVDR(didClient))
		if err != nil {
			return fmt.Errorf("failed to init framework: %w", err)
		}

		verifiablePresentation, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("failed to read input from stdin: %w", err)
		}

		if _, err := framework.VerifyPresentation(verifiablePresentation); err != nil {
			return fmt.Errorf("failed to verify presentation: %w", err)
		}

		credIter, err := framework.GetCredentials(verifiablePresentation)
		if err != nil {
			return fmt.Errorf("failed to get credentials from presentation: %w", err)
		}

		for credIter.HasNext() {
			if err := framework.VerifyCredential(credIter.Next()); err != nil {
				return fmt.Errorf("failed to verify credential in presentation: %w", err)
			}
		}

		fmt.Print(string(verifiablePresentation))

		return nil
	},
}
