package model

import (
	"reflect"
	"testing"
)

func TestNewFailedResult(t *testing.T) {
	type args struct {
		message string
	}
	tests := []struct {
		name string
		args args
		want Result
	}{
		{name: "not ok", args: args{"message"}, want: Result{
			Success: false,
			Message: "message",
			Data:    nil,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFailedResult(tt.args.message); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFailedResult() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewSuccessResult(t *testing.T) {
	type args struct {
		data any
	}
	tests := []struct {
		name string
		args args
		want Result
	}{
		{name: "ok", args: args{1}, want: Result{
			Success: true,
			Message: "",
			Data:    1,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSuccessResult(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSuccessResult() = %v, want %v", got, tt.want)
			}
		})
	}
}
