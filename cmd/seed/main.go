package main

import (
	"database/sql"
	"log/slog"
	"math/rand"
	"os"
	"time"

	"github.com/open-wm/blockehr/pkg/domain/entities"
	"github.com/open-wm/blockehr/pkg/infrastructure"
	repository "github.com/open-wm/blockehr/pkg/repositories"
)

func main() {
	logger := infrastructure.NewSLogger("blockehr")
	slog.SetDefault(logger)

	// this should probably go in a config file or options
	databaseName := "./db/dev.sqlite3"
	os.Remove(databaseName)
	db := infrastructure.NewDatabase(logger, databaseName)
	Seed(logger, db)
}

func Seed(logger *slog.Logger, db *sql.DB) {
	userRepo := repository.NewUserRepository(logger, db)
	// fileRepo := repository.NewFileRepository(logger, db)
	// recordRepo := repository.NewRecordRepository(logger, db)
	appointmentRepo := repository.NewAppointmentRepository(logger, db)
	profileRepo := repository.NewProfileRepository(logger, db)

	doctorIDs := []int{2}
	for i := 0; i < 10; i++ {
		doctor := entities.NewFakeProfile()
		doctor.Role = entities.DOCTOR
		profileRepo.Add(&doctor)
		user := entities.NewFakeUser()
		user.ProfileId = doctor.ID
		userRepo.SaveUser(&user)
		doctorIDs = append(doctorIDs, user.ID)
	}

	p, _ := profileRepo.GetAll()
	for _, v := range p {
		if v.Role == entities.PATIENT {
			for i := 0; i < 100; i++ {
				r := rand.Intn(12)
				r2 := rand.Intn(20)
				r3 := rand.Intn(5)
				statuses := []entities.AppointmentStatus{entities.PENDING, entities.PAID, entities.STARTED, entities.DONE, entities.CANCELED}
				t := time.Now()
				t1 := time.Date(t.Year(), t.Month(), t.Day(), 8+r, 0, 0, t.Nanosecond(), t.Location()).AddDate(0, 0, -7+r2)
				// add appointment for patients
				apppointment := entities.Appointment{
					PatientId:   v.ID,
					DoctorId:    doctorIDs[rand.Intn(len(doctorIDs))],
					Specialty:   "Cardiología",
					StartsAt:    t1,
					Status:      statuses[r3],
					Description: "Checkup",
				}
				appointmentRepo.Add(&apppointment)
			}
		}
	}

}
