package serializing

import(
	"bytes"
	"encoding/gob"
)
/*
适用类型：二进制到struct相互转换

使用方法：
    b, err := serializing.Encode(data)  
    if err != nil {  
       //错误处理 
    }  
    if err := serializing.Decode(b, &to); err != nil {  
        //错误处理
    }

*/

/*
 功能：序列化
*/
func Encode(data interface{})([]byte, error){
    	buf := bytes.NewBuffer(nil)  
        enc := gob.NewEncoder(buf)  
        err := enc.Encode(data)  
        if err != nil {  
            return nil, err  
        }  
        return buf.Bytes(), nil 
}

/*
 功能：反序列化
*/
func Decode(data []byte,to interface{}) error{
   buf := bytes.NewBuffer(data)  
    dec := gob.NewDecoder(buf)  
    return dec.Decode(to)
}