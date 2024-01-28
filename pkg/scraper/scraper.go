package scraper

// Scraper can scrape data from delivery service
type Scraper struct {
	Code  string `json:"code,omitempty"`
	Tasks []Task `json:"tasks,omitempty"`
}

func (s *Scraper) Scrape(args *Args) error {

	for i := range s.Tasks {
		err := s.Tasks[i].Process(args)
		if err != nil {
			return err
		}
	}

	return nil
}
