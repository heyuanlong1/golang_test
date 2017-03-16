package common
import (
   "math/rand"
   "time"
)

func GetRandomString(lens int) string{
   str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
   bytes := []byte(str)
   result := []byte{}
   r := rand.New(rand.NewSource(time.Now().UnixNano()))
   for i := 0; i < lens; i++ {
      result = append(result, bytes[r.Intn(len(bytes))])
   }
   return string(result)
}

//这方式比较特别，按照123456来记忆吧：01月02号 下午3点04分05秒 2006年
func GetFormatTime(format string ) string{
   return time.Now().Local().Format(format)   
}
