package app

type TrackCheckerApp struct {
}

func New() *TrackCheckerApp {
	return &TrackCheckerApp{}
}

func (a *TrackCheckerApp) Run() error {
	return nil
}
