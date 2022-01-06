#!/bin/bash

wget https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz
wget https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz.md5

md5sum -c ffmpeg-release-amd64-static.tar.xz.md5

tar xvf ffmpeg-release-amd64-static.tar.xz --transform 's:^[^/]*:ffmpeg-release:'

mkdir -p layer/bin

cp ffmpeg-release/ffmpeg layer/bin/

rm -rf ffmpeg ffmpeg-release
rm ffmpeg-release-amd64-static.tar.xz ffmpeg-release-amd64-static.tar.xz.md5