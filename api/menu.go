package api

import (
	"net/http"
	"time"
)

type MenuModel struct {
	Id          uint      `json:"_id"`
	Title       string    `json:"title"`       // 名称
	Ingredients []string  `json:"ingredients"` // 食材
	CookMethod  string    `json:"cook_method"` // 烹饪方法
	ImageList   []string  `json:"image_list"`  // 样例图
	Budget      float32   `json:"budget"`      // 预算
	CreateAt    time.Time `json:"create_at"`
	UpdateAt    time.Time `json:"update_at"`
}

func GetMenuList(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GET /menu/list/0/10"))
}

func GetMenuInfo(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GET /menu/info"))
}
