package gc1

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
	//"github.com/hoisie/redis"
	. "github.com/hunterhug/go-image/go_image"
	// "crypto/md5"
	"encoding/json"

	"common"
	CONF "conf"
)

func GetAvatarInfo(w http.ResponseWriter, r *http.Request) {

	type measure struct {
		width, height int
		hasTransed    bool
		uploadPath    string
	}

	type urlRec struct {
		Width      int    `json:"width"`
		Height     int    `json:"height"`
		UploadPath string `json:"uploadPath"`
	}

	type InnerStruct struct {
		Width     int      `json:"width"`
		Height    int      `json:"height"`
		UploadUrl string   `json:"uploadUrl"`
		Urls      []urlRec `json:"urls"`
	}

	type InfoStructSuccess struct {
		Status int         `json:"status"`
		Data   InnerStruct `json:"data"`
	}
	type InfoStructFail struct {
		Status int    `json:"status"`
		Msg    string `json:"msg"`
	}
	var infofail InfoStructFail
	infofail.Status = 1
	var rebnn []byte
	r.ParseForm()
	userID := r.Form.Get("userID")
	sign := r.Form.Get("sign")
	resID := r.Form.Get("resID")
	types := r.Form.Get("types")

	if userID == "" {
		infofail.Msg = "userID is empty"
		rebnn, _ = json.Marshal(infofail)
		fmt.Fprintf(w, fmt.Sprintf("%s", rebnn))
		return
	}
	if sign == "" {
		infofail.Msg = "sign is empty"
		rebnn, _ = json.Marshal(infofail)
		fmt.Fprintf(w, fmt.Sprintf("%s", rebnn))
		return
	}
	if resID == "" {
		infofail.Msg = "resID is empty"
		rebnn, _ = json.Marshal(infofail)
		fmt.Fprintf(w, fmt.Sprintf("%s", rebnn))
		return
	}
	if types == "" {
		infofail.Msg = "types is empty"
		rebnn, _ = json.Marshal(infofail)
		fmt.Fprintf(w, fmt.Sprintf("%s", rebnn))
		return
	}

	log.Print("t1")
	//校验-----------
	sign_str := "resID=" + resID + "&types=" + types + "&userID=" + userID

	var check_url = CONF.User_sign_url

	u, _ := url.Parse(check_url)
	q := u.Query()
	q.Set("userID", userID)
	q.Set("checkstr", sign_str)
	q.Set("sign", sign)
	u.RawQuery = q.Encode()
	res, err := http.Get(u.String())
	if err != nil {
		infofail.Msg = "request fail"
		rebnn, _ = json.Marshal(infofail)
		fmt.Fprintf(w, fmt.Sprintf("%s", rebnn))
		return
	}
	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		infofail.Msg = "request fail !"
		rebnn, _ = json.Marshal(infofail)
		fmt.Fprintf(w, fmt.Sprintf("%s", rebnn))
		return
	}

	log.Print("t2")
	fmt.Printf("%s", result)
	var dat map[string]interface{}
	if err := json.Unmarshal([]byte(result), &dat); err == nil {
		fmt.Println(dat["status"])
	}

	if dat["status"] == nil || dat["status"] != float64(0) {
		infofail.Msg = "签名失败"
		rebnn, _ = json.Marshal(infofail)
		fmt.Fprintf(w, fmt.Sprintf("%s", rebnn))
		log.Printf("check failed, sign str:%s", sign_str)
		return
	}

	newFormat := time.Now().Local().Format("200601") + "avatar"
	uploadDir := CONF.User_avatar_file_path
	childDir := newFormat + "/"

	if runtime.GOOS == "linux" {
		childDir = newFormat + "/"
	}

	err = os.MkdirAll(uploadDir+childDir, 0666)
	if err != nil {
		fmt.Println("mkdir fail")
		infofail.Msg = "mkdir fail"
		rebnn, _ = json.Marshal(infofail)
		fmt.Fprintf(w, fmt.Sprintf("%s", rebnn))
		return
	}
	log.Print("t3")
	//--------------

	tpls := strings.Split(types, ",")
	if tpls == nil {
		infofail.Msg = "types param error"
		rebnn, _ = json.Marshal(infofail)
		fmt.Fprintf(w, fmt.Sprintf("%s", rebnn))
		log.Printf("types param error")
		return
	}

	_mArr := make([]urlRec, 0, 10)
	_mdict := make(map[string]measure)
	for _, tp := range tpls {
		_vs := strings.Split(tp, "*")
		if len(_vs) != 2 {
			infofail.Msg = "types is error:" + tp
			rebnn, _ = json.Marshal(infofail)
			fmt.Fprintf(w, fmt.Sprintf("%s", rebnn))
			log.Printf("types param error" + tp)
			return
		}

		__width, err := strconv.Atoi(_vs[0])
		__height, err2 := strconv.Atoi(_vs[1])

		if err != nil || err2 != nil {
			infofail.Msg = "types is error:" + tp
			rebnn, _ = json.Marshal(infofail)
			fmt.Fprintf(w, fmt.Sprintf("%s", rebnn))
			log.Printf("types param error" + tp)
			return
		}

		_mdict[tp] = measure{__width, __height, false, ""}
	}

	log.Print("t4")
	db, err := common.GetMysql()
	if db == nil || err != nil {
		infofail.Msg = err.Error()
		rebnn, _ = json.Marshal(infofail)
		fmt.Fprintf(w, fmt.Sprintf("%s", rebnn))
		return
	}
	defer db.Close()

	rows, err := db.Query("select `resID`, `userID`, `filePath`, `fileExt`, `width`, `height` from `ko_resource`.`avatar` where `resID` = ?", resID)
	if err != nil {
		infofail.Msg = err.Error() + "db.Query fail"
		rebnn, _ = json.Marshal(infofail)
		fmt.Fprintf(w, fmt.Sprintf("%s", rebnn))
		log.Printf("select avatar resource failed")
		return
	}
	defer rows.Close()

	var _resID, _dx, _dy int
	var _fileExt string
	var _userID, _uploadFilePath string

	_resID = 0

	for rows.Next() {
		if err := rows.Scan(&_resID, &_userID, &_uploadFilePath, &_fileExt, &_dx, &_dy); err != nil {
			infofail.Msg = err.Error() + "db.Scan fail"
			rebnn, _ = json.Marshal(infofail)
			fmt.Fprintf(w, fmt.Sprintf("%s", rebnn))
			log.Printf("select from avatar resource failed")
			return
		}

		break
	}

	if _resID == 0 {
		infofail.Msg = "no this resouce"
		rebnn, _ = json.Marshal(infofail)
		fmt.Fprintf(w, fmt.Sprintf("%s", rebnn))
		log.Printf("no this resouce:", resID)
		return
	}

	if !strings.EqualFold(userID, _userID) {
		infofail.Msg = "resource not belong to this user"
		rebnn, _ = json.Marshal(infofail)
		fmt.Fprintf(w, fmt.Sprintf("%s", rebnn))
		log.Printf("resource :", _resID, " not belong to user:", userID)
		return
	}

	log.Print("t5")
	//查询哪些已经转过码了
	rows1, err := db.Query("select `filePath`, `width`, `height` from `ko_resource`.`avatar_file` where resID=?", resID)
	if err != nil {
		infofail.Msg = err.Error() + "db.Query fail"
		rebnn, _ = json.Marshal(infofail)
		fmt.Fprintf(w, fmt.Sprintf("%s", rebnn))
		log.Printf("select avatar file failed")
		return
	}
	defer rows1.Close()

	var _filePah string
	var __width, __height string

	for rows1.Next() {
		if err := rows1.Scan(&_filePah, &__width, &__height); err != nil {
			infofail.Msg = err.Error() + "db.Scan fail"
			rebnn, _ = json.Marshal(infofail)
			fmt.Fprintf(w, fmt.Sprintf("%s", rebnn))
			log.Printf("select from avatar file failed")
			return
		}

		__tp := __width + "*" + __height
		__m, ok := _mdict[__tp]
		if ok {
			__m.hasTransed = true
			__m.uploadPath = _filePah
			_mdict[__tp] = __m
		}

	}

	timestamp := fmt.Sprintf("%d", time.Now().UnixNano())

	stmt, err := db.Prepare("INSERT INTO `ko_resource`.`avatar_file`(`resID`,`filePath`,`inputTime`,`width`,`height`) VALUES(?, ?, now(), ?, ?)")
	if err != nil {
		infofail.Msg = err.Error() + "db.Prepare fail"
		rebnn, _ = json.Marshal(infofail)
		fmt.Fprintf(w, fmt.Sprintf("%s", rebnn))
		return
	}
	defer stmt.Close()

	log.Print("t6")
	for _, _m := range _mdict {
		if _m.hasTransed {
			_mArr = append(_mArr, urlRec{_m.width, _m.height, _m.uploadPath})
			continue
		}

		_width := _m.width
		_height := _m.height

		save_path := fmt.Sprintf("%s%s%s%dx%d%s", childDir, timestamp, userID, _width, _height, _fileExt)
		err = ThumbnailF2F(uploadDir+_uploadFilePath, uploadDir+save_path, _width, _height)
		if err != nil {
			infofail.Msg = err.Error() + "Thumbnail img fail"
			rebnn, _ = json.Marshal(infofail)
			fmt.Fprintf(w, fmt.Sprintf("%s", rebnn))
			return
		}
		_, err = stmt.Exec(resID, save_path, _width, _height)
		if err != nil {
			infofail.Msg = err.Error()
			rebnn, _ = json.Marshal(infofail)
			fmt.Fprintf(w, fmt.Sprintf("%s", rebnn))
		}

		_m.hasTransed = true
		_m.uploadPath = save_path

		_mArr = append(_mArr, urlRec{_m.width, _m.height, _m.uploadPath})
	}

	log.Print("t7")
	innerData := InnerStruct{
		Width:     _dx,
		Height:    _dy,
		UploadUrl: _uploadFilePath,
		Urls:      _mArr,
	}

	info := InfoStructSuccess{
		Status: 0,
		Data:   innerData,
	}
	bnn, err := json.Marshal(info)
	fmt.Fprintf(w, fmt.Sprintf("%s", bnn))
}
