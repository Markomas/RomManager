package api

import (
	"RomManager/internal/api/romm"
	"RomManager/internal/config"
	"RomManager/internal/db/entity"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Romm struct {
	config *config.Config
}

func New(config *config.Config) *Romm {
	return &Romm{
		config: config,
	}
}

func (r *Romm) GetPlatforms() (romm.Platforms, error) {
	client := &http.Client{}
	platformsUrl := strings.TrimRight(r.config.Romm.Host, "/") + "/api/platforms"
	req, err := http.NewRequest("GET", platformsUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if r.config.Romm.Username != "" && r.config.Romm.Password != "" {
		req.SetBasicAuth(r.config.Romm.Username, r.config.Romm.Password)
	}

	resp, err := client.Do(req)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get platforms, status code: %s", resp.Status)
	}

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Println(string(body))

	var platforms romm.Platforms
	err = json.Unmarshal(body, &platforms)
	if err != nil {
		return nil, err
	}
	return platforms, nil
}

func (r *Romm) GetRomsByPlatform(id int, offset int, perPage int) (romm.Roms, int, error) {
	client := &http.Client{}
	fmt.Println(strings.TrimRight(r.config.Romm.Host, "/") + fmt.Sprintf("/api/roms?with_char_index=true&platform_id=%d&group_by_meta_id=false&order_by=name&order_dir=asc&limit=%d&offset=%d", id, perPage, offset))
	req, err := http.NewRequest("GET", strings.TrimRight(r.config.Romm.Host, "/")+fmt.Sprintf("/api/roms?with_char_index=true&platform_id=%d&group_by_meta_id=false&order_by=name&order_dir=asc&limit=50&offset=%d", id, offset), nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	if r.config.Romm.Username != "" && r.config.Romm.Password != "" {
		req.SetBasicAuth(r.config.Romm.Username, r.config.Romm.Password)
	}

	resp, err := client.Do(req)
	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("failed to get roms for platform, status code: %s", resp.Status)
	}
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}
	log.Println(string(body))
	var roms romm.RomsResponse
	err = json.Unmarshal(body, &roms)
	if err != nil {
		return nil, 0, err
	}
	return roms.Items, roms.Total, nil
}

func (r *Romm) GetRomByID(id int) (*romm.Rom, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", strings.TrimRight(r.config.Romm.Host, "/")+fmt.Sprintf("/api/roms/%d", id), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if r.config.Romm.Username != "" && r.config.Romm.Password != "" {
		req.SetBasicAuth(r.config.Romm.Username, r.config.Romm.Password)
	}

	resp, err := client.Do(req)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get roms for platform, status code: %s", resp.Status)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Println(string(body))
	var rom romm.Rom
	err = json.Unmarshal(body, &rom)
	if err != nil {
		return nil, err
	}
	return &rom, nil
}

func (r *Romm) GetFirmwaresByRomID(id int) ([]romm.Firmware, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", strings.TrimRight(r.config.Romm.Host, "/")+fmt.Sprintf("/api/roms/%d/firmwares", id), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if r.config.Romm.Username != "" && r.config.Romm.Password != "" {
		req.SetBasicAuth(r.config.Romm.Username, r.config.Romm.Password)
	}

	resp, err := client.Do(req)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get firmwares for rom, status code: %s", resp.Status)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Println(string(body))
	var firmwares []romm.Firmware
	err = json.Unmarshal(body, &firmwares)
	if err != nil {
		return nil, err
	}
	return firmwares, nil
}

