package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSourceRTMP(t *testing.T) {
	for _, source := range []string{
		"videoaudio",
		"video",
	} {
		t.Run(source, func(t *testing.T) {
			switch source {
			case "videoaudio", "video":
				cnt1, err := newContainer("nginx-rtmp", "rtmpserver", []string{})
				require.NoError(t, err)
				defer cnt1.close()

				input := "emptyvideoaudio.ts"
				if source == "video" {
					input = "emptyvideo.ts"
				}

				cnt2, err := newContainer("ffmpeg", "source", []string{
					"-re",
					"-stream_loop", "-1",
					"-i", input,
					"-c", "copy",
					"-f", "flv",
					"rtmp://" + cnt1.ip() + "/stream/test",
				})
				require.NoError(t, err)
				defer cnt2.close()

				time.Sleep(1 * time.Second)

				p, ok := testProgram("paths:\n" +
					"  proxied:\n" +
					"    source: rtmp://" + cnt1.ip() + "/stream/test\n" +
					"    sourceOnDemand: yes\n")
				require.Equal(t, true, ok)
				defer p.close()
			}

			time.Sleep(1 * time.Second)

			cnt3, err := newContainer("ffmpeg", "dest", []string{
				"-rtsp_transport", "udp",
				"-i", "rtsp://" + ownDockerIP + ":8554/proxied",
				"-vframes", "1",
				"-f", "image2",
				"-y", "/dev/null",
			})
			require.NoError(t, err)
			defer cnt3.close()
			require.Equal(t, 0, cnt3.wait())
		})
	}
}
