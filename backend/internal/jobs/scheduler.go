package jobs

import (
	"log"
	"time"

	"gorm.io/gorm"
)

type Scheduler struct {
	db   *gorm.DB
	done chan struct{}
}

func New(db *gorm.DB) *Scheduler {
	return &Scheduler{db: db, done: make(chan struct{})}
}

func (s *Scheduler) Start() {
	go s.runActivityReminders()
	go s.runMonthlyReports()
	log.Println("job scheduler started")
}

func (s *Scheduler) Stop() {
	close(s.done)
}

func (s *Scheduler) runActivityReminders() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			s.sendActivityReminders()
		case <-s.done:
			return
		}
	}
}

func (s *Scheduler) runMonthlyReports() {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if time.Now().Day() == 1 {
				s.generateMonthlyReport()
			}
		case <-s.done:
			return
		}
	}
}

func (s *Scheduler) sendActivityReminders() {
	// Stub: query upcoming activities within lookahead window and send notifications
	log.Println("activity reminders: checked")
}

func (s *Scheduler) generateMonthlyReport() {
	log.Println("monthly report: generated for", time.Now().Format("2006-01"))
}
