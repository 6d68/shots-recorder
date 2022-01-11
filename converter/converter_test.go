package main

import (
	"testing"
)

func TestConverter_ToMp4(t *testing.T) {
	type args struct {
		in string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		wantMp4 string
	}{
		{
			name: "Returns with error when neither source nor target path is specified",
			args: args{
				in: "",
			},
			wantErr: true,
			wantMp4: "",
		},
		{
			name: "Creates MP4 from AVI",
			args: args{
				in: "./test_data/recording.avi",
			},
			wantErr: false,
			wantMp4: "./test_data/recording.mp4",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Converter{
				FfMpegPath: "/usr/bin/ffmpeg",
			}

			if mp4, err := c.ToMp4(tt.args.in); (err != nil) != tt.wantErr {
				t.Errorf("ToMp4() error = %v, wantErr %v | mp4 = %v, wantMp4 = %v", err, tt.wantErr, mp4, tt.wantMp4)
			}
		})
	}
}
