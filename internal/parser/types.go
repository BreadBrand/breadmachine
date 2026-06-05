package parser

// RecipeDTO is the output of Parse(). Calculated fields (grams, bakerPercentage,
// densityGPerMl) are not included — the save endpoint computes them at persist time.
type RecipeDTO struct {
	Title            string          `json:"title"`
	Description      string          `json:"description"`
	Instructions     []string        `json:"instructions"`
	DoughIngredients []IngredientDTO `json:"doughIngredients"`
	OtherIngredients []IngredientDTO `json:"otherIngredients"`
	Servings         string          `json:"servings,omitempty"`
	PrepTime         int             `json:"prepTime,omitempty"`
	CookTime         int             `json:"cookTime,omitempty"`
	AdditionalTime   int             `json:"additionalTime,omitempty"`
	Confidence       ConfidenceMeta  `json:"confidence"`
}

// IngredientDTO represents one parsed ingredient line.
// ParseOK=false flags the line for user review in the modal; it is never silently dropped.
// Quantity is the raw token string as it appeared in the source (e.g. "1/2", "200", "2.75")
// so the frontend can display it as-is and convert to a number before saving.
type IngredientDTO struct {
	IngredientName string `json:"ingredientName"`
	Quantity       string `json:"quantity"`
	Unit           string `json:"unit"`
	Phase          string `json:"phase,omitempty"` // otherIngredients only
	RawLine        string `json:"rawLine"`
	ParseOK        bool   `json:"parseOK"`
}

// TitleDetectionMethod records how the title was found, for confidence scoring.
type TitleDetectionMethod int

const (
	TitleEmpty      TitleDetectionMethod = iota // no title found
	TitleFromHeader                             // explicit heading or line above Ingredients
	TitleHeuristic                              // first non-empty line fallback
)

// SectionMap is the output of DetectSections (Stage 2).
type SectionMap struct {
	Title                string
	TitleDetectionMethod TitleDetectionMethod
	Description          string
	MetadataLines        []string
	IngredientGroups     []IngredientGroup
	InstructionLines     []string
	NoLineBreaks         bool // true when input has fewer than 3 newlines
}

// IngredientGroup is one subsection of the ingredients block.
type IngredientGroup struct {
	Phase string // normalised header text, e.g. "starter build", "roasted pepper"
	Lines []string
}

// ConfidenceMeta scores 0.0–1.0 per field. Scores below 0.6 are flagged in the review modal.
type ConfidenceMeta struct {
	Title        float32 `json:"title"`
	Ingredients  float32 `json:"ingredients"`
	Instructions float32 `json:"instructions"`
}
