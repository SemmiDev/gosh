package program

import (
	"context"
	"strings"
)

type SearchService struct {
	programDataStore *ProgramDataStore
}

func NewProgramSearchService(programDataStore *ProgramDataStore) *SearchService {
	return &SearchService{programDataStore: programDataStore}
}

func (s *SearchService) Search(ctx context.Context, keyword string) []Program {
	keyword = strings.TrimSpace(keyword)

	// change the space to underscore
	// because we use ts-vector
	keyword = strings.ReplaceAll(keyword, " ", "_")

	programs, err := s.programDataStore.SearchTerm(ctx, keyword)
	if err != nil {
		return []Program{}
	}

	return programs
}
