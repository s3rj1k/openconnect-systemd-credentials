// SPDX-License-Identifier: MIT
// Copyright 2024 s3rj1k.

package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/cristalhq/otp"
)

const (
	CredentialsDirectoryEnvironKey = "CREDENTIALS_DIRECTORY"

	DefaultTOTPSkew      = 2
	DefaultDirectoryMode = 0o755
)

const OpenConnectConfigTemplate = `
form-entry={{.FormEntryKey}}={{.FormEntryValue}}
`

type ConfigTemplateData struct {
	FormEntryKey   string
	FormEntryValue string
}

func GetOTPAuthKey(otpAuth string) (*otp.Key, error) {
	var otpAuthURL string

	if !strings.HasPrefix(otpAuth, "file:") && !strings.HasPrefix(otpAuth, "key:") {
		return nil, errors.New("invalid \"otp-auth\" prefix, must be \"file:\" or \"key:\"")
	}

	if strings.HasPrefix(otpAuth, "file:") {
		fp := filepath.Clean(strings.TrimPrefix(otpAuth, "file:"))

		b, err := exec.Command("systemd-creds", "--with-key=tpm2", "decrypt", fp, "-").Output()
		if err != nil {
			return nil, fmt.Errorf("error executing \"systemd-creds\": %w", err)
		}

		otpAuthURL = string(bytes.TrimSpace(b))
	}

	if strings.HasPrefix(otpAuth, "key:") {
		key := strings.TrimPrefix(otpAuth, "key:")

		dir := os.Getenv(CredentialsDirectoryEnvironKey)
		if dir == "" {
			return nil, fmt.Errorf("environment variable %q is not set", CredentialsDirectoryEnvironKey)
		}

		b, err := os.ReadFile(filepath.Join(dir, key))
		if err != nil {
			return nil, fmt.Errorf("error reading secret %q: %w", key, err)
		}

		otpAuthURL = string(bytes.TrimSpace(b))
	}

	otpAuthKey, err := otp.ParseKeyFromURL(otpAuthURL)
	if err != nil {
		return nil, errors.New("error parsing OTPAuthURL")
	}

	return otpAuthKey, nil
}

func GenerateOTPCode(otpAuthKey *otp.Key) (string, error) {
	switch otpAuthKey.Type() {
	case "totp":
		totp, err := otp.NewTOTP(otp.TOTPConfig{
			Algo:   otpAuthKey.Algorithm(),
			Digits: otpAuthKey.Digits(),
			Issuer: otpAuthKey.Issuer(),
			Period: otpAuthKey.Period(),
			Skew:   DefaultTOTPSkew,
		})
		if err != nil {
			return "", fmt.Errorf("error creating TOTP config: %w", err)
		}

		code, err := totp.GenerateCode(otpAuthKey.Secret(), time.Now())
		if err != nil {
			return "", fmt.Errorf("error generating TOTP code: %w", err)
		}

		return code, nil
	case "hotp":
		hotp, err := otp.NewHOTP(otp.HOTPConfig{
			Algo:   otpAuthKey.Algorithm(),
			Digits: otpAuthKey.Digits(),
			Issuer: otpAuthKey.Issuer(),
		})
		if err != nil {
			return "", fmt.Errorf("error creating HOTP config: %w", err)
		}

		code, err := hotp.GenerateCode(otpAuthKey.Counter(), otpAuthKey.Secret())
		if err != nil {
			return "", fmt.Errorf("error generating HOTP code: %w", err)
		}

		return code, nil
	}

	return "", errors.New("unknown OTP type")
}

func CreateConfig(location, otpAuth, formEntryKey string) error {
	if location == "" {
		return errors.New("config location is not set")
	}

	data := ConfigTemplateData{
		FormEntryKey: formEntryKey,
	}

	tmpl, err := template.New("config").Parse(strings.TrimSpace(OpenConnectConfigTemplate))
	if err != nil {
		return fmt.Errorf("error parsing config template: %w", err)
	}

	err = os.MkdirAll(filepath.Dir(location), DefaultDirectoryMode)
	if err != nil {
		return fmt.Errorf("error creating config directory: %w", err)
	}

	f, err := os.Create(location)
	if err != nil {
		return fmt.Errorf("error creating config: %w", err)
	}

	defer f.Close()

	otpAuthKey, err := GetOTPAuthKey(otpAuth)
	if err != nil {
		return err
	}

	data.FormEntryValue, err = GenerateOTPCode(otpAuthKey)
	if err != nil {
		return err
	}

	err = tmpl.Execute(f, data)
	if err != nil {
		return fmt.Errorf("error rendering config template: %w", err)
	}

	return nil
}

func main() {
	var (
		flagLocation     string
		flagOTPAuth      string
		flagFormEntryKey string
		flagVersion      bool
	)

	flag.StringVar(&flagLocation, "config", "", "OpenConnect VPN config file location.")
	flag.StringVar(&flagOTPAuth, "otp-auth", "",
		"OTP Authentication URL (file:/path/to/encrypted/file or key:name of systemd credential object).",
	)
	flag.StringVar(&flagFormEntryKey, "form-entry", "main:secondary_password", "OpenConnect VPN config \"form-entry\" key name.")
	flag.BoolVar(&flagVersion, "version", false, "Show build info and exit.")

	flag.Parse()

	if flagVersion {
		fmt.Println(GetVCSBuildInfo())

		os.Exit(0)
	}

	if err := CreateConfig(flagLocation, flagOTPAuth, flagFormEntryKey); err != nil {
		fmt.Printf("Failed to create OpenConnect VPN config: %v\n", err)
	} else {
		fmt.Println("OpenConnect VPN config file created successfully")
	}
}
