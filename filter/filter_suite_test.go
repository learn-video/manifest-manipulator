package filter_test

import (
	"bytes"
	"testing"

	"github.com/mauricioabreu/video-manifest/filter"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var masterPlaylistData = `
#EXTM3U
#EXT-X-VERSION:4
#EXT-X-MEDIA:TYPE=AUDIO,GROUP-ID="audio-aach-96",LANGUAGE="en",NAME="English",DEFAULT=YES,AUTOSELECT=YES,CHANNELS="2"
#EXT-X-STREAM-INF:BANDWIDTH=600000,AVERAGE-BANDWIDTH=600000,CODECS="mp4a.40.5,avc1.64001F",RESOLUTION=384x216,FRAME-RATE=30,AUDIO="audio-aach-96",CLOSED-CAPTIONS=NONE
variant-audio_1=96000-video=249984.m3u8
#EXT-X-STREAM-INF:BANDWIDTH=800000,AVERAGE-BANDWIDTH=800000,CODECS="mp4a.40.5,avc1.64001F",RESOLUTION=768x432,FRAME-RATE=30,AUDIO="audio-aach-96",CLOSED-CAPTIONS=NONE
variant-audio_1=96000-video=1320960.m3u8
#EXT-X-STREAM-INF:BANDWIDTH=1500000,AVERAGE-BANDWIDTH=1500000,CODECS="mp4a.40.5,avc1.64001F",RESOLUTION=1280x720,FRAME-RATE=60,AUDIO="audio-aach-96",CLOSED-CAPTIONS=NONE
variant-audio_1=96000-video=3092992.m3u8
#EXT-X-STREAM-INF:BANDWIDTH=2000000,AVERAGE-BANDWIDTH=2000000,CODECS="mp4a.40.5,avc1.640029",RESOLUTION=1920x1080,FRAME-RATE=60,AUDIO="audio-aach-96",CLOSED-CAPTIONS=NONE
variant-audio_1=96000-video=4686976.m3u8
#EXT-X-I-FRAME-STREAM-INF:BANDWIDTH=37000,CODECS="avc1.64001F",RESOLUTION=384x216,URI="keyframes/variant-video=249984.m3u8"
#EXT-X-I-FRAME-STREAM-INF:BANDWIDTH=193000,CODECS="avc1.64001F",RESOLUTION=768x432,URI="keyframes/variant-video=1320960.m3u8"
#EXT-X-I-FRAME-STREAM-INF:BANDWIDTH=296000,CODECS="avc1.64001F",RESOLUTION=1280x720,URI="keyframes/variant-video=2029952.m3u8"
#EXT-X-I-FRAME-STREAM-INF:BANDWIDTH=684000,CODECS="avc1.640029",RESOLUTION=1920x1080,URI="keyframes/variant-video=4686976.m3u8"
`

func TestFilter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Filter Suite")
}

var _ = Describe("Filter variants by bandwidth", func() {
	When("min is set", func() {
		It("collects all variants with bandwidth greater than min", func() {
			p, err := filter.NewMasterPlaylist(*bytes.NewBufferString(masterPlaylistData))
			Expect(err).ToNot(HaveOccurred())
			Expect(len(p.Playlist.Variants)).To(Equal(8))

			p.FilterBandwidth(filter.BandwidthFilter{
				Min: 800000,
			})

			Expect(len(p.Playlist.Variants)).To(Equal(3))
		})
	})
	When("max is set", func() {
		It("collects all variants with bandwidth lower than max", func() {
			p, err := filter.NewMasterPlaylist(*bytes.NewBufferString(masterPlaylistData))
			Expect(err).ToNot(HaveOccurred())
			Expect(len(p.Playlist.Variants)).To(Equal(8))

			p.FilterBandwidth(filter.BandwidthFilter{
				Max: 1500000,
			})

			Expect(len(p.Playlist.Variants)).To(Equal(7))
		})
	})
	When("min and max are set", func() {
		It("collects all variants with bandwidth between min and max", func() {
			p, err := filter.NewMasterPlaylist(*bytes.NewBufferString(masterPlaylistData))
			Expect(err).ToNot(HaveOccurred())
			Expect(len(p.Playlist.Variants)).To(Equal(8))

			p.FilterBandwidth(filter.BandwidthFilter{
				Min: 800000,
				Max: 1500000,
			})

			Expect(len(p.Playlist.Variants)).To(Equal(2))
		})
	})
})
