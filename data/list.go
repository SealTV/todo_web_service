package data

import (
	"database/sql"
	"log"

	"bitbucket.org/SealTV/go-site/model"
)

// GetAllLists return list of all List in table
func (db *postgresConnector) GetAllLists() (model.ListsCollection, error) {
	rows, err := db.Query("SELECT * FROM lists")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	return parseListsRows(rows)
}

// GetAllListsForUser return user lists
func (db *postgresConnector) GetAllListsForUser(user model.User) (model.ListsCollection, error) {
	return db.GetAllListsForUserId(user.Id)
}

// GetAllListsForUserId return user lists
func (db *postgresConnector) GetAllListsForUserId(user int) (model.ListsCollection, error) {
	rows, err := db.Query("SELECT * FROM lists WHERE user_id = $1", user)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	return parseListsRows(rows)
}

func (db *postgresConnector) GetListById(id int) (model.List, error) {
	var list model.List
	err := db.QueryRow(`SELECT * FROM lists WHERE id = $1`, id).Scan(&list.Id, &list.Name, &list.UserId)
	if err != nil {
		return list, err
	}

	return list, nil
}

func (db *postgresConnector) AddList(list model.List) (model.List, error) {
	err := db.QueryRow(`INSERT
			INTO lists(name, user_id)
			VALUES($1, $2)
			RETURNING id;`,
		list.Name, list.UserId).Scan(&list.Id)
	if err != nil {
		return list, err
	}

	return list, nil
}

func (db *postgresConnector) UpdateList(list model.List) (int64, error) {
	r, err := db.Exec(
		`UPDATE lists
				SET name = $2, user_id = $3
				WHERE id = $1`,
		list.Id, list.Name, list.UserId)
	if err != nil {
		log.Fatal(err)
	}
	return r.RowsAffected()
}

func (db *postgresConnector) DeleteList(list model.List) (int64, error) {
	return db.DeleteListById(list.Id)
}

func (db *postgresConnector) DeleteListById(list int) (int64, error) {
	tx, err := db.Begin()
	r, err := tx.Exec(`DELETE FROM todos WHERE list_id = $1`, list)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	r, err = tx.Exec(`DELETE FROM lists WHERE id = $1`, list)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	return r.RowsAffected()
}

func parseListsRows(rows *sql.Rows) (model.ListsCollection, error) {
	result := model.ListsCollection{}
	for rows.Next() {
		list := model.List{}

		err := rows.Scan(&list.Id, &list.Name, &list.UserId)
		if err != nil {
			return nil, err
		}

		result = append(result, list)
	}
	return result, nil
}
