package costello

import (
	"testing"

	"github.com/markbates/pkger/pkging"
	"github.com/stretchr/testify/require"
)

func InfoTest(t *testing.T, ref *Ref, pkg pkging.Pkger) {
	r := require.New(t)

	info, err := pkg.Info("app")
	r.NoError(err)

	r.Equal(ref.Info, info)
}