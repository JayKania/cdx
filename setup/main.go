package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
)

const shellFunction = `
cdx() {
   local dir=$(command cdx "$@")
   if [ -d "$dir" ]; then
      cd "$dir"
   fi
}
`

const windowsFunction = `
function cdx {
    $dir = & cdx.exe $args
    if (Test-Path $dir -PathType Container) {
        Set-Location $dir
    }
}
`

func main() {
	usr, err := user.Current()
	if err != nil {
		fmt.Println("Error determining user home:", err)
		os.Exit(1)
	}

	fmt.Println("Home directory:", usr.HomeDir)

	switch runtime.GOOS {
	case "darwin", "linux":
		rcFiles := []string{".zshrc", ".bashrc", ".bash_profile"}
		for _, rc := range rcFiles {
			rcPath := filepath.Join(usr.HomeDir, rc)
			if fileExists(rcPath) {
				appendIfMissing(rcPath, shellFunction)
			}
		}
	case "windows":
		profile := os.Getenv("USERPROFILE")
		if profile == "" {
			fmt.Println("❌ Could not determine PowerShell profile path")
			os.Exit(1)
		}
		powershellProfile := filepath.Join(profile, "Documents", "WindowsPowerShell", "Microsoft.PowerShell_profile.ps1")
		// Create directory if it doesn't exist
		os.MkdirAll(filepath.Dir(powershellProfile), 0755)
		appendIfMissing(powershellProfile, windowsFunction)
	default:
		fmt.Printf("❌ OS %s not supported\n", runtime.GOOS)
	}
}

func appendIfMissing(filePath, content string) {
	// Check if file exists
	if !fileExists(filePath) {
		// Create the file if it doesn't exist
		err := os.WriteFile(filePath, []byte(content+"\n"), 0644)
		if err != nil {
			fmt.Printf("❌ Error creating file %s: %v\n", filePath, err)
		} else {
			fmt.Printf("✅ Created file %s and added cdx function\n", filePath)
		}
		return
	}

	// Read existing content
	existing, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("❌ Error reading file %s: %v\n", filePath, err)
		return
	}

	existingContent := string(existing)

	// Check if there's a commented version
	lines := strings.Split(existingContent, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#") && (strings.Contains(trimmed, "cdx()") || strings.Contains(trimmed, "function cdx")) {
			// Uncomment this line and surrounding function
			start := i
			for ; start > 0 && !strings.Contains(lines[start], "function") && !strings.Contains(lines[start], "()"); start-- {
			}
			end := i
			for ; end < len(lines) && !strings.Contains(lines[end], "}"); end++ {
			}
			
			// Uncomment all these lines
			for j := start; j <= end; j++ {
				if strings.HasPrefix(strings.TrimSpace(lines[j]), "#") {
					lines[j] = strings.TrimPrefix(strings.TrimSpace(lines[j]), "#")
				}
			}
			
			// Write back the file
			err = os.WriteFile(filePath, []byte(strings.Join(lines, "\n")), 0644)
			if err != nil {
				fmt.Printf("❌ Error updating %s: %v\n", filePath, err)
			} else {
				fmt.Printf("✅ Uncommented cdx function in %s\n", filePath)
			}
			return
		} else if strings.Contains(trimmed, "cdx()") || strings.Contains(trimmed, "function cdx") {
			fmt.Printf("ℹ️ cdx function already exists in %s\n", filePath)
			return
		}
	}

	// Append the function
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("❌ Error opening %s: %v\n", filePath, err)
		return
	}
	defer f.Close()
	
	_, err = f.WriteString("\n" + content + "\n")
	if err != nil {
		fmt.Printf("❌ Error writing to %s: %v\n", filePath, err)
	} else {
		fmt.Printf("✅ Added cdx function to %s\n", filePath)
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}