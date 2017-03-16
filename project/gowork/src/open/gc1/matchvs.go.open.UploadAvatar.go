package gc1

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"syscall"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	//"github.com/hoisie/redis"
	//"crypto/md5"
	. "github.com/hunterhug/go-image/go_image"

	"encoding/json"

	"common"
	CONF "conf"
)

func UploadAvatar(w http.ResponseWriter, r *http.Request) {
	type InfoStructSuccessInner struct {
		ResID int64 `json:"resID"`
	}
	type InfoStructSuccess struct {
		Status int                    `json:"status"`
		Data   InfoStructSuccessInner `json:"data"`
	}
	type InfoStructFail struct {
		Status int    `json:"status"`
		Code   int    `json:"code"`
		Msg    string `json:"msg"`
	}
	var infofail InfoStructFail
	infofail.Status = 1
	infofail.Code = 200
	var rebnn []byte
	//var jerr error

	uploadDir := CONF.User_avatar_file_path
	fileMaxSize := 1024 * 100 //字节
	var full_file_name string

	newFormat := time.Now().Local().Format("200601") + "avatar"
	childDir := newFormat + "//"
	uploadPath := uploadDir + childDir

	if runtime.GOOS == "linux" {
		childDir = newFormat + "/"
	}
	syscall.Umask(0000);
	err := os.MkdirAll(uploadPath, 0777)
	if err != nil {
		fmt.Println("mkdir fail")
		infofail.Msg = "mkdir fail"
		rebnn, _ = json.Marshal(infofail)
		fmt.Fprintf(w, fmt.Sprintf("%s", rebnn))
		return
	}

	log.Printf("os:%s", runtime.GOOS)

	r.ParseForm()
	userID := r.Form.Get("userID")
	sign := r.Form.Get("sign")
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

	log.Printf("t1")

	//校验-----------
	sign_str := "userID=" + userID
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
	log.Printf("t1.1")
	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		infofail.Msg = "request fail !"
		rebnn, _ = json.Marshal(infofail)
		fmt.Fprintf(w, fmt.Sprintf("%s", rebnn))
		return
	}

	log.Printf("t1.2")
	fmt.Printf("%s", result)
	var dat map[string]interface{}
	if err := json.Unmarshal([]byte(result), &dat); err == nil {
		fmt.Println(dat["status"])
	}

	log.Printf("t1.3")
	//check ifnot post
	if !strings.EqualFold(r.Method, "POST") {
		infofail.Msg = "upload"
		rebnn, _ = json.Marshal(infofail)
		fmt.Fprintf(w, fmt.Sprintf("%s", rebnn))
		log.Printf("upload method error:" + r.Method)
		return
	}
	log.Printf("t1.4")

	//parse the multipart form in the request
	err = r.ParseMultipartForm(100000)
	if err != nil {
		infofail.Msg = err.Error() + strconv.Itoa(http.StatusInternalServerError)
		rebnn, _ = json.Marshal(infofail)
		fmt.Fprintf(w, fmt.Sprintf("%s", rebnn))
		return
	}
	m := r.MultipartForm
	if m == nil {
		infofail.Msg = "multipart form is nil "
		rebnn, _ = json.Marshal(infofail)
		fmt.Fprintf(w, fmt.Sprintf("%s", rebnn))
		return
	}
	log.Printf("t1.5")
	//get the *fileheaders
	files := m.File["myfiles"]
	if files == nil {
		infofail.Msg = "not myfiles"
		rebnn, _ = json.Marshal(infofail)
		fmt.Fprintf(w, fmt.Sprintf("%s", rebnn))
		return
	}

	//for each fileheader, get a handle to the actual file
	file, err := files[0].Open()
	if err != nil {
		infofail.Msg = err.Error() + strconv.Itoa(http.StatusInternalServerError)
		rebnn, _ = json.Marshal(infofail)
		fmt.Fprintf(w, fmt.Sprintf("%s", rebnn))
		log.Printf("upload file open failed")
		return
	}
	defer file.Close()

	if dat["status"] == nil || dat["status"] != float64(0) {
		log.Printf("签名失败")
		infofail.Code = 15
		infofail.Msg = "签名失败"
		rebnn, _ = json.Marshal(infofail)
		fmt.Fprintf(w, fmt.Sprintf("%s", rebnn))
		return
	}

	filenums := 0
	for _, _ = range files {
		filenums = filenums + 1
	}
	if filenums != 1 {
		infofail.Msg = "files upload"
		rebnn, _ = json.Marshal(infofail)
		fmt.Fprintf(w, fmt.Sprintf("%s", rebnn))
		log.Printf("file upload num error:%d", filenums)
		return
	}
	log.Printf("t1.6")

	// 获取文件大小的接口
	type Size interface {
		Size() int64
	}
	if sizeInterface, ok := file.(Size); !ok {
		infofail.Msg = "get file size fail "
		rebnn, _ = json.Marshal(infofail)
		fmt.Fprintf(w, fmt.Sprintf("%s", rebnn))
		return
	} else {
		if sizeInterface.Size() > int64(fileMaxSize) {
			infofail.Msg = " file size > 100kb"
			rebnn, _ = json.Marshal(infofail)
			fmt.Fprintf(w, fmt.Sprintf("%s", rebnn))
			return
		}
	}
	log.Printf("t1.7")
	//------------------------------------------------------------- 生成目标文件 begin
	//获取文件类型
	file_ext := path.Ext(files[0].Filename)
	if !strings.Contains(".jpg.png.jpeg", file_ext) {
		infofail.Msg = " file type is not jpg png jpeg"
		rebnn, _ = json.Marshal(infofail)
		fmt.Fprintf(w, fmt.Sprintf("%s", rebnn))
		log.Printf(" file type is not jpg png jpeg")
		return
	}
	timestamp := fmt.Sprintf("%d", time.Now().UnixNano())
	file_name := timestamp + userID + file_ext
	full_file_name = childDir + file_name
	dst, err := os.Create(uploadDir + full_file_name)
	if err != nil {
		infofail.Msg = err.Error() + strconv.Itoa(http.StatusInternalServerError)
		rebnn, _ = json.Marshal(infofail)
		fmt.Fprintf(w, fmt.Sprintf("%s", rebnn))
		log.Printf("t1.8, create file error:" + uploadDir + full_file_name)
		return
	}

	log.Printf("t2, file ext:" + file_ext + ", file_name:" + file_name + ", full name:" + full_file_name)
	//copy the uploaded file to the destination file
	if _, err := io.Copy(dst, file); err != nil {
		dst.Close()
		infofail.Msg = err.Error() + strconv.Itoa(http.StatusInternalServerError)
		rebnn, _ = json.Marshal(infofail)
		fmt.Fprintf(w, fmt.Sprintf("%s", rebnn))
		log.Printf("copy failed, dst:%s", uploadDir+full_file_name)
		return
	}
	//------------------------------------------------------------- 生成目标文件 end
	//------------------------------------------------------------------获取文件图片尺寸  begin
	dst.Close()
	dst, _ = os.Open(uploadDir + full_file_name)
	defer dst.Close()
	img, _, err := image.Decode(dst)
	if err != nil {
		infofail.Msg = err.Error() + " Decode file fail"
		rebnn, _ = json.Marshal(infofail)
		fmt.Fprintf(w, fmt.Sprintf("%s", rebnn))
		log.Printf("decode failed, dir: %s, err:%s", uploadDir+full_file_name, err)
		return
	}
	bound := img.Bounds()
	dx := bound.Dx()
	dy := bound.Dy()
	//------------------------------------------------------------------获取文件图片尺寸  end
	//------------------------------------------------------------------获取数据库连接    begin
	db, err := common.GetMysql()
	if db == nil || err != nil {
		infofail.Msg = err.Error()
		rebnn, _ = json.Marshal(infofail)
		fmt.Fprintf(w, fmt.Sprintf("%s", rebnn))
		return
	}
	defer db.Close()
	//------------------------------------------------------------------获取数据库连接	  end
	log.Print("t3")
	//----------------------------------------------------------- 选择数据库尺寸 begin
	type measure struct {
		width, height int
	}
	_meas := make([]measure, 0, 10)
	rows1, err := db.Query("select width,height from `ko_resource`.`avatar_measure` where status=1")
	if err != nil {
		infofail.Msg = err.Error() + "db.Query fail"
		rebnn, _ = json.Marshal(infofail)
		fmt.Fprintf(w, fmt.Sprintf("%s", rebnn))
		return
	}
	defer rows1.Close()
	var _width, _height int
	var _rows1Num = 0
	for rows1.Next() {
		if err := rows1.Scan(&_width, &_height); err != nil {
			infofail.Msg = err.Error() + "db.Scan fail"
			rebnn, _ = json.Marshal(infofail)
			fmt.Fprintf(w, fmt.Sprintf("%s", rebnn))
			return
		}
		_meas = append(_meas, measure{_width, _height})
		_rows1Num = _rows1Num + 1
	}
	//-----------------------------------------------------------   选择数据库尺寸 end
	//------------------------------------------------------------------填充默认地址
	_mysql_res, err := db.Exec("INSERT INTO `ko_resource`.`avatar`(`filePath`,`fileExt`,`userID`,`width`, `height`,`inputTime`) VALUES(?, ?, ?, ?,?,now())", full_file_name, file_ext, userID, dx, dy)
	if err != nil {
		infofail.Msg = err.Error() + "db.insert avatar fail"
		rebnn, _ = json.Marshal(infofail)
		fmt.Fprintf(w, fmt.Sprintf("%s", rebnn))
		log.Print("db.insert avatar failed,", err)
		return
	}
	_res_id, err := _mysql_res.LastInsertId()
	if err != nil {
		infofail.Msg = err.Error() + "db.insert avatar fail"
		rebnn, _ = json.Marshal(infofail)
		fmt.Fprintf(w, fmt.Sprintf("%s", rebnn))
		log.Print("get insert avatar id failed,", err)
		return
	}

	log.Print("t4")
	//------------------------------------------------------------------填充默认地址
	//--------------------------------------------------------------------------------准备插入转换地址
	stmt, err := db.Prepare("INSERT INTO `ko_resource`.`avatar_file`(`resID`,`filePath`,`inputTime`,`width`,`height`) VALUES(?, ?, now(), ?, ?)")
	if err != nil {
		infofail.Msg = err.Error() + "db.Prepare fail"
		rebnn, _ = json.Marshal(infofail)
		fmt.Fprintf(w, fmt.Sprintf("%s", rebnn))
		return
	}
	defer stmt.Close()
	//--------------------------------------------------------------------------------准备插入转换地址
	//-----------------------------------------------------------begin 插入默认尺寸
	_, err = stmt.Exec(_res_id, full_file_name, dx, dy)
	if err != nil {
		infofail.Msg = err.Error()
		rebnn, _ = json.Marshal(infofail)
		fmt.Fprintf(w, fmt.Sprintf("%s", rebnn))
	}
	//-----------------------------------------------------------end   插入默认尺寸
	for _, _mea := range _meas {
		_width = _mea.width
		_height = _mea.height
		if dx == _width && dy == _height {
			continue
		} else {
			save_path := fmt.Sprintf("%s%s%s%dx%d%s", childDir, timestamp, userID, _width, _height, file_ext)
			err = ThumbnailF2F(uploadDir+full_file_name, uploadDir+save_path, _width, _height)
			if err != nil {
				infofail.Msg = err.Error() + "Thumbnail img fail"
				rebnn, _ = json.Marshal(infofail)
				fmt.Fprintf(w, fmt.Sprintf("%s", rebnn))
				return
			}
			_, err = stmt.Exec(_res_id, save_path, _width, _height)
			if err != nil {
				infofail.Msg = err.Error()
				rebnn, _ = json.Marshal(infofail)
				fmt.Fprintf(w, fmt.Sprintf("%s", rebnn))
			}
		}
	}
	log.Print("t5")
	//display success message.
	infoInner := InfoStructSuccessInner{
		ResID: _res_id,
	}
	info := InfoStructSuccess{
		Status: 0,
		Data:   infoInner,
	}
	bnn, err := json.Marshal(info)
	fmt.Fprintf(w, fmt.Sprintf("%s", bnn))

}
