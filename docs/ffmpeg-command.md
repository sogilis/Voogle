# FFMPEG

When you have to process a video, `ffmpeg` is the only tool you need and want.
This document contains several command that can be useful to process videos for streaming purpose.

Ffmpeg is insanely complete, but also really hard to use.

(?) -> It might need verification

## HLS

Watching videos on internet is not as simple as watching pictures, especially if you want to optimise your users bandwidth and your resources.

To solve this issue, two video format were invented: DASH and HLS. Both formats are optimise for the streaming of videos but it seems that HLS is more popular, support more devices and was originally supported by Apple.

### HLS in details

To transform a video from for example AVI to HLS (aside from encoding), we split the original video in small segment of a few seconds (and sometimes with the sound in separates files). 
Then we write a master file where we list all the segment, their encoding, time and many more. With this, the user can ask only the required segment and not load the whole video then play it.

```text
#EXTM3U
#EXT-X-VERSION:3
#EXT-X-TARGETDURATION:6
#EXT-X-MEDIA-SEQUENCE:0
#EXTINF:6.006000,
part0.ts
#EXTINF:6.006000,
part1.ts
#EXTINF:6.006000,
part2.ts
#EXTINF:6.006000,
part3.ts
#EXTINF:4.713356,
part4.ts
#EXT-X-ENDLIST
```
*Master that list the segments (Probably not 100% accurate)*

But if you want to make you video more accessible, your original video in 4k might not work for everyone. 
For this you can convert you video to lower resolution.

At the end, you will have a master manifest that reference all the resolution and that link to a master for each resolution you have.

```text
#EXTM3U
#EXT-X-VERSION:3
#EXT-X-STREAM-INF:BANDWIDTH=1240800,RESOLUTION=640x480,CODECS="avc1.64001e,mp4a.40.2"
v0/part_index.m3u8

#EXT-X-STREAM-INF:BANDWIDTH=2340800,RESOLUTION=1280x720,CODECS="avc1.64001f,mp4a.40.2"
v1/part_index.m3u8

#EXT-X-STREAM-INF:BANDWIDTH=4540800,RESOLUTION=1920x1080,CODECS="avc1.640032,mp4a.40.2"
v2/part_index.m3u8

#EXT-X-STREAM-INF:BANDWIDTH=8940800,RESOLUTION=3840x2160,CODECS="avc1.640033,mp4a.40.2"
v3/part_index.m3u8
```
*A master of master (with 4 resolutions (480p, 720p, 1080, 4k))*

And the folder you generate will look like this:

```bash
$tree seagle

seagle
├── master.m3u8
├── v0
│    ├── segment0.ts
│    ├── ...
│    ├── segment7.ts
│    └── segment_index.m3u8
└── v1
    ├── segment0.ts
    ├──  ...
    ├── segment7.ts
    └── segment_index.m3u8
```
*Two are missing to be clearer. It will be the same structure with a v2 and v3*

### Watching flow

When you open a video, the player will do something like this:
```text
player -> master.m3u8
player -> v0/segment_index.m3u8
player -> v0/segment0.ts
...
*Quality change to 4k*
player -> v3/segment_index.m3u8
player -> v3/segmentn.ts
...
```

## FFMPEG convert Video to HLS
This command can be improved a lot. Here is a few improvement idea:
* Extract the sound from each part (This way, each part will be lighter)
* Bitrate might need to tweaked for each resolution (and adapted to the type of video (higher bitrate for a video with a lot of movement))
* Each quality level can be named more clearly (Here, it will be `480p` -> `v0`, `720p` -> `v1`, ... )
* Change the preset to be faster
* Use graphical acceleration

```bash
ffmpeg -y -i <filepath> \
              -pix_fmt yuv420p \ # Encoding pixel formation
              -vcodec libx264 \ # Encoding video codec
              -preset slow \ # It define desired the encoding speed. The slower equals better quality
              -g 48 -sc_threshold 0 \ 
              -map 0:0 -map 0:1 -map 0:0 -map 0:1 -map 0:0 -map 0:1 -map 0:0 -map 0:1 \ # For each resolution, it maps the sound with it
              -s:v:0 640x480 -c:v:0 libx264 -b:v:0 1000k \ # `480p` encoded with `libx264` at a bitrate of `1000k`
              -s:v:1 1280x720 -c:v:1 libx264 -b:v:1 2000k  \ # `720p` encoded with `libx264` at a bitrate of `2000k`
              -s:v:2 1920x1080 -c:v:2 libx264 -b:v:2 4000k  \ # `1080` encoded with `libx264` at a bitrate of `4000k`
              -s:v:3 3840x2160 -c:v:3 libx264 -b:v:3 8000k  \ # `2160` encoded with `libx264` at a bitrate of `8000k`
              -c:a aac -b:a 128k -ac 2 \  # Sound encoding
              -var_stream_map "v:0,a:0 v:1,a:1 v:2,a:2 v:3,a:3" \ # For each resolution with add them to the master (?)
              -master_pl_name master.m3u8 \ # Master file name
              -f hls -hls_time 6 -hls_list_size 0 \
              -hls_segment_filename "v%v/segment%d.ts" \ # Determine how are save each segment, here resolution/part(count).ts
              v%v/segment_index.m3u8 # each quality level are registered in a sub master file
```
*(If you want to use this command, remove the comments, sorry)*

## FFMPEG extract resolution
It returns the resolution of the video with a pattern like this: "WidthxHeight" (1024x700)

```bash
   ffprobe -v error -select_streams v:0 -show_entries stream=width,height -of csv=s=x:p=0 <filepath>
```

## FFMPEG screenshot
It extracts one frame at the given time code and write it to the provided path

```bash
ffmpeg -i <filepath> -ss 00:00:01.000 -vframes 1 miniature.jpeg
```
