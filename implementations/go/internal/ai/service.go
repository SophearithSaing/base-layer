package ai

type AIModel string

const (
	Gemma_4   AIModel = "google/gemma-4-31B-it"
	KimiK_2_6 AIModel = "moonshotai/Kimi-K2.6"
	GLM_5_2   AIModel = "zai-org/GLM-5.2"
)

type Service struct {
	repo   *Repository
	apiKey string
}

func NewService(repo *Repository, apiKey string) *Service {
	return &Service{repo: repo, apiKey: apiKey}
}
