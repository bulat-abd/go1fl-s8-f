package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	// randSource источник псевдо случайных чисел.
	// Для повышения уникальности в качестве seed
	// используется текущее время в unix формате (в виде числа)
	randSource = rand.NewSource(time.Now().UnixNano())
	// randRange использует randSource для генерации случайных чисел
	randRange = rand.New(randSource)
)

// getTestParcel возвращает тестовую посылку
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	db, err := sql.Open("sqlite", "file:test.db?mode=memory&cache=shared")
	require.NoError(t, err)
	defer db.Close()
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS "parcel"
	(
    	number     integer
       	 constraint parcel_pk
            primary key autoincrement,
    	client     integer      not null,
    	status     VARCHAR(128) not null,
    	address    VARCHAR(512) not null,
    	created_at text         not null
	);`
	_, err = db.Exec(createTableSQL)
	require.NoError(t, err)
	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	p, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, p.Client, parcel.Client)
	require.Equal(t, p.Status, parcel.Status)
	require.Equal(t, p.Address, parcel.Address)
	require.Equal(t, p.CreatedAt, parcel.CreatedAt)

	err = store.Delete(id)
	require.NoError(t, err)
	_ ,err = store.Get(id)
	require.Error(t, err)
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	db, err := sql.Open("sqlite", "file:test.db?mode=memory&cache=shared")// настройте подключение к БД
	require.NoError(t, err)
	defer db.Close()
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS "parcel"
	(
    	number     integer
       	 constraint parcel_pk
            primary key autoincrement,
    	client     integer      not null,
    	status     VARCHAR(128) not null,
    	address    VARCHAR(512) not null,
    	created_at text         not null
	);`
	_, err = db.Exec(createTableSQL)
	require.NoError(t, err)
	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	newAddress := "new test address"
	err = store.SetAddress(id, newAddress)
    require.NoError(t, err)

	updatedParcel, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, updatedParcel.Address, newAddress)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	db, err := sql.Open("sqlite", "file:test.db?mode=memory&cache=shared")// настройте подключение к БД
	require.NoError(t, err)
	defer db.Close()
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS "parcel"
	(
    	number     integer
       	 constraint parcel_pk
            primary key autoincrement,
    	client     integer      not null,
    	status     VARCHAR(128) not null,
    	address    VARCHAR(512) not null,
    	created_at text         not null
	);`
	_, err = db.Exec(createTableSQL)
	require.NoError(t, err)
	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	err = store.SetStatus(id, ParcelStatusDelivered)
	require.NoError(t, err)

	updatedParcel, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, updatedParcel.Status, ParcelStatusDelivered)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	db, err := sql.Open("sqlite", "file:test.db?mode=memory&cache=shared")// настройте подключение к БД
	require.NoError(t, err)
	defer db.Close()
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS "parcel"
	(
    	number     integer
       	 constraint parcel_pk
            primary key autoincrement,
    	client     integer      not null,
    	status     VARCHAR(128) not null,
    	address    VARCHAR(512) not null,
    	created_at text         not null
	);`
	_, err = db.Exec(createTableSQL)
	require.NoError(t, err)
	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i])
		require.NoError(t, err)
		require.NotEmpty(t, id)
		parcels[i].Number = id
		parcelMap[id] = parcels[i]
	}

	storedParcels, err := store.GetByClient(client) // получите список посылок по идентификатору клиента, сохранённого в переменной client
	require.NoError(t, err)
	require.Len(t, storedParcels, len(parcels))
	for _, parcel := range storedParcels {
		require.Contains(t, parcelMap, parcel.Number)
		require.Equal(t, parcel, parcelMap[parcel.Number])
	}
}
