package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sy_backend/db"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

type MenuModel struct {
	Id          string    `json:"_id"`
	Title       string    `json:"title"` // 名称
	Type        string    `json:"type"`  // 荤、素、汤
	Tag         []string  `json:"tag"`
	Ingredients []string  `json:"ingredients"` // 食材
	CookMethod  string    `json:"cook_method"` // 烹饪方法
	ImagePath   string    `json:"image_path"`  // 样例图
	Budget      float32   `json:"budget"`      // 预算
	ModifyAt    time.Time `json:"modify_at"`
}

func GetMenuList(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	filter := query.Get("type")

	page, _ := strconv.Atoi(query.Get("page"))
	limit, _ := strconv.Atoi(query.Get("limit"))

	b := strings.Builder{}
	b.WriteString("SELECT * FROM menu ")

	if filter != "" {
		b.WriteString(fmt.Sprintf("WHERE type='%s' ", filter))
	}
	b.WriteString("ORDER BY RANDOM() LIMIT ? OFFSET ?")

	rows, err := db.Conn.Query(b.String(), limit, limit*page)
	if err != nil {
		log.Println(err)
		return
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return
	}
	var menuList []MenuModel

	for rows.Next() {
		var menu MenuModel
		var tag, ingre, modify string

		err := rows.Scan(
			&menu.Id, &menu.Title, &menu.Type, &tag, &ingre, &menu.CookMethod,
			&menu.ImagePath, &menu.Budget, &modify,
		)
		if err != nil {
			log.Println(err)
			return
		}
		menu.Tag = strings.Split(tag, ",")
		menu.Ingredients = strings.Split(ingre, ",")

		menu.ModifyAt, _ = time.Parse("2006-01-02 15:04:05", modify)

		menuList = append(menuList, menu)
	}
	JSON(w, menuList)
}

func PostMenu(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	tp := r.FormValue("type")
	tag := r.FormValue("tag")
	ingredients := r.FormValue("ingredients") // 食材
	cookMethod := r.FormValue("cook_method")  // 烹饪方法
	budget := r.FormValue("budget")
	files := r.MultipartForm.File["image"]

	id, _ := gonanoid.New()

	path, err := NewResource(menu, id, files[0])
	if err != nil {
		log.Println(err)
		return
	}
	_, err = db.Conn.Exec(
		`INSERT INTO menu(
			_id, title, type, tag, ingredients, cook_method, image_path, budget)
		VALUES(?, ?, ?, ?, ?, ?, ?, ?)`,
		id,
		title,
		tp,
		tag,
		ingredients,
		cookMethod,
		path,
		budget,
	)
	if err != nil {
		log.Println(err)
		return
	}
	JSON(w, "OK")
}

func DeleteMenu(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	row := db.Conn.QueryRow("SELECT image_path FROM menu WHERE _id=?", id)

	var path string
	_ = row.Scan(&path)

	if err := DeleteResource(path); err != nil {
		log.Println(err)
		return
	}
	_, err := db.Conn.Exec("DELETE FROM menu WHERE _id=?", id)
	if err != nil {
		log.Println(err)
		return
	}
}

func PatchMenu(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	field := r.PathValue("field")

	var value string

	switch field {
	case "title":
		value = r.FormValue("title")
	case "ingredients":
		value = r.FormValue("ingredients")
	case "cook-method":
		value = r.FormValue("cook_method")
	case "budget":
		value = r.FormValue("budget")
	case "type":
		value = r.FormValue("type")
	case "tag":
		value = r.FormValue("tag")
	}
	if field == "cook-method" {
		field = "cook_method"
	}
	sql := fmt.Sprintf(
		"UPDATE menu SET %s='%s', update_at='%s' WHERE _id=?",
		field,
		value,
		time.Now().Format("2006-01-02 15:04:05"),
	)
	_, err := db.Conn.Exec(sql, id)
	if err != nil {
		log.Println(err)
		return
	}
}
