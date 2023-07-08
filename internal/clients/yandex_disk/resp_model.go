package yandex_disk

// YandexResp https://yandex.com/dev/disk/api/reference/public.html
type YandexResp struct {
	Embedded struct {
		Sort      string `json:"sort"`
		PublicKey string `json:"public_key"`
		Items     []Item `json:"items"`
		Limit     int    `json:"limit"`
		Offset    int    `json:"offset"`
		Path      string `json:"path"`
		Total     int    `json:"total"`
	} `json:"_embedded"`
}

type Item struct {
	Path string `json:"path"`
	Type string `json:"type"`
	Name string `json:"name"`
	File string `json:"file,omitempty"`
	MD5  string `json:"md5"`
}
