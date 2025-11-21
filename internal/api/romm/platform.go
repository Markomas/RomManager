package romm

type Platforms []Platform

type Platform struct {
	ID             int        `json:"id"`
	Slug           string     `json:"slug"`
	FSSlug         string     `json:"fs_slug"`
	RomCount       int        `json:"rom_count"`
	Name           string     `json:"name"`
	IgdbSlug       string     `json:"igdb_slug"`
	MobySlug       string     `json:"moby_slug"`
	HltbSlug       string     `json:"hltb_slug"`
	CustomName     string     `json:"custom_name"`
	IgdbID         int        `json:"igdb_id"`
	SgdbID         int        `json:"sgdb_id"`
	MobyID         int        `json:"moby_id"`
	LaunchboxID    int        `json:"launchbox_id"`
	SsID           int        `json:"ss_id"`
	RaID           int        `json:"ra_id"`
	HasheousID     int        `json:"hasheous_id"`
	TgdbID         int        `json:"tgdb_id"`
	FlashpointID   int        `json:"flashpoint_id"`
	Category       string     `json:"category"`
	Generation     int        `json:"generation"`
	FamilyName     string     `json:"family_name"`
	FamilySlug     string     `json:"family_slug"`
	URL            string     `json:"url"`
	URLLogo        string     `json:"url_logo"`
	Firmware       []Firmware `json:"firmware"`
	AspectRatio    string     `json:"aspect_ratio"`
	CreatedAt      string     `json:"created_at"`
	UpdatedAt      string     `json:"updated_at"`
	FSSizeBytes    int64      `json:"fs_size_bytes"`
	IsUnidentified bool       `json:"is_unidentified"`
	IsIdentified   bool       `json:"is_identified"`
	MissingFromFS  bool       `json:"missing_from_fs"`
	DisplayName    string     `json:"display_name"`
}
