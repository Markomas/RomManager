package romm

import "time"

type SaveState struct {
	ID             int         `json:"id"`
	RomID          int         `json:"rom_id"`
	UserID         int         `json:"user_id"`
	FileName       string      `json:"file_name"`
	FileNameNoTags string      `json:"file_name_no_tags"`
	FileNameNoExt  string      `json:"file_name_no_ext"`
	FileExtension  string      `json:"file_extension"`
	FilePath       string      `json:"file_path"`
	FileSizeBytes  int64       `json:"file_size_bytes"`
	FullPath       string      `json:"full_path"`
	DownloadPath   string      `json:"download_path"`
	MissingFromFS  bool        `json:"missing_from_fs"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
	Emulator       string      `json:"emulator"`
	Screenshot     *Screenshot `json:"screenshot"`
}

type Screenshot struct {
	ID             int       `json:"id"`
	RomID          int       `json:"rom_id"`
	UserID         int       `json:"user_id"`
	FileName       string    `json:"file_name"`
	FileNameNoTags string    `json:"file_name_no_tags"`
	FileNameNoExt  string    `json:"file_name_no_ext"`
	FileExtension  string    `json:"file_extension"`
	FilePath       string    `json:"file_path"`
	FileSizeBytes  int64     `json:"file_size_bytes"`
	FullPath       string    `json:"full_path"`
	DownloadPath   string    `json:"download_path"`
	MissingFromFS  bool      `json:"missing_from_fs"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
