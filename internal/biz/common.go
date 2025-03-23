package biz

var DEFAULT_PAGE_SIZE int32 = 10
var MINIMUM_PAGE_SIZE int32 = 1

type Pagination struct {
	Page int32 `json:"page"`
	Size int32 `json:"size"`
}

type Difficulty int

const (
	EASY Difficulty = iota
	MEDIUM
	HARD
	EXPERT
)

func (d Difficulty) String(difficulty int) string {
	switch difficulty {
	case 0:
		return "Easy"
	case 1:
		return "Medium"
	case 2:
		return "Hard"
	case 3:
		return "Expert"
	default:
		return "Unknown"
	}
}

func DifficultyFromString(difficulty string) Difficulty {
	switch difficulty {
	case "Easy":
		return EASY
	case "Medium":
		return MEDIUM
	case "Hard":
		return HARD
	case "Expert":
		return EXPERT
	default:
		return EASY
	}
}

type Audit struct {
	CreatedBy string `json:"created_by"`
	UpdatedBy string `json:"updated_by"`
	DeletedBy string `json:"deleted_by"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
}

func PaginationOrDefault(pagination *Pagination) *Pagination {
	if pagination == nil {
		return &Pagination{
			Page: 0,
			Size: 10,
		}
	}
	if pagination.Page < 0 {
		pagination.Page = 0
	}
	if pagination.Size < MINIMUM_PAGE_SIZE {
		pagination.Size = DEFAULT_PAGE_SIZE
	}
	return pagination
}
