package state

import (
	"sync"
)

type importState struct {
	mu                   sync.Mutex
	UploadingSrsCnt      int
	UploadingSrsCntTotal int

	UploadingEpsCnt      int
	UploadingEpsCntTotal int

	ProcessingSrsCnt      int
	ProcessingSrsCntTotal int

	ProcessingEpsCnt      int
	ProcessingEpsCntTotal int

	Stage    int
	StageCnt int

	uploadActive bool

	uploadError error
}

var Import = &importState{}

func (s *importState) Json() map[string]any {
	s.mu.Lock()
	defer s.mu.Unlock()
	return map[string]any{
		"UploadingSrsCnt":       s.UploadingSrsCnt,
		"UploadingEpsCnt":       s.UploadingEpsCnt,
		"ProcessingSrsCnt":      s.ProcessingSrsCnt,
		"ProcessingEpsCnt":      s.ProcessingEpsCnt,
		"UploadingSrsCntTotal":  s.UploadingSrsCntTotal,
		"UploadingEpsCntTotal":  s.UploadingEpsCntTotal,
		"ProcessingSrsCntTotal": s.ProcessingSrsCntTotal,
		"ProcessingEpsCntTotal": s.ProcessingEpsCntTotal,
		"Stage":                 s.Stage,
		"StageCnt":              s.StageCnt,
		"uploadActive":          s.uploadActive,
		"uploadError":           s.uploadError,
	}
}

func (s *importState) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.UploadingSrsCnt = 0
	s.UploadingEpsCnt = 0

	s.UploadingSrsCntTotal = 0
	s.UploadingEpsCntTotal = 0

	s.ProcessingSrsCnt = 0
	s.ProcessingEpsCnt = 0

	s.ProcessingSrsCntTotal = 0
	s.ProcessingEpsCntTotal = 0

	s.uploadError = nil
}

func (s *importState) SetUploadError(err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.uploadError = err
}

func (s *importState) GetUploadError() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.uploadError
}

func (s *importState) SetUploadActive(active bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.uploadActive = active
}

func (s *importState) GetUploadActive() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.uploadActive
}
