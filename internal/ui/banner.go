package ui

import (
	"fmt"
	"strings"
)

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorBold   = "\033[1m"
	ColorDim    = "\033[2m"

	// Bright colors
	ColorBrightBlue   = "\033[94m"
	ColorBrightCyan   = "\033[96m"
	ColorBrightGreen  = "\033[92m"
	ColorBrightYellow = "\033[93m"
)

const archieBanner = `
  ██████╗  ██████╗   ██████╗ ██╗  ██╗ ██╗ ███████╗
  ██╔══██╗ ██╔══██╗ ██╔════╝ ██║  ██║ ██║ ██╔════╝
  ███████║ ██████╔╝ ██║      ███████║ ██║ █████╗
  ██╔══██║ ██╔══██╗ ██║      ██╔══██║ ██║ ██╔══╝
  ██║  ██║ ██║  ██║ ╚██████╗ ██║  ██║ ██║ ███████╗
  ╚═╝  ╚═╝ ╚═╝  ╚═╝  ╚═════╝ ╚═╝  ╚═╝ ╚═╝ ╚══════╝
`

// ShowWelcomeBanner displays the welcome banner for archie init
func ShowWelcomeBanner() {
	// Clear screen (optional)
	// fmt.Print("\033[H\033[2J")

	fmt.Println()

	// Print banner with color gradient
	printColoredBanner()

	fmt.Println()

	// Welcome message
	centerPrint(ColorBold + ColorBrightCyan + "Welcome to ARCHIE" + ColorReset)
	fmt.Println()

	// Tagline
	centerPrint(ColorDim + "AI-Powered Technical Design Documentation System" + ColorReset)
	fmt.Println()

	// Separator
	printSeparator()
	fmt.Println()

	// Agent configuration prompt
	centerPrint(ColorBold + ColorBrightYellow + "🤖 Please Configure Your AI Agent" + ColorReset)
	fmt.Println()
	fmt.Println()
}

// printColoredBanner prints the ARCHIE banner with gradient colors
func printColoredBanner() {
	lines := strings.Split(archieBanner, "\n")
	colors := []string{
		ColorBrightBlue,
		ColorBrightCyan,
		ColorCyan,
		ColorBrightCyan,
		ColorBrightBlue,
		ColorBlue,
	}

	for i, line := range lines {
		if line == "" {
			continue
		}

		// Use color gradient
		colorIndex := i % len(colors)
		fmt.Println(colors[colorIndex] + line + ColorReset)
	}
}

// centerPrint prints text centered (assuming 80 char width)
func centerPrint(text string) {
	// Remove color codes for length calculation
	cleanText := removeColorCodes(text)
	width := 80
	padding := (width - len(cleanText)) / 2

	if padding > 0 {
		fmt.Print(strings.Repeat(" ", padding))
	}
	fmt.Println(text)
}

// printSeparator prints a decorative separator
func printSeparator() {
	separator := strings.Repeat("─", 70)
	fmt.Println(ColorDim + "  " + separator + ColorReset)
}

// removeColorCodes removes ANSI color codes from text for length calculation
func removeColorCodes(text string) string {
	// Simple removal of common ANSI codes
	result := text
	codes := []string{
		ColorReset, ColorRed, ColorGreen, ColorYellow, ColorBlue,
		ColorPurple, ColorCyan, ColorWhite, ColorBold, ColorDim,
		ColorBrightBlue, ColorBrightCyan, ColorBrightGreen, ColorBrightYellow,
	}

	for _, code := range codes {
		result = strings.ReplaceAll(result, code, "")
	}

	return result
}

// ShowSuccess displays a success message
func ShowSuccess(message string) {
	fmt.Println()
	fmt.Println(ColorBrightGreen + "  ✅ " + message + ColorReset)
	fmt.Println()
}

// ShowError displays an error message
func ShowError(message string) {
	fmt.Println()
	fmt.Println(ColorRed + "  ❌ " + message + ColorReset)
	fmt.Println()
}

// ShowInfo displays an info message
func ShowInfo(message string) {
	fmt.Println(ColorBrightCyan + "  ℹ️  " + message + ColorReset)
}

// ShowStep displays a step message
func ShowStep(step int, total int, message string) {
	fmt.Printf(ColorBrightYellow+"  [%d/%d] "+ColorReset+"%s\n", step, total, message)
}

// ShowSubStep displays a sub-step with an arrow
func ShowSubStep(message string) {
	fmt.Println(ColorDim + "      → " + ColorReset + message)
}

// PrintBox prints text in a box
func PrintBox(title string, content []string) {
	width := 60

	// Top border
	fmt.Println(ColorBrightCyan + "  ╔" + strings.Repeat("═", width-2) + "╗" + ColorReset)

	// Title
	if title != "" {
		padding := (width - len(title) - 4) / 2
		fmt.Print(ColorBrightCyan + "  ║ " + ColorReset)
		fmt.Print(strings.Repeat(" ", padding))
		fmt.Print(ColorBold + title + ColorReset)
		fmt.Print(strings.Repeat(" ", width-len(title)-padding-4))
		fmt.Println(ColorBrightCyan + " ║" + ColorReset)

		// Separator after title
		fmt.Println(ColorBrightCyan + "  ╠" + strings.Repeat("═", width-2) + "╣" + ColorReset)
	}

	// Content
	for _, line := range content {
		cleanLine := removeColorCodes(line)
		padding := width - len(cleanLine) - 4
		fmt.Print(ColorBrightCyan + "  ║ " + ColorReset)
		fmt.Print(line)
		if padding > 0 {
			fmt.Print(strings.Repeat(" ", padding))
		}
		fmt.Println(ColorBrightCyan + " ║" + ColorReset)
	}

	// Bottom border
	fmt.Println(ColorBrightCyan + "  ╚" + strings.Repeat("═", width-2) + "╝" + ColorReset)
	fmt.Println()
}
