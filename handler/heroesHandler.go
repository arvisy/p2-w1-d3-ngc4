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

func GetHeroes(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	db, err := config.GetDB()
	if err != nil {
		log.Fatal("Failed connecting to Database")
	}
	defer db.Close()

	ctx := context.Background()
	var hero []entity.Heroes

	query := `
		SELECT ID, Name, Universe, Skill, ImageURL FROM heroes
	`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		h := entity.Heroes{}
		err := rows.Scan(&h.ID, &h.Name, &h.Universe, &h.Skill, &h.ImageURL)
		if err != nil {
			panic(err)
		}

		hero = append(hero, h)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hero)
}

func GetHeroesByID(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	db, err := config.GetDB()
	if err != nil {
		log.Fatal("Failed connecting to Database")
	}
	defer db.Close()

	ctx := context.Background()
	var hero entity.Heroes

	id := p.ByName("id")

	query := `
		SELECT ID, Name, Universe, Skill, ImageURL FROM heroes WHERE ID = ?
	`

	row := db.QueryRowContext(ctx, query, id)
	err = row.Scan(&hero.ID, &hero.Name, &hero.Universe, &hero.Skill, &hero.ImageURL)

	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hero)
}

func CreateHero(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	db, err := config.GetDB()
	if err != nil {
		log.Fatal("Failed connecting to Database")
	}
	defer db.Close()

	ctx := context.Background()
	var hero entity.Heroes

	err = json.NewDecoder(r.Body).Decode(&hero)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadGateway)
		return
	}

	query := `
		INSERT INTO heroes (Name, Universe, Skill, ImageURL)
		VALUES (?, ?, ?, ?)
	`

	result, err := db.ExecContext(ctx, query, hero.Name, hero.Universe, hero.Skill, hero.ImageURL)
	if err != nil {
		http.Error(w, "Failed to create Hero", http.StatusBadGateway)
		log.Println("Error", err)
		return
	}

	id, _ := result.LastInsertId()

	hero.ID = int(id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(hero)
}

func DeleteHeroByID(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	db, err := config.GetDB()
	if err != nil {
		log.Fatal("Failed connecting to Database")
	}
	defer db.Close()

	ctx := context.Background()

	id := p.ByName("id")
	HeroID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid Hero ID", http.StatusBadGateway)
		return
	}

	existingHero, err := GetHByID(ctx, db, HeroID)
	if err != nil {
		http.Error(w, "Failed to retrieve existing Hero ID", http.StatusBadGateway)
		return
	}

	if existingHero.ID == 0 {
		http.NotFound(w, r)
		return
	}

	err = DeleteHero(ctx, db, HeroID)
	if err != nil {
		http.Error(w, "Failed to delete Hero", http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func GetHByID(ctx context.Context, db *sql.DB, id int) (entity.Heroes, error) {
	var hero entity.Heroes

	query := `
		SELECT ID, Name, Universe, Skill, ImageURL FROM heroes
		WHERE ID = ?
	`

	row := db.QueryRowContext(ctx, query, id)
	err := row.Scan(&hero.ID, &hero.Name, &hero.Universe, &hero.Skill, &hero.ImageURL)
	return hero, err
}

func DeleteHero(ctx context.Context, db *sql.DB, id int) error {
	query := `
        DELETE FROM heroes
        WHERE ID = ?
    `
	_, err := db.ExecContext(ctx, query, id)
	return err
}

func UpdateHeroByID(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	db, err := config.GetDB()
	if err != nil {
		log.Fatal("Failed connecting to Database")
	}
	defer db.Close()

	ctx := context.Background()

	id := p.ByName("id")
	HeroID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid Hero ID", http.StatusBadRequest)
		return
	}

	existingHero, err := GetHByID(ctx, db, HeroID)
	if err != nil {
		http.Error(w, "Failed to retrieve hero", http.StatusInternalServerError)
		log.Println("Error:", err)
		return
	}

	var updatedHero entity.Heroes
	err = json.NewDecoder(r.Body).Decode(&updatedHero)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	existingHero.Name = updatedHero.Name
	existingHero.Universe = updatedHero.Universe
	existingHero.Skill = updatedHero.Skill
	existingHero.ImageURL = updatedHero.ImageURL

	err = updateHero(ctx, db, existingHero)
	if err != nil {
		http.Error(w, "Failed to Crime Event", http.StatusInternalServerError)
		log.Println("Error:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existingHero)
}

func updateHero(ctx context.Context, db *sql.DB, hero entity.Heroes) error {
	query := `
        UPDATE heroes
        SET Name = ?, Universe = ?, Skill = ?, ImageURL = ?
        WHERE ID = ?
    `
	_, err := db.ExecContext(ctx, query, hero.Name, hero.Universe, hero.Skill, hero.ImageURL, hero.ID)
	return err
}
