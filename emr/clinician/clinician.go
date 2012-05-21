package clinician

import (
	"os"
	"fmt"
	"encoding/json"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/thrsafe"
)

func printOK() {
	fmt.Println("OK")
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func checkedResult(rows []mysql.Row, res mysql.Result, err error) ([]mysql.Row,
	mysql.Result) {
	checkError(err)
	return rows, res
}

type Clinician struct {
	Fname	string
	Lname	string
}

func str2js(c Clinician) string {
	cjs, err := json.Marshal(c)	
	if err != nil {
        return "Json failed"
    }
    return string(cjs)
}

func GetClinican() string{
	user := "tmaniaci"
	pass := "password"
	dbname := "openemr"
	proto := "tcp"
	addr := "dev.itelehome.com:3306"

	db := mysql.New(proto, "", addr, user, pass, dbname)

	fmt.Printf("Connect to %s:%s... ", proto, addr)
	checkError(db.Connect())
	printOK()

	fmt.Println("Select from OpenEMR Users... ")
	rows, res := checkedResult(db.Query("select fname, lname from users"))
	
	jstr := "["
		
	for i, row := range rows {
		clinician := Clinician{}
		clinician.Fname = string(row.Str(res.Map("fname")))
		clinician.Lname = string(row.Str(res.Map("lname")))	
        
        jstr += str2js(clinician)   	
    	if i < (len(rows)-1) {jstr += ","}  			
	}
	jstr += "]"
	fmt.Println("\njsarray:", jstr)	
	
	printOK()

	fmt.Print("Close connection... ")
	checkError(db.Close())
	printOK()
	
	return jstr
}