package reboot

import (
	"github.com/juju/1.25-upgrade/juju1/api/base/testing"
)

// PatchFacadeCall patches the State's facade such that
// FacadeCall method calls are diverted to the provided
// function.
func PatchFacadeCall(p testing.Patcher, st State, f func(request string, params, response interface{}) error) {
	st0 := st.(*state) // *state is the only implementation of State.
	testing.PatchFacadeCall(p, &st0.facade, f)
}
