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

func GetVillain(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	db, err := config.GetDB()
	if err != nil {
		log.Fatal("Failed connecting to Database")
	}
	defer db.Close()

	ctx := context.Background()
	var villain []entity.Villain

	query := `
		SELECT ID, Name, Universe, ImageURL FROM villain
	`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		v := entity.Villain{}
		err := rows.Scan(&v.ID, &v.Name, &v.Universe, &v.ImageURL)
		if err != nil {
			panic(err)
		}

		villain = append(villain, v)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(villain)
}

func GetVillainByID(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	db, err := config.GetDB()
	if err != nil {
		log.Fatal("Failed connecting to Database")
	}
	defer db.Close()

	ctx := context.Background()
	var villain entity.Villain

	id := p.ByName("id")

	query := `
		SELECT ID, Name, Universe, ImageURL FROM villain WHERE ID = ?
	`

	row := db.QueryRowContext(ctx, query, id)
	err = row.Scan(&villain.ID, &villain.Name, &villain.Universe, &villain.ImageURL)

	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(villain)
}

func CreateVillain(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	db, err := config.GetDB()
	if err != nil {
		log.Fatal("Failed connecting to Database")
	}
	defer db.Close()

	ctx := context.Background()
	var villain entity.Villain

	err = json.NewDecoder(r.Body).Decode(&villain)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadGateway)
		return
	}

	query := `
		INSERT INTO villain (Name, Universe, ImageURL)
		VALUES (?, ?, ?)
	`

	result, err := db.ExecContext(ctx, query, villain.Name, villain.Universe, villain.ImageURL)
	if err != nil {
		http.Error(w, "Failed to create Villain", http.StatusBadGateway)
		log.Println("Error", err)
		return
	}

	id, _ := result.LastInsertId()

	villain.ID = int(id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(villain)
}

func DeleteVillainByID(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	db, err := config.GetDB()
	if err != nil {
		log.Fatal("Failed connecting to Database")
	}
	defer db.Close()

	ctx := context.Background()

	id := p.ByName("id")
	villainID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid Villain ID", http.StatusBadGateway)
		return
	}

	existingVillain, err := GetVByID(ctx, db, villainID)
	if err != nil {
		http.Error(w, "Failed to retrieve existing Villain ID", http.StatusBadGateway)
		return
	}

	if existingVillain.ID == 0 {
		http.NotFound(w, r)
		return
	}

	err = DeleteVillain(ctx, db, villainID)
	if err != nil {
		http.Error(w, "Failed to delete Villain", http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func GetVByID(ctx context.Context, db *sql.DB, id int) (entity.Villain, error) {
	var villain entity.Villain

	query := `
		SELECT ID, Name, Universe, ImageURL FROM villain
		WHERE ID = ?
	`

	row := db.QueryRowContext(ctx, query, id)
	err := row.Scan(&villain.ID, &villain.Name, &villain.Universe, &villain.ImageURL)
	return villain, err
}

func DeleteVillain(ctx context.Context, db *sql.DB, id int) error {
	query := `
        DELETE FROM villain
        WHERE ID = ?
    `
	_, err := db.ExecContext(ctx, query, id)
	return err
}

func UpdateVillainByID(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	db, err := config.GetDB()
	if err != nil {
		log.Fatal("Failed connecting to Database")
	}
	defer db.Close()

	ctx := context.Background()

	id := p.ByName("id")
	villainID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid Villain ID", http.StatusBadRequest)
		return
	}

	existingVillain, err := GetVByID(ctx, db, villainID)
	if err != nil {
		http.Error(w, "Failed to retrieve Villain", http.StatusInternalServerError)
		log.Println("Error:", err)
		return
	}

	var updateVillain entity.Villain
	err = json.NewDecoder(r.Body).Decode(&updateVillain)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	existingVillain.Name = updateVillain.Name
	existingVillain.Universe = updateVillain.Universe
	existingVillain.ImageURL = updateVillain.ImageURL

	err = updateVillainDB(ctx, db, existingVillain)
	if err != nil {
		http.Error(w, "Failed to Villain", http.StatusInternalServerError)
		log.Println("Error:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existingVillain)
}

func updateVillainDB(ctx context.Context, db *sql.DB, villain entity.Villain) error {
	query := `
        UPDATE villain
        SET Name = ?, Universe = ?, ImageURL = ?
        WHERE ID = ?
    `
	_, err := db.ExecContext(ctx, query, villain.Name, villain.Universe, villain.ImageURL, villain.ID)
	return err
}
