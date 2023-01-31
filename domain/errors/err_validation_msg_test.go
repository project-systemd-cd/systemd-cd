package errors

import (
	"errors"
	"testing"
)

func TestErrValidationMsg_As(t *testing.T) {
	type fields struct {
		Message string
	}
	type args struct {
		t any
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "true",
			fields: fields{},
			args: args{
				t: ErrValidationMsg{},
			},
			want: true,
		},
		{
			name:   "true",
			fields: fields{},
			args: args{
				t: ErrValidation{},
			},
			want: true,
		},
		{
			name:   "true",
			fields: fields{},
			args: args{
				t: &ErrValidationMsg{},
			},
			want: true,
		},
		{
			name:   "true",
			fields: fields{},
			args: args{
				t: &ErrValidation{},
			},
			want: true,
		},
		{
			name:   "true",
			fields: fields{},
			args: args{
				t: func() any { e := &ErrValidationMsg{}; return &e }(),
			},
			want: true,
		},
		{
			name:   "true",
			fields: fields{},
			args: args{
				t: func() any { e := &ErrValidation{}; return &e }(),
			},
			want: true,
		},
		{
			name:   "false",
			fields: fields{},
			args: args{
				t: &ErrNotFound{},
			},
			want: false,
		},
		{
			name:   "false",
			fields: fields{},
			args: args{
				t: func() any { e := errors.New(""); return &e },
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := ErrValidationMsg{
				Msg: tt.fields.Message,
			}
			if got := e.As(tt.args.t); got != tt.want {
				t.Errorf("ErrValidationMsg.As() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrValidationMsg_Error(t *testing.T) {
	type fields struct {
		Message string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "ok",
			fields: fields{
				Message: "hoge",
			},
			want: "hoge",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := ErrValidationMsg{
				Msg: tt.fields.Message,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("ErrValidationMsg.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
