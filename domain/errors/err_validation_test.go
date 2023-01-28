package errors

import (
	"errors"
	"testing"
)

func TestErrValidation_As(t *testing.T) {
	type fields struct {
		Property    string
		Given       *string
		Description string
	}
	type args struct {
		t interface{}
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
				t: ErrValidation{},
			},
			want: true,
		},
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
				t: &ErrValidation{},
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
				t: func() any { e := &ErrValidation{}; return &e }(),
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
			e := &ErrValidation{
				Property:    tt.fields.Property,
				Given:       tt.fields.Given,
				Description: tt.fields.Description,
			}
			if got := e.As(tt.args.t); got != tt.want {
				t.Errorf("ErrValidation.As() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrValidation_Error(t *testing.T) {
	type fields struct {
		Property    string
		Given       *string
		Description string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "nil",
			fields: fields{
				Property:    "Name",
				Given:       nil,
				Description: "given nil",
			},
			want: "the property \"Name\" given 'nil' given nil.",
		},
		{
			name: "non-nil pointer",
			fields: fields{
				Property:    "Name",
				Given:       new(string),
				Description: "given non-nil pointer",
			},
			want: "the property \"Name\" given \"\" given non-nil pointer.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ErrValidation{
				Property:    tt.fields.Property,
				Given:       tt.fields.Given,
				Description: tt.fields.Description,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("ErrValidation.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
