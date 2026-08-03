package services

type _state struct {
	IUploadingSrsCnt      int
	IUploadingSrsCntTotal int

	IUploadingEpsCnt      int
	IUploadingEpsCntTotal int

	IUploadActive bool

	IUploadError error
}

var State = &_state{}
