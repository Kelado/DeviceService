package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/Kelado/DeviceService/models"
)

const (
	defaultDSN = "file:devices.db?cache=shared&mode=rwc"
	timeFormat = time.RFC3339
)

type SQLiteRepo struct {
	db        *sql.DB
	tableName string
}

type SQLiteRepoConfig struct {
	DSN string
}

func NewSQLiteDeviceRepo(config *SQLiteRepoConfig) *SQLiteRepo {
	dsn := defaultDSN
	if config != nil && config.DSN != "" {
		dsn = config.DSN
	}

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Unable to ping SQLiteRepo database ... %v", err)
	}

	repo := &SQLiteRepo{
		db:        db,
		tableName: "devices",
	}

	repo.Init()

	return repo
}

func (r *SQLiteRepo) Add(device *models.DeviceModel) error {
	statement, err := r.db.Prepare(fmt.Sprintf(`INSERT INTO %s (id, name, brand, created_at) VALUES (?, ?, ?, ?)`, r.tableName))
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(device.ID, device.Name, device.Brand, device.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *SQLiteRepo) GetById(id string) (*models.DeviceModel, error) {
	query := fmt.Sprintf(`SELECT id, name, brand, created_at FROM %s WHERE id = ?`, r.tableName)
	var device models.DeviceModel

	err := r.db.QueryRow(query, id).Scan(&device.ID, &device.Name, &device.Brand, &device.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("device with ID %s not found", id)
		}
		return nil, err
	}

	return &device, nil
}

func (r *SQLiteRepo) ListAll() ([]models.DeviceModel, error) {
	rows, err := r.db.Query(fmt.Sprintf("SELECT id, name, brand, created_at FROM %s", r.tableName))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []models.DeviceModel

	for rows.Next() {
		var device models.DeviceModel

		if err := rows.Scan(&device.ID, &device.Name, &device.Brand, &device.CreatedAt); err != nil {
			return nil, err
		}

		if err != nil {
			return nil, err
		}

		devices = append(devices, device)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return devices, nil
}

func (r *SQLiteRepo) Update(device *models.DeviceModel) error {
	if device.ID == "" {
		return fmt.Errorf("device ID must be set")
	}

	statement, err := r.db.Prepare(fmt.Sprintf(`UPDATE %s SET name = ?, brand = ? WHERE id = ?`, r.tableName))
	if err != nil {
		return err
	}
	defer statement.Close()

	result, err := statement.Exec(device.Name, device.Brand, device.ID)
	if err != nil {
		return err
	}

	// Check the number of rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		// No rows were affected, meaning no record with the given ID was found
		return fmt.Errorf("device with ID %s not found", device.ID)
	}

	return nil
}

func (r *SQLiteRepo) Delete(id string) error {
	statement, err := r.db.Prepare(fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, r.tableName))
	if err != nil {
		return err
	}
	defer statement.Close()

	result, err := statement.Exec(id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("device with ID %s not found", id)
	}

	return nil
}

func (r *SQLiteRepo) SearchByBrand(brand string) ([]models.DeviceModel, error) {
	query := fmt.Sprintf(`SELECT id, name, brand, created_at FROM %s WHERE brand = ?`, r.tableName)

	rows, err := r.db.Query(query, brand)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []models.DeviceModel

	for rows.Next() {
		var device models.DeviceModel

		err := rows.Scan(&device.ID, &device.Name, &device.Brand, &device.CreatedAt)
		if err != nil {
			return nil, err
		}

		devices = append(devices, device)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return devices, nil
}

func (r *SQLiteRepo) Init() {
	r.migrateTable()
}

func (r *SQLiteRepo) migrateTable() {
	statement, err := r.db.Prepare(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
        id TEXT PRIMARY KEY,
        name TEXT,
        brand TEXT,
        created_at DATETIME
    )`, r.tableName))

	if err != nil {
		log.Fatal("Error preparing SQL statement:", err)
	}

	_, err = statement.Exec()
	if err != nil {
		log.Fatal("Error executing SQL statement:", err)
	}

	log.Printf("Table %s created", r.tableName)
}
