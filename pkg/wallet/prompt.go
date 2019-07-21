package wallet

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

// promptPass prompts the user for a passphrase with the given prefix.  The
// function will ask the user to confirm the passphrase and will repeat the
// prompts until they enter a matching response.
func promptPass(reader *bufio.Reader, prefix string, confirm bool) ([]byte, error) {
	// Prompt the user until they enter a passphrase.
	prompt := fmt.Sprintf("%s: ", prefix)
	for {
		fmt.Print(prompt)
		pass, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return nil, err
		}
		fmt.Print("\n")
		pass = bytes.TrimSpace(pass)
		if len(pass) == 0 {
			continue
		}

		if !confirm {
			return pass, nil
		}

		fmt.Print("Confirm passphrase: ")
		confirm, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return nil, err
		}
		fmt.Print("\n")
		confirm = bytes.TrimSpace(confirm)
		if !bytes.Equal(pass, confirm) {
			fmt.Println("The entered passphrases do not match")
			continue
		}

		return pass, nil
	}
}

// promptList prompts the user with the given prefix, list of valid responses,
// and default list entry to use.  The function will repeat the prompt to the
// user until they enter a valid response.
func promptList(reader *bufio.Reader, prefix string, validResponses []string, defaultEntry string) (string, error) {
	// Setup the prompt according to the parameters.
	validStrings := strings.Join(validResponses, "/")
	var prompt string
	if defaultEntry != "" {
		prompt = fmt.Sprintf("%s (%s) [%s]: ", prefix, validStrings,
			defaultEntry)
	} else {
		prompt = fmt.Sprintf("%s (%s): ", prefix, validStrings)
	}

	// Prompt the user until one of the valid responses is given.
	for {
		fmt.Print(prompt)
		reply, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		reply = strings.TrimSpace(strings.ToLower(reply))
		if reply == "" {
			reply = defaultEntry
		}

		for _, validResponse := range validResponses {
			if reply == validResponse {
				return reply, nil
			}
		}
	}
}

// promptListBool prompts the user for a boolean (yes/no) with the given prefix.
// The function will repeat the prompt to the user until they enter a valid
// reponse.
func promptListBool(reader *bufio.Reader, prefix string, defaultEntry string) (bool, error) {
	// Setup the valid responses.
	valid := []string{"n", "no", "y", "yes"}
	response, err := promptList(reader, prefix, valid, defaultEntry)
	if err != nil {
		return false, err
	}
	return response == "yes" || response == "y", nil
}

// PromptProvidePrivatePass prompts the user for a passphrase
func PromptProvidePrivatePass(reader *bufio.Reader) ([]byte, error) {
	return promptPass(reader, "Enter the passphrase for your new wallet", true)
}

// PromptProvideSecret prompts the user for a secret, such as a private key
func PromptProvideSecret(reader *bufio.Reader, promptMessage string) ([]byte, error) {
	return promptPass(reader, promptMessage, false)
}
