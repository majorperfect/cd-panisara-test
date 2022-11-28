package models

type findingType string

const (
	SAST findingType = "sast"
)

func (ft findingType) String() string {
	return string(ft)
}

type Metadata struct {
	Description string `json:"description"`
	Severity    string `json:"severity"`
}

type Location struct {
	Path     string   `json:"path"`
	Position Position `json:"position"`
}

type Position struct {
	Begin Begin `json:"begin"`
}
type Begin struct {
	Line uint32 `json:"line"`
}

type Finding struct {
	Type     findingType `json:"path"`
	RuleId   string      `json:"rule_id"`
	Location Location    `json:"location"`
	Metadata Metadata    `json:"metadata"`
}
