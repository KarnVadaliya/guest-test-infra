package regression

import (
	"github.com/GoogleCloudPlatform/guest-test-infra/imagetest"
)

// Name is the name of the test package. It must match the directory name.
var Name = "regression"

const (
	vm1Name = "vm1"
	vm2Name = "vm2"
	vm1IP   = "192.168.0.2"
	vm2IP   = "192.168.0.3"
)

// TestSetup sets up the test workflow.
func TestSetup(t *imagetest.TestWorkflow) error {
	vm1, err := t.CreateTestVM("vm1")
	if err != nil {
		return err
	}
	vm1.RunTests("TestKernelStart")
	return nil
}
