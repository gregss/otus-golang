package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in:          "not struct",
			expectedErr: errors.New("is`t type struct"),
		},
		{
			in:          App{Version: "ver10"},
			expectedErr: nil,
		},
		{
			in:          App{Version: "ver1"},
			expectedErr: errors.New("field: Version, err: len"),
		},
		{
			in: User{
				ID:     "e169cb96-223c-41ac-8fa0-af5e1614cfba",
				Name:   "name",
				Age:    25,
				Email:  "email@email.com",
				Role:   "admin",
				Phones: []string{"phone1", "phone2"},
				meta:   []byte{30, 50},
			},
			expectedErr: errors.New("field: Phones, err: len"),
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)

			if tt.expectedErr == nil {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			}
		})
	}
}
