package main

import (
    "bufio"
    "flag"
    "fmt"
    "os"
    "strings"
    
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
        topLeft:     "┌", topRight: "┐", bottomLeft: "└", bottomRight: "┘",
        horizontal:  "─", vertical: "│", titleLeft: "┘", titleRight: "└",
    },
    2: {
        topLeft:     "╭", topRight: "╮", bottomLeft: "╰", bottomRight: "╯",
        horizontal:  "─", vertical: "│", titleLeft: "╯", titleRight: "╰",
    },
    3: {
        topLeft:     "╔", topRight: "╗", bottomLeft: "╚", bottomRight: "╝",
        horizontal:  "═", vertical: "║", titleLeft: "╝", titleRight: "╚",
    },
}

// visualLength returns the visual width of the string considering the character widths in different writing systems.
func visualLength(s string) int {
    length := 0
    for _, r := range s {
        rw := runewidth.RuneWidth(r)
        if rw > 1 {
            // For characters wider than 1, increase the width.
            length += rw
        } else {
            // Otherwise, assume a width of 1.
            length++
        }
    }
    return length
}

// max returns the larger of two integers.
func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}

func main() {
    // Read parameters.
    // Note: Only styles 1 to 3 are available.
    styleNum := flag.Int("n", 1, "Box style (1-3)")
    title := flag.String("t", "", "Box title")
    // Default: left-aligned; with -c, centered.
    center := flag.Bool("c", false, "Center text")
    flag.Parse()

    style, ok := styles[*styleNum]
    if !ok {
        fmt.Fprintln(os.Stderr, "Invalid style number. Please use 1-3")
        os.Exit(1)
    }

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

    // Standard padding: 1 space on each side.
    minPadding := 2
    // Inner width: longest content + 2 (1 left and 1 right).
    innerWidth := maxContentWidth + minPadding

    // If a title is provided, create a decorated title.
    var titleDecor string
    if *title != "" {
        // Example: "┘ *title* └" for style 1 or "╝ *title* ╚" for style 3.
        titleDecor = style.titleLeft + " *" + *title + "* " + style.titleRight
        // Adjust innerWidth if the decorated title is longer.
        if visualLength(titleDecor) > innerWidth {
            innerWidth = visualLength(titleDecor)
        }
    }

    // Top box line (with title if available)
    if *title != "" {
        remaining := innerWidth - visualLength(titleDecor)
        leftFill := remaining / 2
        rightFill := remaining - leftFill
        leftHor := strings.Repeat(style.horizontal, leftFill)
        rightHor := strings.Repeat(style.horizontal, rightFill)
        fmt.Printf("%s%s%s%s%s\n",
            style.topLeft,
            leftHor,
            titleDecor,
            rightHor,
            style.topRight)
    } else {
        fmt.Printf("%s%s%s\n",
            style.topLeft,
            strings.Repeat(style.horizontal, innerWidth),
            style.topRight)
    }

    // Output content lines – left-aligned (default) or centered (-c).
    for _, line := range lines {
        pad := innerWidth - visualLength(line)
        if *center {
            // Centered: equal padding on both sides.
            leftPad := pad / 2
            rightPad := pad - leftPad
            fmt.Printf("%s%s%s%s%s\n",
                style.vertical,
                strings.Repeat(" ", leftPad),
                line,
                strings.Repeat(" ", rightPad),
                style.vertical)
        } else {
            // Left-aligned: 1 space on the left, the rest on the right.
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

    // Bottom box line.
    fmt.Printf("%s%s%s\n",
        style.bottomLeft,
        strings.Repeat(style.horizontal, innerWidth),
        style.bottomRight)
}
