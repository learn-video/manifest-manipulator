package filter

import "github.com/grafov/m3u8"

func FilterBandwidth(mp *m3u8.MasterPlaylist, min uint32) *m3u8.MasterPlaylist {
	variants := make([]*m3u8.Variant, 0)
	for _, variant := range mp.Variants {
		if variant.Bandwidth >= min {
			variants = append(variants, variant)
		}
	}
	mp.Variants = variants
	return mp
}
