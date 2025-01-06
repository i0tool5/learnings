Convert MP4 file into WEBM

```sh
gst-launch-1.0 -v filesrc location="bar.mp4" ! qtdemux name=demux demux.video_0 ! queue ! decodebin ! videoconvert ! videoscale ! vp8enc ! webmmux name="mux" ! filesink location='foo.webm' demux.audio_0 ! queue ! decodebin ! audioconvert ! audioresample ! vorbisenc ! mux.
```
