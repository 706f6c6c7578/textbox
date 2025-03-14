package main

import (
    "bufio"
    "flag"
    "fmt"
    "os"
    "strings"
    "unicode/utf8"

    "github.com/mattn/go-runewidth"
)

// BoxStyle contains the characters for the various frame components.
type BoxStyle struct {
    topLeft     string
    topRight    string
    bottomLeft  string
    bottomRight string
    horizontal  string
    vertical    string
    titleLeft   string
    titleRight  string
}

// Different styles to choose from.
var styles = map[int]BoxStyle{
    1: {
        topLeft: "┌", topRight: "┐", bottomLeft: "└", bottomRight: "┘",
        horizontal: "─", vertical: "│", titleLeft: "┘", titleRight: "└",
    },
    2: {
        topLeft: "╭", topRight: "╮", bottomLeft: "╰", bottomRight: "╯",
        horizontal: "─", vertical: "│", titleLeft: "╯", titleRight: "╰",
    },
    3: {
        topLeft: "╔", topRight: "╗", bottomLeft: "╚", bottomRight: "╝",
        horizontal: "═", vertical: "║", titleLeft: "╝", titleRight: "╚",
    },
}

// visualLength returns the visual width of the string considering the character widths in different writing systems.
func visualLength(s string) int {
    return runewidth.StringWidth(s)
}

// max returns the larger of two integers.
func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}

func repeatChar(char string, count int) string {
    result := ""
    for i := 0; i < count; i++ {
        result += char
    }
    return result
}

func main() {
    // Read parameters.
    styleNum := flag.Int("n", 1, "Box style (1-4)")
    customChar := flag.String("f", "", "Custom UTF-8 character for style 4")
    title := flag.String("t", "", "Box title")
    center := flag.Bool("c", false, "Center text")
    flag.Parse()

    var style BoxStyle

    // Validate style number and custom character.
    if *styleNum >= 1 && *styleNum <= 3 {
        style = styles[*styleNum]
    } else if *styleNum == 4 {
        // Trim whitespace and validate rune count.
        utfChar := strings.TrimSpace(*customChar)
        if utf8.RuneCountInString(utfChar) != 1 {
            fmt.Fprintln(os.Stderr, "Error: For -n 4, exactly one UTF-8 character must be provided with -f.")
            os.Exit(1)
        }
        style = BoxStyle{
            topLeft: utfChar, topRight: utfChar, bottomLeft: utfChar, bottomRight: utfChar,
            horizontal: utfChar, vertical: utfChar, titleLeft: utfChar, titleRight: utfChar,
        }
    } else {
        fmt.Fprintln(os.Stderr, "Invalid style number or missing custom character. Please use -n 1-4 or provide a custom character with -f.")
        os.Exit(1)
    }

    // Read input lines and calculate maximum content width.
    var lines []string
    scanner := bufio.NewScanner(os.Stdin)
    maxContentWidth := 0
    for scanner.Scan() {
        line := scanner.Text()
        lines = append(lines, line)
        if l := visualLength(line); l > maxContentWidth {
            maxContentWidth = l
        }
    }

    minPadding := 2
    innerWidth := maxContentWidth + minPadding

    // Handle title decoration.
    var titleDecor string
    if *title != "" {
        titleDecor = style.titleLeft + " " + *title + " " + style.titleRight
        if visualLength(titleDecor) > innerWidth {
            innerWidth = visualLength(titleDecor)
        }
    }

    // Generate the top border.
    if *title != "" {
        remaining := innerWidth - visualLength(titleDecor)
        leftFill := remaining / 2
        rightFill := remaining - leftFill
        leftHor := repeatChar(style.horizontal, leftFill/visualLength(style.horizontal))
        rightHor := repeatChar(style.horizontal, rightFill/visualLength(style.horizontal))
        fmt.Printf("%s%s%s%s%s\n",
            style.topLeft,
            leftHor,
            titleDecor,
            rightHor,
            style.topRight)
    } else {
        lineWidth := innerWidth / visualLength(style.horizontal)
        fmt.Printf("%s%s%s\n",
            style.topLeft,
            repeatChar(style.horizontal, lineWidth),
            style.topRight)
    }

    // Print the content.
    for _, line := range lines {
        pad := innerWidth - visualLength(line)
        if *center {
            leftPad := pad / 2
            rightPad := pad - leftPad
            fmt.Printf("%s%s%s%s%s\n",
                style.vertical,
                strings.Repeat(" ", leftPad),
                line,
                strings.Repeat(" ", rightPad),
                style.vertical)
        } else {
            leftPad := 1
            rightPad := pad - leftPad
            if rightPad < 0 {
                rightPad = 0
            }
            fmt.Printf("%s%s%s%s%s\n",
                style.vertical,
                strings.Repeat(" ", leftPad),
                line,
                strings.Repeat(" ", rightPad),
                style.vertical)
        }
    }

    // Generate the bottom border.
    lineWidth := innerWidth / visualLength(style.horizontal)
    fmt.Printf("%s%s%s\n",
        style.bottomLeft,
        repeatChar(style.horizontal, lineWidth),
        style.bottomRight)
}