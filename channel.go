package twitch

import "time"

type Channel struct {
	ID                           string    `json:"_id"`
	Name                         string    `json:"name"`
	DisplayName                  string    `json:"display_name"`
	Email                        string    `json:"email,omitempty"`
	Mature                       bool      `json:"mature"`
	Status                       string    `json:"status"`
	Language                     string    `json:"language"`
	BroadcasterLanguage          string    `json:"broadcaster_language"`
	Game                         string    `json:"game"`
	Partner                      bool      `json:"partner"`
	Logo                         string    `json:"logo"`
	VideoBanner                  string    `json:"video_banner"`
	ProfileBanner                string    `json:"profile_banner"`
	ProfileBannerBackgroundColor string    `json:"profile_banner_background_color"`
	URL                          string    `json:"url"`
	Views                        uint      `json:"views"`
	Followers                    uint      `json:"followers"`
	BroadcasterType              string    `json:"broadcaster_type,omitempty"`
	StreamKey                    string    `json:"stream_key,omitempty"`
	CreatedAt                    time.Time `json:"created_at"`
	UpdatedAt                    time.Time `json:"updated_at"`
}
