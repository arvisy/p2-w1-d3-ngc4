package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"ngc4/config"
	"ngc4/entity"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func GetCrimeEvent(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	db, err := config.GetDB()
	if err != nil {
		log.Fatal("Failed connecting to Database")
	}
	defer db.Close()

	ctx := context.Background()
	var crimeEvent []entity.CrimeEvent

	query := `
		SELECT ID, HeroID, VillainID, Description, DateTime FROM crimeevent
	`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		ce := entity.CrimeEvent{}
		err := rows.Scan(&ce.ID, &ce.HeroID, &ce.VillainID, &ce.Description, &ce.DateTime)
		if err != nil {
			panic(err)
		}

		crimeEvent = append(crimeEvent, ce)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(crimeEvent)
}

func GetCrimeEventByID(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	db, err := config.GetDB()
	if err != nil {
		log.Fatal("Failed connecting to Database")
	}
	defer db.Close()

	ctx := context.Background()
	var crimeEvent entity.CrimeEvent

	id := p.ByName("id")

	query := `
		SELECT ID, HeroID, VillainID, Description, DateTime FROM crimeevent WHERE ID = ?
	`

	row := db.QueryRowContext(ctx, query, id)
	err = row.Scan(&crimeEvent.ID, &crimeEvent.HeroID, &crimeEvent.VillainID, &crimeEvent.Description, &crimeEvent.DateTime)

	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(crimeEvent)
}

func CreateCrimeEvent(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	db, err := config.GetDB()
	if err != nil {
		log.Fatal("Failed connecting to Database")
	}
	defer db.Close()

	ctx := context.Background()
	var crimeEvent entity.CrimeEvent

	err = json.NewDecoder(r.Body).Decode(&crimeEvent)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadGateway)
		return
	}

	query := `
		INSERT INTO crimeevent (HeroID, VillainID, Description, DateTime)
		VALUES (?, ?, ?, ?)
	`

	result, err := db.ExecContext(ctx, query, crimeEvent.HeroID, crimeEvent.VillainID, crimeEvent.Description, crimeEvent.DateTime)
	if err != nil {
		http.Error(w, "Failed to create Crime Event", http.StatusBadGateway)
		log.Println("Error", err)
		return
	}

	id, _ := result.LastInsertId()

	crimeEvent.ID = int(id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(crimeEvent)
}

func DeleteCrimeEventByID(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	db, err := config.GetDB()
	if err != nil {
		log.Fatal("Failed connecting to Database")
	}
	defer db.Close()

	ctx := context.Background()

	id := p.ByName("id")
	crimeEventID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid Crime Event ID", http.StatusBadGateway)
		return
	}

	existingCrimeEvent, err := GetCEByID(ctx, db, crimeEventID)
	if err != nil {
		http.Error(w, "Failed to retrieve existing Crime Event ID", http.StatusBadGateway)
		return
	}

	if existingCrimeEvent.ID == 0 {
		http.NotFound(w, r)
		return
	}

	err = DeleteCrime(ctx, db, crimeEventID)
	if err != nil {
		http.Error(w, "Failed to delete Crime Event", http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func GetCEByID(ctx context.Context, db *sql.DB, id int) (entity.CrimeEvent, error) {
	var crimeEvent entity.CrimeEvent

	query := `
		SELECT ID, HeroID, VillainID, Description, DateTime FROM crimeevent
		WHERE ID = ?
	`

	row := db.QueryRowContext(ctx, query, id)
	err := row.Scan(&crimeEvent.ID, &crimeEvent.HeroID, &crimeEvent.VillainID, &crimeEvent.Description, &crimeEvent.DateTime)
	return crimeEvent, err
}

func DeleteCrime(ctx context.Context, db *sql.DB, id int) error {
	query := `
        DELETE FROM crimeevent
        WHERE ID = ?
    `
	_, err := db.ExecContext(ctx, query, id)
	return err
}

func UpdateCrimeEventByID(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	db, err := config.GetDB()
	if err != nil {
		log.Fatal("Failed connecting to Database")
	}
	defer db.Close()

	ctx := context.Background()

	id := p.ByName("id")
	crimeEventID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid Crime Event ID", http.StatusBadRequest)
		return
	}

	existingCrimeEvent, err := GetCEByID(ctx, db, crimeEventID)
	if err != nil {
		http.Error(w, "Failed to retrieve Crime Event", http.StatusInternalServerError)
		log.Println("Error:", err)
		return
	}

	var updatedCrimeEvent entity.CrimeEvent
	err = json.NewDecoder(r.Body).Decode(&updatedCrimeEvent)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	existingCrimeEvent.HeroID = updatedCrimeEvent.HeroID
	existingCrimeEvent.VillainID = updatedCrimeEvent.VillainID
	existingCrimeEvent.Description = updatedCrimeEvent.Description
	existingCrimeEvent.DateTime = updatedCrimeEvent.DateTime

	err = updateCrimeE(ctx, db, existingCrimeEvent)
	if err != nil {
		http.Error(w, "Failed to Crime Event", http.StatusInternalServerError)
		log.Println("Error:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existingCrimeEvent)
}

func updateCrimeE(ctx context.Context, db *sql.DB, crimeEvent entity.CrimeEvent) error {
	query := `
        UPDATE crimeevent
        SET HeroID = ?, VillainID = ?, Description = ?, DateTime = ?
        WHERE ID = ?
    `
	_, err := db.ExecContext(ctx, query, crimeEvent.HeroID, crimeEvent.VillainID, crimeEvent.Description, crimeEvent.DateTime, crimeEvent.ID)
	return err
}
