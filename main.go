package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"Transaction_system/requests"

	"github.com/joho/godotenv"
)

const (
	collectURL      = "https://demo.campay.net/api/collect/"
	transactionURL  = "https://demo.campay.net/api/transaction/"
	countryCode     = "237"
	dividerWidth    = 60
	requiredLength  = 9
	firstDigit      = 6
)

// Provider represents a mobile money provider
type Provider string

const (
	ProviderMTN    Provider = "MTN"
	ProviderOrange Provider = "Orange"
	ProviderInvalid Provider = ""
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Load environment variables
	if err := loadEnvFile(); err != nil {
		return err
	}

	// Get and validate API key
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		return fmt.Errorf("API_KEY environment variable is not set")
	}
	fmt.Println("API key loaded successfully")

	// Get phone number with validation
	phoneNumber, err := getValidPhoneNumber()
	if err != nil {
		return err
	}

	printDivider()

	// Get amount
	amount, err := getValidAmount()
	if err != nil {
		return err
	}

	printDivider()

	// Get reference
	reference, err := getReference()
	if err != nil {
		return err
	}

	printDivider()

	// Process payment
	return processPayment(apiKey, phoneNumber, amount, reference)
}

// loadEnvFile loads the .env file if it exists
func loadEnvFile() error {
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			return fmt.Errorf("failed to load .env file: %w", err)
		}
	}
	return nil
}

// getValidPhoneNumber prompts user for a valid phone number
func getValidPhoneNumber() (string, error) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Please enter your mobile money number: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("failed to read input: %w", err)
		}

		phoneNumber := strings.TrimSpace(input)

		if isValidPhoneNumber(phoneNumber) {
			return phoneNumber, nil
		}

		fmt.Println("Invalid number.")
		fmt.Println("Please enter a 9-digit number starting with 6")
		fmt.Printf("Detected provider: %s\n", getProvider(phoneNumber))
	}
}

// isValidPhoneNumber validates the phone number format and provider
func isValidPhoneNumber(number string) bool {
	if len(number) != requiredLength {
		return false
	}

	if number[0] != '0'+firstDigit {
		return false
	}

	provider := getProvider(number)
	return provider == ProviderMTN || provider == ProviderOrange
}

// getProvider determines the mobile money provider based on the phone number
func getProvider(number string) Provider {
	if len(number) < 3 {
		return ProviderInvalid
	}

	secondDigit := number[1] - '0'
	thirdDigit := number[2] - '0'

	switch secondDigit {
	case 5:
		if thirdDigit >= 0 && thirdDigit <= 4 {
			return ProviderMTN
		}
		if thirdDigit >= 5 && thirdDigit <= 8 {
			return ProviderOrange
		}
	case 7:
		if thirdDigit >= 0 && thirdDigit <= 9 {
			return ProviderMTN
		}
	case 8:
		if thirdDigit >= 0 && thirdDigit <= 3 {
			return ProviderMTN
		}
	case 9:
		if thirdDigit >= 0 && thirdDigit <= 9 {
			return ProviderOrange
		}
	}

	return ProviderInvalid
}

// getValidAmount prompts user for a valid amount
func getValidAmount() (int, error) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("\nPlease enter amount: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return 0, fmt.Errorf("failed to read input: %w", err)
		}

		amountStr := strings.TrimSpace(input)
		amount, err := strconv.Atoi(amountStr)
		if err != nil {
			fmt.Println("Invalid amount. Please enter a valid number.")
			continue
		}

		if amount <= 0 {
			fmt.Println("Amount must be greater than zero.")
			continue
		}

		return amount, nil
	}
}

// getReference prompts user for a payment reference
func getReference() (string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("\nReference: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read input: %w", err)
	}

	reference := strings.TrimSpace(input)
	if reference == "" {
		return "", fmt.Errorf("reference cannot be empty")
	}

	return reference, nil
}

// processPayment sends the payment request and checks status
func processPayment(apiKey, phoneNumber string, amount int, description string) error {
	body := map[string]string{
		"from":        countryCode + phoneNumber,
		"amount":      strconv.Itoa(amount),
		"description": description,
	}

	fmt.Println("\nSending payment request...")
	reference := requests.PaymentRequest(collectURL, apiKey, body)

	fmt.Println()
	printDivider()

	if reference == "" {
		return fmt.Errorf("could not retrieve payment reference")
	}

	statusURL := transactionURL + reference + "/"
	requests.GetTransactionStatus(statusURL, apiKey)
	printDivider()

	return nil
}

// printDivider prints a horizontal line separator
func printDivider() {
	fmt.Println(strings.Repeat("-", dividerWidth))
}