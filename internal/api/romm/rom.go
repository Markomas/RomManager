package romm

type Roms []Rom
type RomsResponse struct {
	Items  []Rom `json:"items"`
	Total  int   `json:"total"`
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
}

type Rom struct {
	ID                  int                `json:"id"`
	IgdbID              int                `json:"igdb_id"`
	SgdbID              int                `json:"sgdb_id"`
	MobyID              *int               `json:"moby_id"`
	SsID                int                `json:"ss_id"`
	RaID                int                `json:"ra_id"`
	LaunchboxID         *int               `json:"launchbox_id"`
	HasheousID          int                `json:"hasheous_id"`
	TgdbID              *int               `json:"tgdb_id"`
	FlashpointID        *int               `json:"flashpoint_id"`
	HltbID              *int               `json:"hltb_id"`
	GamelistID          *int               `json:"gamelist_id"`
	PlatformID          int                `json:"platform_id"`
	PlatformSlug        string             `json:"platform_slug"`
	PlatformFsSlug      string             `json:"platform_fs_slug"`
	PlatformCustomName  string             `json:"platform_custom_name"`
	PlatformDisplayName string             `json:"platform_display_name"`
	FsName              string             `json:"fs_name"`
	FsNameNoTags        string             `json:"fs_name_no_tags"`
	FsNameNoExt         string             `json:"fs_name_no_ext"`
	FsExtension         string             `json:"fs_extension"`
	FsPath              string             `json:"fs_path"`
	FsSizeBytes         int64              `json:"fs_size_bytes"`
	Name                string             `json:"name"`
	Slug                string             `json:"slug"`
	Summary             string             `json:"summary"`
	AlternativeNames    []string           `json:"alternative_names"`
	YoutubeVideoID      string             `json:"youtube_video_id"`
	Metadatum           Metadatum          `json:"metadatum"`
	IgdbMetadata        IgdbMetadata       `json:"igdb_metadata"`
	MobyMetadata        MobyMetadata       `json:"moby_metadata"`
	SsMetadata          SsMetadata         `json:"ss_metadata"`
	LaunchboxMetadata   LaunchboxMetadata  `json:"launchbox_metadata"`
	HasheousMetadata    HasheousMetadata   `json:"hasheous_metadata"`
	FlashpointMetadata  FlashpointMetadata `json:"flashpoint_metadata"`
	HltbMetadata        HltbMetadata       `json:"hltb_metadata"`
	GamelistMetadata    GamelistMetadata   `json:"gamelist_metadata"`
	PathCoverSmall      string             `json:"path_cover_small"`
	PathCoverLarge      string             `json:"path_cover_large"`
	URLCover            string             `json:"url_cover"`
	HasManual           bool               `json:"has_manual"`
	PathManual          *string            `json:"path_manual"`
	URLManual           string             `json:"url_manual"`
	IsIdentifying       bool               `json:"is_identifying"`
	IsUnidentified      bool               `json:"is_unidentified"`
	IsIdentified        bool               `json:"is_identified"`
	Revision            string             `json:"revision"`
	Regions             []string           `json:"regions"`
	Languages           []string           `json:"languages"`
	Tags                []string           `json:"tags"`
	CrcHash             string             `json:"crc_hash"`
	Md5Hash             string             `json:"md5_hash"`
	Sha1Hash            string             `json:"sha1_hash"`
	Multi               bool               `json:"multi"`
	HasSimpleSingleFile bool               `json:"has_simple_single_file"`
	HasNestedSingleFile bool               `json:"has_nested_single_file"`
	HasMultipleFiles    bool               `json:"has_multiple_files"`
	Files               []File             `json:"files"`
	FullPath            string             `json:"full_path"`
	CreatedAt           string             `json:"created_at"`
	UpdatedAt           string             `json:"updated_at"`
	MissingFromFs       bool               `json:"missing_from_fs"`
	Siblings            []Sibling          `json:"siblings"`
	RomUser             RomUser            `json:"rom_user"`
	MergedRaMetadata    MergedRaMetadata   `json:"merged_ra_metadata"`
}

type Metadatum struct {
	RomID            int      `json:"rom_id"`
	Genres           []string `json:"genres"`
	Franchises       []string `json:"franchises"`
	Collections      []string `json:"collections"`
	Companies        []string `json:"companies"`
	GameModes        []string `json:"game_modes"`
	AgeRatings       []string `json:"age_ratings"`
	FirstReleaseDate int64    `json:"first_release_date"`
	AverageRating    float32  `json:"average_rating"`
}

type IgdbMetadataAgeRating struct {
	//{"rating":"T","category":"ESRB","rating_cover_url":"https://www.igdb.com/icons/rating_icons/esrb/esrb_t.png"}
	Rating         string `json:"rating"`
	Category       string `json:"category"`
	RatingCoverURL string `json:"rating_cover_url"`
}

