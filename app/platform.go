package main

import (
	"errors"
	"os/exec"
	"strings"
)

// openTarget opens a URL or folder in the Windows host from inside WSL2:
// URLs in the default Windows browser, folders in Explorer. It falls back to
// xdg-open when the Windows interop tools aren't available.
func openTarget(target string) error {
	isURL := strings.Contains(target, "://")
	if p, err := exec.LookPath("wslview"); err == nil {
		return exec.Command(p, target).Run()
	}
	if p, err := exec.LookPath("explorer.exe"); err == nil {
		arg := target
		if !isURL {
			if wp, err := exec.LookPath("wslpath"); err == nil {
				if out, err := exec.Command(wp, "-w", target).Output(); err == nil {
					arg = strings.TrimSpace(string(out))
				}
			}
		}
		// explorer.exe reports exit status 1 even when it succeeded.
		_ = exec.Command(p, arg).Run()
		return nil
	}
	if p, err := exec.LookPath("xdg-open"); err == nil {
		return exec.Command(p, target).Run()
	}
	return errors.New("no way to open " + target + " (install wslu, or open it yourself)")
}
