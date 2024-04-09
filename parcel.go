package main

import (
	"database/sql"
	"fmt"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p

	// верните идентификатор последней добавленной записи
	query := "INSERT INTO parcel (client, status, address, created_at) VALUES (?, ?, ?, ?)"
	res, err := s.db.Exec(query, p.Client, p.Status, p.Address, p.CreatedAt)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil

}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка
	query := "SELECT * FROM parcel WHERE number = ?"

	res := s.db.QueryRow(query, number)

	// заполните объект Parcel данными из таблицы
	p := Parcel{}

	err := res.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return Parcel{}, err
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк
	r, err := s.db.Query("SELECT number, client, status, address, created_at FROM parcel WHERE client = :client", sql.Named("client", client))
	defer r.Close()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// заполните срез Parcel данными из таблицы
	var res []Parcel

	for r.Next() {
		p := Parcel{}
		err := r.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		res = append(res, p)
	}

	err = r.Err()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :number", sql.Named("status", status), sql.Named("number", number))

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered

	res, err := s.Get(number)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if res.Status != ParcelStatusRegistered {
		return nil
	}

	_, err = s.db.Exec("UPDATE parcel SET address = :address WHERE number = :number", sql.Named("address", address), sql.Named("number", number))
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	res, err := s.Get(number)

	if err != nil {
		fmt.Println(err)
		return err
	}

	if res.Status != ParcelStatusRegistered {
		return nil

	}

	_, err = s.db.Exec("DELETE FROM parcel WHERE number = :number", sql.Named("number", number))
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
