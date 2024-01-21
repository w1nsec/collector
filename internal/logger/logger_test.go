package logger

import (
	"github.com/rs/zerolog"
	"testing"
)

func Test_selectLevel(t *testing.T) {

	tests := []struct {
		name  string
		level string
		want  zerolog.Level
	}{
		{
			name:  "Test INFO 1",
			level: "INF",
			want:  zerolog.InfoLevel,
		},
		{
			name:  "Test INFO 2",
			level: "INFO",
			want:  zerolog.InfoLevel,
		},
		{
			name:  "Test DBG 1",
			level: "DBG",
			want:  zerolog.DebugLevel,
		},
		{
			name:  "Test DBG 2",
			level: "DEBUG",
			want:  zerolog.DebugLevel,
		},
		{
			name:  "Test WARN 1",
			level: "WRN",
			want:  zerolog.WarnLevel,
		},
		{
			name:  "Test WARN 2",
			level: "WARNING",
			want:  zerolog.WarnLevel,
		},
		{
			name:  "Test WARN 3",
			level: "WARN",
			want:  zerolog.WarnLevel,
		},
		{
			name:  "Test ERR 1",
			level: "ERR",
			want:  zerolog.ErrorLevel,
		},
		{
			name:  "Test ERR 2",
			level: "ERROR",
			want:  zerolog.ErrorLevel,
		},
		{
			name:  "Test OTHER ",
			level: "some other level",
			want:  zerolog.DebugLevel,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := selectLevel(tt.level); got != tt.want {
				t.Errorf("selectLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}
