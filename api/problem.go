package api

type Problem struct {
	ID       string            `json:"id"`
	TrackID  string            `json:"track_id"`
	Language string            `json:"language"`
	Slug     string            `json:"slug"`
	Name     string            `json:"name"`
	IsFresh  bool              `json:"fresh"`
	Files    map[string]string `json:"files"`
}

func (p *Problem) Save(dir string) error {
	return nil
}
