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
	rootCmd.AddCommand(verifyCredentialCmd)
}

var verifyCredentialCmd = &cobra.Command{
	Use:   "verify-credential",
	Short: "Verify credential",
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

		verifiableCredential, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("failed to read input from stdin: %w", err)
		}

		if err := framework.VerifyCredential(verifiableCredential); err != nil {
			return fmt.Errorf("failed to verify credential: %w", err)
		}

		fmt.Print(string(verifiableCredential))

		return nil
	},
}
