package authn

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateUsername(t *testing.T) {
	testCases := []struct {
		name        string
		input       string
		expectedErr error
	}{
		{
			name:        "happy__simple",
			input:       "test username",
			expectedErr: nil,
		},
		{
			name:        "error__empty_string",
			input:       "",
			expectedErr: fmt.Errorf("empty string"),
		},
		{
			name:        "error__non-ascii",
			input:       "❤️",
			expectedErr: fmt.Errorf("contains non-ascii characters"),
		},
		{
			name:        "error__non-printable_1",
			input:       "\x00",
			expectedErr: fmt.Errorf("contains non-ascii characters"),
		},
		{
			name:        "error__non-printable_2",
			input:       "\n",
			expectedErr: fmt.Errorf("contains non-ascii characters"),
		},
		{
			name:        "error__non-printable_3",
			input:       "\x7F", // del
			expectedErr: fmt.Errorf("contains non-ascii characters"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			err := ValidateUsername(tc.input)
			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}

		})
	}

}
