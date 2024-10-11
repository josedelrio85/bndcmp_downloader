package bandcamp

import (
	"strings"

	"github.com/josedelrio85/bndcmp_downloader/internal/model"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type TrAlbum struct {
	ForTheCurious              string      `json:"for the curious"`
	Current                    Current     `json:"current"`
	PreorderCount              any         `json:"preorder_count"`
	HasAudio                   bool        `json:"hasAudio"`
	ArtID                      int64       `json:"art_id"`
	Packages                   []any       `json:"packages"`
	DefaultPrice               float64     `json:"defaultPrice"`
	FreeDownloadPage           any         `json:"freeDownloadPage"`
	Free                       int64       `json:"FREE"`
	Paid                       int64       `json:"PAID"`
	Artist                     string      `json:"artist"`
	ItemType                   string      `json:"item_type"`
	ID                         int64       `json:"id"`
	LastSubscriptionItem       any         `json:"last_subscription_item"`
	HasDiscounts               bool        `json:"has_discounts"`
	IsBonus                    any         `json:"is_bonus"`
	PlayCapData                PlayCapData `json:"play_cap_data"`
	IsPurchased                any         `json:"is_purchased"`
	ItemsPurchased             any         `json:"items_purchased"`
	IsPrivateStream            any         `json:"is_private_stream"`
	IsBandMember               any         `json:"is_band_member"`
	LicensedVersionIds         any         `json:"licensed_version_ids"`
	PackageAssociatedLicenseID any         `json:"package_associated_license_id"`
	HasVideo                   any         `json:"has_video"`
	TralbumSubscriberOnly      bool        `json:"tralbum_subscriber_only"`
	AlbumIsPreorder            bool        `json:"album_is_preorder"`
	AlbumReleaseDate           string      `json:"album_release_date"`
	Trackinfo                  []TrackInfo `json:"trackinfo"`
	PlayingFrom                string      `json:"playing_from"`
	AlbumURL                   string      `json:"album_url"`
	AlbumUpsellURL             string      `json:"album_upsell_url"`
	URL                        string      `json:"url"`
}

type Current struct {
	Audit               int64   `json:"audit"`
	Title               string  `json:"title"`
	NewDate             string  `json:"new_date"`
	ModDate             string  `json:"mod_date"`
	PublishDate         string  `json:"publish_date"`
	Private             any     `json:"private"`
	Killed              any     `json:"killed"`
	DownloadPref        int64   `json:"download_pref"`
	RequireEmail        any     `json:"require_email"`
	IsSetPrice          any     `json:"is_set_price"`
	SetPrice            float64 `json:"set_price"`
	MinimumPrice        float64 `json:"minimum_price"`
	MinimumPriceNonzero float64 `json:"minimum_price_nonzero"`
	RequireEmail0       any     `json:"require_email_0"`
	Artist              any     `json:"artist"`
	About               any     `json:"about"`
	Credits             any     `json:"credits"`
	AutoRepriced        any     `json:"auto_repriced"`
	NewDescFormat       int64   `json:"new_desc_format"`
	BandID              int64   `json:"band_id"`
	SellingBandID       int64   `json:"selling_band_id"`
	ArtID               any     `json:"art_id"`
	DownloadDescID      any     `json:"download_desc_id"`
	TrackNumber         int64   `json:"track_number"`
	ReleaseDate         any     `json:"release_date"`
	FileName            any     `json:"file_name"`
	Lyrics              any     `json:"lyrics"`
	AlbumID             int64   `json:"album_id"`
	EncodingsID         int64   `json:"encodings_id"`
	PendingEncodingsID  any     `json:"pending_encodings_id"`
	LicenseType         int64   `json:"license_type"`
	Isrc                string  `json:"isrc"`
	PreorderDownload    any     `json:"preorder_download"`
	Streaming           int64   `json:"streaming"`
	ID                  int64   `json:"id"`
	Type                string  `json:"type"`
}

type TrackInfo struct {
	ID                int64   `json:"id"`
	TrackID           int64   `json:"track_id"`
	File              File    `json:"file"`
	Artist            any     `json:"artist"`
	Title             string  `json:"title"`
	EncodingsID       int64   `json:"encodings_id"`
	LicenseType       int64   `json:"license_type"`
	Private           any     `json:"private"`
	TrackNum          int64   `json:"track_num"`
	AlbumPreorder     bool    `json:"album_preorder"`
	UnreleasedTrack   bool    `json:"unreleased_track"`
	TitleLink         string  `json:"title_link"`
	HasLyrics         bool    `json:"has_lyrics"`
	HasInfo           bool    `json:"has_info"`
	Streaming         int64   `json:"streaming"`
	IsDownloadable    bool    `json:"is_downloadable"`
	HasFreeDownload   any     `json:"has_free_download"`
	FreeAlbumDownload bool    `json:"free_album_download"`
	Duration          float64 `json:"duration"`
	Lyrics            any     `json:"lyrics"`
	SizeofLyrics      int64   `json:"sizeof_lyrics"`
	IsDraft           bool    `json:"is_draft"`
	VideoSourceType   any     `json:"video_source_type"`
	VideoSourceID     any     `json:"video_source_id"`
	VideoMobileURL    any     `json:"video_mobile_url"`
	VideoPosterURL    any     `json:"video_poster_url"`
	VideoID           any     `json:"video_id"`
	VideoCaption      any     `json:"video_caption"`
	VideoFeatured     any     `json:"video_featured"`
	AltLink           any     `json:"alt_link"`
	EncodingError     any     `json:"encoding_error"`
	EncodingPending   any     `json:"encoding_pending"`
	PlayCount         int64   `json:"play_count"`
	IsCapped          bool    `json:"is_capped"`
	TrackLicenseID    any     `json:"track_license_id"`
}

type File struct {
	Mp3128 string `json:"mp3-128"`
}

type PlayCapData struct {
	StreamingLimitsEnabled bool  `json:"streaming_limits_enabled"`
	StreamingLimit         int64 `json:"streaming_limit"`
}

func (t *TrAlbum) ToTrack() *model.Track {
	if t == nil {
		return nil
	}

	track := model.Track{
		Title:       t.Current.Title,
		TrackNumber: t.Current.TrackNumber,
		Artist:      t.Artist,
		Album:       t.getAlbumName(),
		URL:         t.URL,
	}

	if len(t.Trackinfo) > 0 {
		track.DownloadURL = t.Trackinfo[0].File.Mp3128
	}

	return &track
}

func (t *TrAlbum) getAlbumName() *string {
	if t == nil {
		return nil
	}

	album := strings.Replace(t.AlbumURL, "/album/", "", -1)
	album = strings.Replace(album, "-", " ", -1)
	words := strings.Fields(album)
	caser := cases.Title(language.Und)
	for i, word := range words {
		words[i] = caser.String(word)
	}
	album = strings.Join(words, " ")
	return &album
}