type IgdbMetadata struct {
	TotalRating      string                  `json:"total_rating"`
	AggregatedRating string                  `json:"aggregated_rating"`
	FirstReleaseDate int64                   `json:"first_release_date"`
	YoutubeVideoID   string                  `json:"youtube_video_id"`
	Genres           []string                `json:"genres"`
	Franchises       []string                `json:"franchises"`
	AlternativeNames []string                `json:"alternative_names"`
	Collections      []string                `json:"collections"`
	Companies        []string                `json:"companies"`
	GameModes        []string                `json:"game_modes"`
	AgeRatings       []IgdbMetadataAgeRating `json:"age_ratings"`
	Platforms        []IgdbPlatform          `json:"platforms"`
	Expansions       []interface{}           `json:"expansions"`
	Dlcs             []interface{}           `json:"dlcs"`
	Remasters        []interface{}           `json:"remasters"`
	Remakes          []interface{}           `json:"remakes"`
	ExpandedGames    []interface{}           `json:"expanded_games"`
	Ports            []interface{}           `json:"ports"`
	SimilarGames     []SimilarGame           `json:"similar_games"`
}

type AgeRating struct {
	Rating         string `json:"rating"`
	Category       string `json:"category"`
	RatingCoverURL string `json:"rating_cover_url"`
}

type IgdbPlatform struct {
	IgdbID int    `json:"igdb_id"`
	Name   string `json:"name"`
}

type SimilarGame struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	Type     string `json:"type"`
	CoverURL string `json:"cover_url"`
}

type File struct {
	ID            int     `json:"id"`
	RomID         int     `json:"rom_id"`
	FileName      string  `json:"file_name"`
	FilePath      string  `json:"file_path"`
	FileSizeBytes int64   `json:"file_size_bytes"`
	FullPath      string  `json:"full_path"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
	LastModified  string  `json:"last_modified"`
	CrcHash       string  `json:"crc_hash"`
	Md5Hash       string  `json:"md5_hash"`
	Sha1Hash      string  `json:"sha1_hash"`
	Category      *string `json:"category"`
}

type RomUser struct {
	ID              int     `json:"id"`
	UserID          int     `json:"user_id"`
	RomID           int     `json:"rom_id"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
	LastPlayed      *string `json:"last_played"`
	NoteRawMarkdown string  `json:"note_raw_markdown"`
	NoteIsPublic    bool    `json:"note_is_public"`
	IsMainSibling   bool    `json:"is_main_sibling"`
	Backlogged      bool    `json:"backlogged"`
	NowPlaying      bool    `json:"now_playing"`
	Hidden          bool    `json:"hidden"`
	Rating          int     `json:"rating"`
	Difficulty      int     `json:"difficulty"`
	Completion      int     `json:"completion"`
	Status          *string `json:"status"`
	UserUsername    string  `json:"user__username"`
}

