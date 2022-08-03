//go:build cit

package regression

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/guest-test-infra/imagetest/utils"
)

const (
	gceMTU = 1500
)

// TestKernelStart test
func TestKernelStart(t *testing.T) {
	cmd := exec.Command("ps", "x")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("ps command failed %v", err)
	}
	for _, line := range strings.Split(string(out), "\n") {
		if strings.Contains(line, fmt.Sprintf("dhclient %s", iface.Name)) {
			return
		}
	}
	t.Fatalf("failed finding dhclient process")
}
