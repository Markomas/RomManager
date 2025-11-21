package api

import (
	"RomManager/internal/api/romm"
	"RomManager/internal/config"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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
