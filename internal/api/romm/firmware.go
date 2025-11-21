package romm

type Firmware struct {
	ID             int    `json:"id"`
	FileName       string `json:"file_name"`
	FileNameNoTags string `json:"file_name_no_tags"`
	FileNameNoExt  string `json:"file_name_no_ext"`
	FileExtension  string `json:"file_extension"`
	FilePath       string `json:"file_path"`
	FileSizeBytes  int64  `json:"file_size_bytes"`
	FullPath       string `json:"full_path"`
	IsVerified     bool   `json:"is_verified"`
	CRCHash        string `json:"crc_hash"`
	MD5Hash        string `json:"md5_hash"`
	SHA1Hash       string `json:"sha1_hash"`
	MissingFromFS  bool   `json:"missing_from_fs"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}
