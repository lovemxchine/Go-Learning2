package main

import (
	"fmt"
	"log"
	"time"
)

func createTask(task *CreateTask) error {
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	// EXEC() คือการรันคำสั่งเฉยๆไม่ต้องการ return ค่าไรมาเช็คแค่ error หรือ ไม่
	_, err := db.Exec(
		`INSERT INTO public.tasks(user_id, task_name, created_at , updated_at , status) VALUES ($1, $2, $3, $4, $5);`,
		task.User_Id, task.Name, currentTime, currentTime, false)

	return err
}

func getTask(id int) (Task, error) {
	var t Task
	row := db.QueryRow(
		`SELECT id ,user_id,task_name ,created_at ,updated_at ,status FROM tasks WHERE id = $1;`, id,
	)
	err := row.Scan(&t.Id, &t.User_Id, &t.Name, &t.CreateTime, &t.UpdateTime, &t.Status)

	if err != nil {
		return Task{}, err
	}

	return t, nil
}

func getTasks() ([]Task, error) {
	rows, err := db.Query("SELECT id , user_id, task_name , created_at , updated_at ,status FROM tasks ORDER BY id ASC;")
	if err != nil {
		log.Fatal(err)
	}
	var tasks []Task
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.Id, &t.User_Id, &t.Name, &t.CreateTime, &t.UpdateTime, &t.Status)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func updateTask(id int, task *UpdateTask) error {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	// EXEC() คือการรันคำสั่งเฉยๆไม่ต้องการ return ค่าไรมาเช็คแค่ error หรือ ไม่
	_, err := db.Exec(
		"UPDATE public.tasks SET  task_name=$1, updated_at=$2 , status=$3 WHERE id = $4;",
		task.Name,
		currentTime,
		task.Status,
		id,
	)
	return err
}

func deleteTask(id int) error {
	result, err := db.Exec(
		"DELETE FROM tasks WHERE id = $1",
		id,
	)

	if err != nil {
		return err
	}
	// RowsAffected คือนับว่ามี rows ไหนที่มีการเปลี่ยนแปลงบ้างถ้าไม่มีสักอันคือไม่มี id ที่ลบไป
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no task with id %d found", id)
	}

	return nil
}
