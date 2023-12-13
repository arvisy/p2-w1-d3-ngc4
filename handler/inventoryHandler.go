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

func GetInventory(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	db, err := config.GetDB()
	if err != nil {
		log.Fatal("Failed to connect")
	}
	defer db.Close()

	ctx := context.Background()
	var item []entity.Item

	query := `SELECT ID, Name, ItemCode, Stock, Description, Status FROM item`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		i := entity.Item{}
		err := rows.Scan(&i.ID, &i.Name, &i.ItemCode, &i.Stock, &i.Description, &i.Status)
		if err != nil {
			panic(err)
		}
		item = append(item, i)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func GetInventoryByID(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	db, err := config.GetDB()
	if err != nil {
		log.Fatal("Failed to connect")
	}
	defer db.Close()

	ctx := context.Background()
	var item entity.Item

	id := p.ByName("id")

	query := `SELECT ID, Name, ItemCode, Stock, Description, Status FROM item WHERE ID = ?`

	row := db.QueryRowContext(ctx, query, id)
	err = row.Scan(&item.ID, &item.Name, &item.ItemCode, &item.Stock, &item.Description, &item.Status)

	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func CreateInventory(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	db, err := config.GetDB()
	if err != nil {
		log.Fatal("Failed to connect")
	}
	defer db.Close()

	ctx := context.Background()

	var newItem entity.Item
	err = json.NewDecoder(r.Body).Decode(&newItem)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	query := `
        INSERT INTO item (ID, Name, ItemCode, Stock, Description, Status)
        VALUES (?, ?, ?, ?, ?, ?)
    `

	_, err = db.ExecContext(ctx, query, newItem.ID, newItem.Name, newItem.ItemCode, newItem.Stock, newItem.Description, newItem.Status)
	if err != nil {
		http.Error(w, "Failed to create inventory item", http.StatusBadRequest)
		log.Println("Error", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newItem)
}

func UpdateInventoryID(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	db, err := config.GetDB()
	if err != nil {
		log.Fatal("Failed to connect")
	}
	defer db.Close()

	ctx := context.Background()

	id := p.ByName("id")
	itemID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	existingItem, err := getItemByID(ctx, db, itemID)
	if err != nil {
		http.Error(w, "Failed to retrieve existing inventory item", http.StatusInternalServerError)
		log.Println("Error:", err)
		return
	}

	var updatedItem entity.Item
	err = json.NewDecoder(r.Body).Decode(&updatedItem)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	existingItem.Name = updatedItem.Name
	existingItem.ItemCode = updatedItem.ItemCode
	existingItem.Stock = updatedItem.Stock
	existingItem.Description = updatedItem.Description
	existingItem.Status = updatedItem.Status

	err = updateItem(ctx, db, existingItem)
	if err != nil {
		http.Error(w, "Failed to update inventory item", http.StatusInternalServerError)
		log.Println("Error:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existingItem)
}

func getItemByID(ctx context.Context, db *sql.DB, id int) (entity.Item, error) {
	var item entity.Item
	query := `
        SELECT ID, Name, ItemCode, Stock, Description, Status
        FROM item
        WHERE ID = ?
    `
	row := db.QueryRowContext(ctx, query, id)
	err := row.Scan(&item.ID, &item.Name, &item.ItemCode, &item.Stock, &item.Description, &item.Status)
	return item, err
}

func updateItem(ctx context.Context, db *sql.DB, item entity.Item) error {
	query := `
        UPDATE item
        SET Name = ?, ItemCode = ?, Stock = ?, Description = ?, Status = ?
        WHERE ID = ?
    `
	_, err := db.ExecContext(ctx, query, item.Name, item.ItemCode, item.Stock, item.Description, item.Status, item.ID)
	return err
}

func DeleteInventoryByID(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	db, err := config.GetDB()
	if err != nil {
		log.Fatal("Failed to connect")
	}
	defer db.Close()

	ctx := context.Background()

	id := p.ByName("id")
	itemID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	existingItem, err := getItemByID(ctx, db, itemID)
	if err != nil {
		http.Error(w, "Failed to retrieve existing inventory item", http.StatusBadRequest)
		log.Println("Error", err)
		return
	}

	if existingItem.ID == 0 {
		http.NotFound(w, r)
		return
	}

	err = deleteItem(ctx, db, itemID)
	if err != nil {
		http.Error(w, "Failed to delete inventory item", http.StatusBadRequest)
		log.Println("Error", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func deleteItem(ctx context.Context, db *sql.DB, id int) error {
	query := `
        DELETE FROM item
        WHERE ID = ?
    `
	_, err := db.ExecContext(ctx, query, id)
	return err
}