type MobyMetadata map[string]interface{}
type SsMetadata struct {
	//{
	//        "bezel_url": null,
	//        "box2d_url": "https://neoclone.screenscraper.fr/api2/mediaJeu.php?devid=zurdi15&devpassword=xTJwoOFjOQG&softname=romm&ssid=markom&sspassword=k72f8ks4L1GuYw&systemeid=10&jeuid=46335&media=box-2D(us)",
	//        "box2d_side_url": "https://neoclone.screenscraper.fr/api2/mediaJeu.php?devid=zurdi15&devpassword=xTJwoOFjOQG&softname=romm&ssid=markom&sspassword=k72f8ks4L1GuYw&systemeid=10&jeuid=46335&media=box-2D-side(us)",
	//        "box2d_back_url": "https://neoclone.screenscraper.fr/api2/mediaJeu.php?devid=zurdi15&devpassword=xTJwoOFjOQG&softname=romm&ssid=markom&sspassword=k72f8ks4L1GuYw&systemeid=10&jeuid=46335&media=box-2D-back(us)",
	//        "box3d_url": "https://neoclone.screenscraper.fr/api2/mediaJeu.php?devid=zurdi15&devpassword=xTJwoOFjOQG&softname=romm&ssid=markom&sspassword=k72f8ks4L1GuYw&systemeid=10&jeuid=46335&media=box-3D(us)",
	//        "fanart_url": "https://neoclone.screenscraper.fr/api2/mediaJeu.php?devid=zurdi15&devpassword=xTJwoOFjOQG&softname=romm&ssid=markom&sspassword=k72f8ks4L1GuYw&systemeid=10&jeuid=46335&media=fanart",
	//        "fullbox_url": "https://neoclone.screenscraper.fr/api2/mediaJeu.php?devid=zurdi15&devpassword=xTJwoOFjOQG&softname=romm&ssid=markom&sspassword=k72f8ks4L1GuYw&systemeid=10&jeuid=46335&media=box-texture(us)",
	//        "logo_url": null,
	//        "manual_url": null,
	//        "marquee_url": "https://neoclone.screenscraper.fr/api2/mediaJeu.php?devid=zurdi15&devpassword=xTJwoOFjOQG&softname=romm&ssid=markom&sspassword=k72f8ks4L1GuYw&systemeid=10&jeuid=46335&media=screenmarquee(wor)",
	//        "miximage_url": "https://neoclone.screenscraper.fr/api2/mediaJeu.php?devid=zurdi15&devpassword=xTJwoOFjOQG&softname=romm&ssid=markom&sspassword=k72f8ks4L1GuYw&systemeid=10&jeuid=46335&media=mixrbv1(us)",
	//        "physical_url": "https://neoclone.screenscraper.fr/api2/mediaJeu.php?devid=zurdi15&devpassword=xTJwoOFjOQG&softname=romm&ssid=markom&sspassword=k72f8ks4L1GuYw&systemeid=10&jeuid=46335&media=support-2D(eu)",
	//        "screenshot_url": "https://neoclone.screenscraper.fr/api2/mediaJeu.php?devid=zurdi15&devpassword=xTJwoOFjOQG&softname=romm&ssid=markom&sspassword=k72f8ks4L1GuYw&systemeid=10&jeuid=46335&media=ss(wor)",
	//        "steamgrid_url": "https://neoclone.screenscraper.fr/api2/mediaJeu.php?devid=zurdi15&devpassword=xTJwoOFjOQG&softname=romm&ssid=markom&sspassword=k72f8ks4L1GuYw&systemeid=10&jeuid=46335&media=steamgrid",
	//        "title_screen_url": "https://neoclone.screenscraper.fr/api2/mediaJeu.php?devid=zurdi15&devpassword=xTJwoOFjOQG&softname=romm&ssid=markom&sspassword=k72f8ks4L1GuYw&systemeid=10&jeuid=46335&media=sstitle(wor)",
	//        "video_url": "https://neoclone.screenscraper.fr/api2/mediaVideoJeu.php?devid=zurdi15&devpassword=xTJwoOFjOQG&softname=romm&ssid=markom&sspassword=k72f8ks4L1GuYw&systemeid=10&jeuid=46335&media=video",
	//        "video_normalized_url": "https://neoclone.screenscraper.fr/api2/mediaVideoJeu.php?devid=zurdi15&devpassword=xTJwoOFjOQG&softname=romm&ssid=markom&sspassword=k72f8ks4L1GuYw&systemeid=10&jeuid=46335&media=video-normalized",
	//        "bezel_path": null,
	//        "box3d_path": null,
	//        "miximage_path": null,
	//        "physical_path": null,
	//        "marquee_path": null,
	//        "logo_path": null,
	//        "video_path": null,
	//        "ss_score": "",
	//        "first_release_date": 915148800,
	//        "alternative_names": [
	//          "10-pin Bowling",
	//          "10-Pin Bowling",
	//          "10-Pin Bowling"
	//        ],
	//        "companies": [
	//          "Majesco",
	//          "Morning Star Multimedia"
	//        ],
	//        "franchises": [],
	//        "game_modes": [],
	//        "genres": [
	//          "Sports / Bowling",
	//          "Sports"
	//        ]
	//      }
	VideoUrl      string `json:"video_url"`
	ScreenshotUrl string `json:"screenshot_url"`
}
type LaunchboxMetadata map[string]interface{}
type HasheousMetadata map[string]interface{}
type FlashpointMetadata map[string]interface{}
type HltbMetadata map[string]interface{}
type GamelistMetadata map[string]interface{}
type Sibling struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	FsNameNoTags   string `json:"fs_name_no_tags"`
	FsNameNoExt    string `json:"fs_name_no_ext"`
	SortComparator string `json:"sort_comparator"`
}
type MergedRaMetadata struct {
	FirstReleaseDate int64         `json:"first_release_date"`
	Genres           []string      `json:"genres"`
	Companies        []string      `json:"companies"`
	Achievements     []Achievement `json:"achievements"`
}
type Achievement struct {
	RaID               int    `json:"ra_id"`
	Title              string `json:"title"`
	Description        string `json:"description"`
	Points             int    `json:"points"`
	NumAwarded         int    `json:"num_awarded"`
	NumAwardedHardcore int    `json:"num_awarded_hardcore"`
	BadgeID            string `json:"badge_id"`
	BadgeURLLock       string `json:"badge_url_lock"`
	BadgePathLock      string `json:"badge_path_lock"`
	BadgeURL           string `json:"badge_url"`
	BadgePath          string `json:"badge_path"`
	DisplayOrder       int    `json:"display_order"`
	Type               string `json:"type"`
}