func (r *Romm) DownloadRomm(rommID int, progress func(progress float64)) (*romm.Rom, *string, error) {
	rom, err := r.GetRomByID(rommID)
	if err != nil {
		return rom, nil, err
	}

	//TODO download supporting images

	rom.PlatformFsSlug = r.getPlatformFolder(rom.PlatformFsSlug)
	outputFilePath := r.getLocalRomPath(rom)
	err = os.MkdirAll(filepath.Dir(outputFilePath), 0755)
	if err != nil {
		return rom, nil, err
	}

	if _, err := os.Stat(outputFilePath); err == nil {
		if rom.Sha1Hash != "" {
			existingFileHash, err := calculateFileSha1(outputFilePath)
			if err != nil {
				return rom, &outputFilePath, fmt.Errorf("failed to calculate hash for existing file: %w", err)
			}
			if strings.EqualFold(existingFileHash, rom.Sha1Hash) {
				log.Printf("File %s already exists with correct SHA1 hash, skipping download.", outputFilePath)
				return rom, &outputFilePath, nil
			}
			log.Printf("File %s exists but SHA1 hash is different. Re-downloading.", outputFilePath)
		} else {
			log.Printf("File %s already exists, but no remote hash to compare. Skipping download.", outputFilePath)
			return rom, &outputFilePath, nil
		}
	} else if !os.IsNotExist(err) {
		return rom, &outputFilePath, err
	}

	downloadUrl := fmt.Sprintf("%s/api/roms/%d/content/output.zip", strings.TrimRight(r.config.Romm.Host, "/"), rommID)

	req, err := http.NewRequest("GET", downloadUrl, nil)
	if err != nil {
		return rom, &outputFilePath, err
	}
	if r.config.Romm.Username != "" && r.config.Romm.Password != "" {
		req.SetBasicAuth(r.config.Romm.Username, r.config.Romm.Password)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return rom, &outputFilePath, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return rom, &outputFilePath, fmt.Errorf("failed to download rom, status code: %s", resp.Status)
	}

	out, err := os.Create(outputFilePath)
	if err != nil {
		return rom, &outputFilePath, err
	}
	defer out.Close()

	counter := &writeCounter{
		total:      resp.ContentLength,
		progressFn: progress,
	}
	_, err = io.Copy(out, io.TeeReader(resp.Body, counter))
	if err != nil {
		return rom, &outputFilePath, err
	}

	return rom, &outputFilePath, nil
}

func (r *Romm) getLocalRomPath(rom *romm.Rom) string {
	platformFolder := r.getPlatformFolder(rom.PlatformFsSlug)
	return fmt.Sprintf("%s/%s/%s", r.config.System.RomsPath, platformFolder, rom.FsName)
}

func (r *Romm) getPlatformFolder(platformFsSlug string) string {
	platformFolder := platformFsSlug
	if mappedFolder, ok := r.config.PlatformFolderMapping[platformFsSlug]; ok {
		platformFolder = mappedFolder
	}
	return platformFolder
}

func calculateFileSha1(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha1.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func (r *Romm) GetSaveStates(rommId int, platformId int) ([]romm.SaveState, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", strings.TrimRight(r.config.Romm.Host, "/")+fmt.Sprintf("/api/states?rom_id=%d&platform_id=%d", rommId, platformId), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if r.config.Romm.Username != "" && r.config.Romm.Password != "" {
		req.SetBasicAuth(r.config.Romm.Username, r.config.Romm.Password)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var saveStates []romm.SaveState
	err = json.Unmarshal(body, &saveStates)
	if err != nil {
		return nil, err
	}

	return saveStates, nil
}

func (r *Romm) DownloadSaveStateToTmp(state romm.SaveState, rom entity.Rom) (*string, error) {
	destinationPath := filepath.Join(
		r.config.System.TmpPath,
		"savestates",
		fmt.Sprintf(
			"%s/%s.state%d",
			r.getPlatformFolder(rom.PlatformSlug),
			rom.FsNameNoExt,
			state.ID+100,
		),
	)

	if err := os.MkdirAll(filepath.Dir(destinationPath), 0o755); err != nil {
		return nil, fmt.Errorf("creating savestate directory: %w", err)
	}

	_ = strings.TrimRight(r.config.Romm.Host, "/") + state.DownloadPath

	return &destinationPath, nil
}
