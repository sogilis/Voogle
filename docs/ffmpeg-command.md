# FFMPEG

When you have to process a video, `ffmpeg` is the only tool you need and want.
This document contains several command that can be useful to process videos for streaming purpose.

Ffmpeg is insanely complete, but also really hard to use.

(?) -> It might need verification

## HLS

regarder des vidéos sur internet n'est pas aussi simple qu'afficher une image, surtout si on veut sauvegarder nos ressources et tenir compte de la connexion des utilisateurs.

Pour ça, il y a au moins deux formats de vidéo qui ont été conçus récemment: le Dash et le HLS. Le HLS est devenu le format le plus courant et le plus supporté et il a été conçu et originellement porté par Apple.

Mais vous êtes probablement en train de vous demander ce qu'est ce format.

Pour le HLS, on vient découper la vidéo originale en petit segment de quelques secondes (on peut faire pareil avec le son, mais la taille étant généralement plus petite que la vidéo, ce n'est pas forcément nécessaire).

Pour gérer la partie adaptatif, on convertit la vidéo originale vers des résolutions inférieures. Cette étape n'est pas obligatoire pour le streaming de la vidéo, c'est juste du confort pour l'utilisateur.

Une fois les segments générés, on vient généré un manifeste principal et un manifeste par sous résolution si on veut être adaptatif.

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
              -hls_segment_filename "v%v/part%d.ts" \ # Determine how are save each segment, here resolution/part(count).ts
              v%v/part_index.m3u8 # each quality level are registered in a sub master file
```

The tree will look like this:
```bash
$tree seagle

seagle
├── master.m3u8
├── v0
│    ├── part0.ts
│    ├── ...
│    ├── part7.ts
│    └── part_index.m3u8
└── v1
    ├── part0.ts
    ├──  ...
    ├── part7.ts
    └── part_index.m3u8

```

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
