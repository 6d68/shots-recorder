package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/modfy/fluent-ffmpeg"
	"os"
	"path"
)

type AviToMp4Converter interface {
	ToMp4(in string) (string, error)
}

type Converter struct {
	FfMpegPath string
}

func (c Converter) ToMp4(path string) (string, error) {
	_, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return "", errors.New("file not exists")
	}

	mp4 := setExtensions(path, ".mp4")
	_, err = os.Create(mp4)
	if err != nil {
		return "", fmt.Errorf("error creating file for mp4 conversion: %w", err)
	}

	buf := &bytes.Buffer{}
	err = fluentffmpeg.NewCommand(c.FfMpegPath).
		InputPath(path).
		OutputFormat("mp4").
		OutputPath(mp4).
		Overwrite(true).
		OutputLogs(buf).
		Run()

	if err != nil {
		return "", fmt.Errorf("failed to convert %s to %s due to %s: %w ", path, mp4, buf.String(), err)
	}

	return mp4, nil
}

func setExtensions(fileName string, newExt string) string {
	ext := path.Ext(fileName)
	return fileName[0:len(fileName)-len(ext)] + newExt
}
