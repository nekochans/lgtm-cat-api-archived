package handler

import (
	"testing"

	"github.com/nekochans/lgtm-cat-api/domain"
)

var bearerAuthorizer = NewBearerAuthorizer(&domain.JwtValidatorMock{})

func TestExtractAccessToken(t *testing.T) {
	const wantErr, noErr = true, false
	cases := []struct {
		name      string
		header    string
		want      string
		expectErr bool
	}{
		{
			name:      "Success extract access token",
			header:    "bearer access_token",
			want:      "access_token",
			expectErr: noErr,
		},
		{
			name:      "Failure bearer authorization Header not set",
			header:    "",
			want:      "",
			expectErr: wantErr,
		},
		{
			name:      "Failure invalid bearer authorization header",
			header:    "beare access_token",
			want:      "",
			expectErr: wantErr,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := bearerAuthorizer.extractAccessToken(tt.header)

			if tt.expectErr {
				if err == nil {
					t.Fatal("expected to return an error, but no error")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected err = %s", err)
				}
			}

			if got != tt.want {
				t.Errorf("\nwant\n%+v\ngot\n%+v", tt.want, got)
			}
		})
	}
}
