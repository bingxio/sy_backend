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
	Title       string    `json:"title"`       // 名称
	Type        string    `json:"type"`        // 荤、素、汤
	Ingredients []string  `json:"ingredients"` // 食材
	CookMethod  string    `json:"cook_method"` // 烹饪方法
	ImageList   []string  `json:"image_list"`  // 样例图
	Budget      float32   `json:"budget"`      // 预算
	CreateAt    time.Time `json:"create_at"`
	UpdateAt    time.Time `json:"update_at"`
}

func GetRandomMenu(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	filter := query.Get("type")

	QL := "SELECT * FROM menu %s ORDER BY RANDOM() LIMIT 1"
	if filter != "" {
		QL = fmt.Sprintf(QL, "WHERE type=?")
	}
	row := db.Conn.QueryRow(QL, filter)
	if err := row.Err(); err != nil {
		log.Println(err)
	}
	var menu MenuModel
	var ingre, image, create, update string

	err := row.Scan(
		&menu.Id, &menu.Title, &menu.Type, &ingre, &menu.CookMethod,
		&image, &menu.Budget, &create, &update,
	)
	if err != nil {
		log.Println(err)
	}
	menu.Ingredients = strings.Split(ingre, ",")
	menu.ImageList = strings.Split(image, ",")

	menu.CreateAt, _ = time.Parse("2006-01-02 15:04:05", create)
	menu.UpdateAt, _ = time.Parse("2006-01-02 15:04:05", update)

	JSON(w, menu)
}

func GetMenuList(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.PathValue("page"))
	limit, _ := strconv.Atoi(r.PathValue("limit"))

	rows, err := db.Conn.Query(
		"SELECT * FROM menu LIMIT ? OFFSET ?", limit, limit*page)
	if err != nil {
		log.Println(err)
	}
	var menuList []MenuModel

	for rows.Next() {
		var menu MenuModel
		var ingre, image, create, update string

		err := rows.Scan(
			&menu.Id, &menu.Title, &menu.Type, &ingre, &menu.CookMethod,
			&image, &menu.Budget, &create, &update,
		)
		if err != nil {
			log.Println(err)
		}
		menu.Ingredients = strings.Split(ingre, ",")
		menu.ImageList = strings.Split(image, ",")

		menu.CreateAt, _ = time.Parse("2006-01-02 15:04:05", create)
		menu.UpdateAt, _ = time.Parse("2006-01-02 15:04:05", update)

		menuList = append(menuList, menu)
	}
	JSON(w, menuList)
}

func commaSlice(s []string) string {
	b := strings.Builder{}

	for i := 0; i < len(s); i++ {
		b.WriteString(s[i])

		if i != len(s)-1 {
			b.WriteByte(',')
		}
	}
	return b.String()
}

func PostMenu(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	tp := r.FormValue("type")
	ingredients := r.FormValue("ingredients") // 食材
	cookMethod := r.FormValue("cook_method")  // 烹饪方法
	budget := r.FormValue("budget")

	files := r.MultipartForm.File["image"]
	image := []string{}

	for _, v := range files {
		image = append(image, NewResource(menu, v))
	}
	imageList := strings.Builder{}

	for i := 0; i < len(image); i++ {
		imageList.WriteString(image[i])

		if i != len(image)-1 {
			imageList.WriteByte(',')
		}
	}
	nano, _ := gonanoid.New()

	_, err := db.Conn.Exec(
		`INSERT INTO menu(
			_id, title, type, ingredients, cook_method, image_list, budget)
		VALUES(?, ?, ?, ?, ?, ?, ?)`,
		nano,
		title,
		tp,
		ingredients,
		cookMethod,
		imageList.String(),
		budget,
	)
	if err != nil {
		log.Println(err)
	}
	JSON(w, "OK")
}

func DeleteMenu(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	row := db.Conn.QueryRow("SELECT image_list FROM menu WHERE _id=?", id)

	var imageList string
	_ = row.Scan(&imageList)

	for _, v := range strings.Split(imageList, ",") {
		DeleteResource(v)
	}
	_, err := db.Conn.Exec("DELETE FROM menu WHERE _id=?", id)
	if err != nil {
		log.Println(err)
	}
}

func pushMenuImage(r *http.Request) {
	id := r.PathValue("id")
	row := db.Conn.QueryRow("SELECT image_list FROM menu WHERE _id=?", id)
	var imageList string

	_ = row.Scan(&imageList)
	_, header, err := r.FormFile("image")

	if err != nil {
		log.Println(err)
	}
	path := NewResource(menu, header)
	var slice []string

	if len(imageList) != 0 {
		slice = strings.Split(imageList, ",")
	}
	slice = append(slice, path)

	_, err = db.Conn.Exec(
		"UPDATE menu SET image_list=? WHERE _id=?", commaSlice(slice), id)
	if err != nil {
		log.Println(err)
	}
}

func popMenuImage(r *http.Request) {
	id := r.PathValue("id")
	row := db.Conn.QueryRow("SELECT image_list FROM menu WHERE _id=?", id)

	var imageList string
	_ = row.Scan(&imageList)

	slice := strings.Split(imageList, ",")
	back := len(slice) - 1

	DeleteResource(slice[back])
	slice = slice[:back]

	_, err := db.Conn.Exec(
		"UPDATE menu SET image_list=? WHERE _id=?", commaSlice(slice), id)
	if err != nil {
		log.Println(err)
	}
}

// 对示例图片的操作
func MenuImage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		pushMenuImage(r) // 添加一张
	}
	if r.Method == "DELETE" {
		popMenuImage(r) // 删一张
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
	}
}
