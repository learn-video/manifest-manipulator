package filter_test

import (
	"bytes"
	"testing"

	"github.com/mauricioabreu/manifest-manipulator/filter"
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

var variantPlaylistData = `
#EXTM3U
#EXT-X-VERSION:4
#EXT-X-MEDIA-SEQUENCE:320035356
#EXT-X-INDEPENDENT-SEGMENTS
#EXT-X-TARGETDURATION:8
#EXT-X-PROGRAM-DATE-TIME:2020-09-15T13:32:55Z
#EXTINF:5, no desc
variant-audio_1=96000-video=249984-320035684.ts
#EXTINF:5, no desc
variant-audio_1=96000-video=249984-320035685.ts
#EXTINF:5, no desc
variant-audio_1=96000-video=249984-320035686.ts
#EXTINF:5, no desc
variant-audio_1=96000-video=249984-320035687.ts
#EXTINF:4.1333, no desc
variant-audio_1=96000-video=249984-320035688.ts
#EXT-X-DATERANGE:ID="4026531847",START-DATE="2020-09-15T14:00:39.133333Z",PLANNED-DURATION=60,SCTE35-OUT=0xFC3025000000000BB800FFF01405F00000077FEFFE0AF311F0FE005265C0000101010000817C918E
#EXT-X-CUE-OUT:60
#EXT-X-PROGRAM-DATE-TIME:2020-09-15T14:00:39.133333Z
#EXTINF:5.8666, no desc
variant-audio_1=96000-video=249984-320035689.ts
#EXTINF:5, no desc
variant-audio_1=96000-video=249984-320035690.ts
#EXTINF:5, no desc
variant-audio_1=96000-video=249984-320035691.ts
#EXTINF:5, no desc
variant-audio_1=96000-video=249984-320035692.ts
#EXTINF:5, no desc
variant-audio_1=96000-video=249984-320035693.ts
#EXTINF:5, no desc
variant-audio_1=96000-video=249984-320035694.ts
#EXTINF:5, no desc
variant-audio_1=96000-video=249984-320035695.ts
#EXTINF:5, no desc
variant-audio_1=96000-video=249984-320035696.ts
#EXTINF:5, no desc
variant-audio_1=96000-video=249984-320035697.ts
#EXTINF:5, no desc
variant-audio_1=96000-video=249984-320035698.ts
#EXTINF:5, no desc
variant-audio_1=96000-video=249984-320035699.ts
#EXTINF:4.1333, no desc
variant-audio_1=96000-video=249984-320035700.ts
#EXT-X-CUE-IN
#EXT-X-PROGRAM-DATE-TIME:2020-09-15T14:01:39.133333Z
#EXTINF:5.8666, no desc
variant-audio_1=96000-video=249984-320035701.ts
#EXTINF:5, no desc
variant-audio_1=96000-video=249984-320035702.ts
#EXTINF:5, no desc
variant-audio_1=96000-video=249984-320035703.ts
`

func TestFilter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Filter Suite")
}

var _ = Describe("Filter master playlist", func() {
	Describe("Filter bandwidth", func() {
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
	Describe("Filter by frame rate", func() {
		It("collects all variants with the given frame rate", func() {
			p, err := filter.NewMasterPlaylist(*bytes.NewBufferString(masterPlaylistData))
			Expect(err).ToNot(HaveOccurred())
			Expect(len(p.Playlist.Variants)).To(Equal(8))

			p.FilterFrameRate(60)

			Expect(len(p.Playlist.Variants)).To(Equal(2))
		})
	})
	Describe("Set first playlist", func() {
		It("sets the selected playlist to be the first", func() {
			p, err := filter.NewMasterPlaylist(*bytes.NewBufferString(masterPlaylistData))
			Expect(err).ToNot(HaveOccurred())

			p.SetFirst(2)

			variants := p.Playlist.Variants
			Expect(len(variants)).To(Equal(8))
			Expect(int(variants[0].Bandwidth)).To(Equal(800000))
			Expect(int(variants[1].Bandwidth)).To(Equal(600000))
			Expect(int(variants[2].Bandwidth)).To(Equal(1500000))
			Expect(int(variants[3].Bandwidth)).To(Equal(2000000))
		})
	})
})

var _ = Describe("Filter variant playlist", func() {
	Describe("Filter DVR", func() {
		It("removes some segments from the variant playlist", func() {
			v, err := filter.NewVariant(*bytes.NewBufferString(variantPlaylistData))
			Expect(err).ToNot(HaveOccurred())
			Expect(int(v.Playlist.Count())).To(Equal(20))

			v.FilterDVR(20)

			Expect(int(v.Playlist.Count())).To(Equal(16))
		})
	})
})
