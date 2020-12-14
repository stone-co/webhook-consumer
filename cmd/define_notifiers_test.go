package main

import (
	"reflect"
	"testing"
)

func Test_extractNotifiersFromConfig(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    []string
		wantErr bool
	}{
		{
			name: "Valid config with one notifier",
			args: "stdout",
			want: []string{"stdout"},
		},
		{
			name: "Valid config with two notifiers",
			args: "stdout;proxy",
			want: []string{"stdout", "proxy"},
		},
		{
			name: "Valid config with three notifiers",
			args: "stdout;proxy;redis",
			want: []string{"stdout", "proxy", "redis"},
		},
		{
			name: "Config is case insensitive",
			args: "stDOut;PROxy",
			want: []string{"stdout", "proxy"},
		},
		{
			name: "Spaces all trimmed when necessary",
			args: "      stdout   ;     proxy     ",
			want: []string{"stdout", "proxy"},
		},
		{
			name: "Extra ';' are ignored",
			args: ";;      stdout   ;    ; proxy  ;   ",
			want: []string{"stdout", "proxy"},
		},
		{
			name:    "Empty config must fail",
			args:    "",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid notifier must fail",
			args:    "xpto;stdout",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Duplicate notifier must fail",
			args:    "stdout;stdout",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractNotifiersFromConfig(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractNotifiersFromConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractNotifiersFromConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
