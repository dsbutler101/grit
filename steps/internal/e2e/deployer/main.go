package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"cloud.google.com/go/auth/credentials"
	oslogin "cloud.google.com/go/oslogin/apiv1"
	"cloud.google.com/go/oslogin/apiv1/osloginpb"
	"cloud.google.com/go/oslogin/common/commonpb"
	"golang.org/x/crypto/ssh"
)

const (
	e2eCloudEnv = "E2E_CLOUD"

	colorReset       = "\033[0m"
	colorLightRed    = "\033[91m"
	colorLightYellow = "\033[93m"
	colorLightCyan   = "\033[96m"
)

var (
	workdir         = flag.String("workdir", "", "working directory")
	command         = flag.String("command", "", "command to run")
	tfTarget        = flag.String("tf-target", "", "TF target directory")
	additionalFlags = flag.String("additional-flags", "", "additional flags to pass to the command")

	sshCommands = map[string]bool{
		"shutdown":        true,
		"wait-healthy":    true,
		"wait-terminated": true,
	}
)

func main() {
	flag.Parse()

	ctx, cancelFn := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancelFn()

	args := []string{
		"run", ".", *command,
		"--tf-target", *tfTarget,
	}

	af := strings.Split(*additionalFlags, " ")
	args = append(args, af...)

	if _, ok := sshCommands[*command]; isGoogleCloud() && ok {
		cred, err := setupGoogleOSLoginSSHKey(ctx)
		if err != nil {
			printlnColor(colorLightRed, "FATAL: %v", err)
			os.Exit(1)
		}

		args = append(args,
			"--ssh-key-file", cred.SSHKeyFile,
			"--ssh-username", cred.Username,
		)
	}

	printlnColor(colorLightCyan, "Executing go @%s with %v", *workdir, args)

	cmd := exec.CommandContext(ctx, "go", args...)
	cmd.Dir = *workdir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		printlnColor(colorLightRed, "Error running command: %v", err)
		os.Exit(1)
	}
}

func printlnColor(color string, format string, a ...any) {
	fmt.Printf("%s%s%s\n", color, fmt.Sprintf(format, a...), colorReset)
}

func isGoogleCloud() bool {
	return os.Getenv(e2eCloudEnv) == "google"
}

type osLoginCred struct {
	Username   string
	SSHKeyFile string
}

func setupGoogleOSLoginSSHKey(ctx context.Context) (osLoginCred, error) {
	var cred osLoginCred

	if !isGoogleCloud() {
		return cred, nil
	}

	const (
		privKeyName = "id_rsa"
	)

	printlnColor(colorLightYellow, "Preparing Google OS Login SSH key...")

	dir, err := os.MkdirTemp("", "e2d-deployer-test")
	if err != nil {
		return cred, fmt.Errorf("creating temporary directory: %w", err)
	}

	privKeyFile := filepath.Join(dir, privKeyName)

	priv, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return cred, fmt.Errorf("generating private key: %w", err)
	}

	privPEM := pem.EncodeToMemory(&pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   x509.MarshalPKCS1PrivateKey(priv),
	})

	err = os.WriteFile(privKeyFile, privPEM, 0600)
	if err != nil {
		return cred, fmt.Errorf("writing private key at %s: %w", privKeyFile, err)
	}

	pubKey, err := ssh.NewPublicKey(priv.Public())
	if err != nil {
		return cred, fmt.Errorf("generating public key: %w", err)
	}

	pbPubKey := commonpb.SshPublicKey{
		Key:                string(ssh.MarshalAuthorizedKey(pubKey)),
		ExpirationTimeUsec: time.Now().UTC().Add(24 * time.Hour).UnixMicro(),
	}

	creds, err := credentials.DetectDefault(&credentials.DetectOptions{})
	if err != nil {
		return cred, fmt.Errorf("detecting default Google Cloud credentials: %w", err)
	}

	var sa struct {
		ClientEmail string `json:"client_email"`
	}

	err = json.Unmarshal(creds.JSON(), &sa)
	if err != nil {
		return cred, fmt.Errorf("unmarshalling Google Cloud credentials: %w", err)
	}

	pbOsLoginReq := osloginpb.ImportSshPublicKeyRequest{
		Parent:       fmt.Sprintf("users/%s", sa.ClientEmail),
		SshPublicKey: &pbPubKey,
	}

	client, err := oslogin.NewClient(ctx)
	if err != nil {
		return cred, fmt.Errorf("creating oslogin client: %w", err)
	}

	defer client.Close()

	pbOsLoginResp, err := client.ImportSshPublicKey(ctx, &pbOsLoginReq)
	if err != nil {
		return cred, fmt.Errorf("importing ssh public key: %w", err)
	}

	cred.Username = fmt.Sprintf("sa_%s", pbOsLoginResp.GetLoginProfile().Name)
	cred.SSHKeyFile = privKeyFile

	return cred, nil
}
