package state

import (
	"fmt"
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

// for debugging purposes
func (s *importState) Set(mp map[string]any) {
	s.mu.Lock()
	defer s.mu.Unlock()

	getInt := func(key string) int {
		if val, ok := mp[key].(float64); ok {
			return int(val)
		}
		if val, ok := mp[key].(int); ok {
			return val
		}
		return 0
	}

	s.UploadingSrsCnt = getInt("UploadingSrsCnt")
	s.UploadingEpsCnt = getInt("UploadingEpsCnt")
	s.ProcessingSrsCnt = getInt("ProcessingSrsCnt")
	s.ProcessingEpsCnt = getInt("ProcessingEpsCnt")
	s.UploadingSrsCntTotal = getInt("UploadingSrsCntTotal")
	s.UploadingEpsCntTotal = getInt("UploadingEpsCntTotal")
	s.ProcessingSrsCntTotal = getInt("ProcessingSrsCntTotal")
	s.ProcessingEpsCntTotal = getInt("ProcessingEpsCntTotal")
	s.Stage = getInt("Stage")
	s.StageCnt = getInt("StageCnt")

	if b, ok := mp["uploadActive"].(bool); ok {
		s.uploadActive = b
	} else {
		s.uploadActive = false
	}
	errStr, ok := mp["uploadError"].(string)
	if ok {
		s.uploadError = fmt.Errorf("%v", errStr)
	}
}

func (s *importState) JsonMap() map[string]any {
	s.mu.Lock()
	defer s.mu.Unlock()

	var errStr *string
	if s.uploadError != nil {
		_s := s.uploadError.Error()
		errStr = &_s
	}
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
		"uploadError":           errStr,
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
