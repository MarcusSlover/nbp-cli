// Package cli provides the command-line interface logic for the NBP tool.
// It handles argument parsing, date validation, and result presentation.
package cli

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/MarcusSlover/nbp-cli/internal/nbp"
	"github.com/spf13/cobra"
)

var (
	mode    string
	taxFlag bool
)

var rootCmd = &cobra.Command{
	Use:   "nbp [currency] [date|today]",
	Short: "NBP Exchange Rate CLI",
	Long:  `A simple and efficient Go-based command-line tool to fetch currency exchange rates from the National Bank of Poland (NBP).`,
	Args:  cobra.MaximumNArgs(2),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			// Complete currency codes
			codes, err := nbp.FetchAvailableCurrencies()
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}
			return codes, cobra.ShellCompDirectiveNoFileComp
		}
		if len(args) == 1 {
			// Complete "today" for the second argument
			if strings.HasPrefix("today", strings.ToLower(toComplete)) {
				return []string{"today"}, cobra.ShellCompDirectiveNoFileComp
			}
			// Always return NoFileComp even if no match to avoid fallback to files
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			cmd.Help()
			os.Exit(1)
		}

		currency := strings.ToUpper(args[0])
		inputDate := strings.ToLower(args[1])
		isTaxMode := mode == "tax" || taxFlag

		var targetDate time.Time
		if inputDate == "today" {
			targetDate = time.Now()
		} else {
			var err error
			targetDate, err = time.Parse("2006-01-02", inputDate)
			if err != nil {
				fmt.Printf("Error: Date must be YYYY-MM-DD, got: %s\n", inputDate)
				os.Exit(1)
			}
		}

		if isTaxMode {
			handleTaxMode(currency, targetDate)
		} else {
			handleNormalMode(currency, inputDate)
		}
	},
}

// Run parses command-line flags and arguments, then executes the requested mode.
// It terminates the program with an exit code 1 if the input is invalid or an error occurs.
func Run() {
	rootCmd.Flags().StringVarP(&mode, "mode", "m", "normal", "Mode to use: 'normal' for current rate, 'tax' for previous business day")
	rootCmd.Flags().BoolVarP(&taxFlag, "tax", "t", false, "Shortcut for --mode tax")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// handleNormalMode fetches and prints the currency rate for a specific date or "today".
func handleNormalMode(currency, inputDate string) {
	data, err := nbp.FetchNormalRate(currency, inputDate)
	if err != nil {
		fmt.Printf("Error: Could not find a rate for %s on %s. (Note: NBP doesn't publish on weekends)\n", currency, inputDate)
		return
	}

	PrintResult("Normal Rate", data.Rates[0], data.Code, nil)
}

// handleTaxMode calculates and prints the tax-compliant rate (last business day before targetDate).
func handleTaxMode(currency string, targetDate time.Time) {
	data, err := nbp.FetchTaxRate(currency, targetDate)
	if err != nil {
		fmt.Printf("Error: Could not calculate tax rate. Check your currency code: %v\n", err)
		return
	}

	// Use the last element in the range (closest to the transaction date)
	lastBusinessDayRate := data.Rates[len(data.Rates)-1]
	PrintResult("Tax Rate (Previous Business Day)", lastBusinessDayRate, data.Code, &targetDate)
}

// PrintResult displays the formatted exchange rate information to the standard output.
// label describes the mode (Normal or Tax).
// rate contains the NBP rate details.
// code is the ISO 4217 currency code.
// transactionDate is the date of the transaction (optional, used in tax mode).
func PrintResult(label string, rate nbp.Rate, code string, transactionDate *time.Time) {
	fmt.Printf("--- %s ---\n", label)
	if transactionDate != nil {
		fmt.Printf("Transaction Date: %s\n", transactionDate.Format("2006-01-02"))
	}
	fmt.Printf("NBP Rate Date:    %s\n", rate.EffectiveDate)
	fmt.Printf("Table Number:     %s\n", rate.No)
	fmt.Printf("Value:            1 %s = %.4f PLN\n", code, rate.Mid)
}
