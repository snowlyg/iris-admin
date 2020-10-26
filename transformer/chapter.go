package transformer

type Chapter struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Doc       Doc    `json:"doc"`
	CreatedAt string `json:"created_at"`
}
