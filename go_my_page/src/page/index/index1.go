package index

import(
"net/http"
"common"  
"html/template"  
"log"
	"database/sql"
)

func Home(w http.ResponseWriter, r *http.Request) {
	var err error
	var db *sql.DB

	db, err = common.GetMysql()
	if  db == nil || err != nil {
		
	}
	defer common.CloseMysql(db)

	//tpl := template.New("some template") //创建一个模板
	tpl, err1 := template.ParseFiles("../tmpl/index1.html") //解析模板文件
	if err1 != nil {
		log.Printf("failed:" + err1.Error())
	}
	tpl.Execute(w,nil) //执行模板的merger操作
}