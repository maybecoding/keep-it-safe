package pwd

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPwd(t *testing.T) {
	type args struct {
		pwd      string
		pwdCheck string
		cost     int
	}

	type want struct {
		ok  bool
		err bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{"Test #1", args{"123", "123", bcrypt.DefaultCost}, want{true, true}},
		{"Test #2", args{"123", "12", bcrypt.DefaultCost}, want{false, true}},
		{"Test #3", args{"SuperComplexPassword!23u348u3#*#@#", "SuperComplexPassword!23u348u3#*#@#", bcrypt.DefaultCost}, want{true, true}},
		{"Test #4", args{"SuperComplexPassword!23u348u3#*#@#", "SuperComplexPassword!23u348u3#*#@", bcrypt.DefaultCost}, want{false, true}},
		{"Test #5", args{"", "", bcrypt.DefaultCost}, want{true, true}},

		{"Test #6", args{"SuperComplexPassword!23u348u3#*#@#", "SuperComplexPassword!23u348u3#*#@#", bcrypt.MinCost}, want{true, true}},
		{"Test #8", args{"SuperComplexPassword!23u348u3#*#@#", "SuperComplexPassword!23u348u3#*#@#", bcrypt.MaxCost + 1}, want{false, false}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHash, err := Hash([]byte(tt.args.pwd), tt.args.cost)
			if tt.want.err {
				require.NoError(t, err)
				require.NotEqual(t, gotHash, []byte{})

				got := Check([]byte(tt.args.pwdCheck), gotHash)
				require.Equal(t, got, tt.want.ok)
			} else {
				require.Error(t, err, "Must be error")
			}
		})
	}
}
