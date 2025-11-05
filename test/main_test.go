package code

import (
	"code"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetPathSizeHumanTxt(t *testing.T) {
	ex1, err1 := code.GetPathSize("../testdata/example/ex1.txt", true, true, true)
	if err1 != nil {
		t.Errorf(`Error in testing of getting path size %v`, err1)
	}
	require.Equal(t, "132B", ex1, "Expected format: size\\tfilename")
}

func TestGetPathSizeHumanDocx(t *testing.T) {
	ex1, err1 := code.GetPathSize("../testdata/example/ex2.docx", true, true, true)
	if err1 != nil {
		t.Errorf(`Error in testing of getting path size %v`, err1)
	}
	require.Equal(t, "17.6KB", ex1, "Expected format: size\\tfilename")
}
