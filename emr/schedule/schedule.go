package schedule

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

type Schedule struct {
		Fname string 
		Lname string 
		Date string 
		Start string 
		End string 
		Status string 
		Comment string
	}
	
func str2js(c Schedule) string {
	cjs, err := json.Marshal(c)	
	if err != nil {
        return "Json failed"
    }
    return string(cjs)
}

func build_sql() string {
	A := "SELECT "	+	"fname, lname, pc_eventDate, pc_startTime, pc_endTime, pc_apptstatus, pc_hometext ";	
	B := "FROM "	+	"patient_data, openemr_postcalendar_events ";
	C := "WHERE "  	+	"pc_aid = '4' ";
	D := "AND " 	+	"pc_eventDate = '2012-05-02' ";
	E := "AND " 	+	"pid = pc_pid";

	sql := A + B + C + D + E
	
	return sql
}

func GetSchedule() string {
	
	user := "tmaniaci"; pass := "password"; dbname := "openemr"; proto := "tcp"
	addr := "dev.itelehome.com:3306"
	db := mysql.New(proto, "", addr, user, pass, dbname)

	fmt.Printf("Connect to %s:%s... ", proto, addr)
	checkError(db.Connect())
	printOK()

	fmt.Println("Select from OpenEMR Users... ")
	rows, res := checkedResult(db.Query(build_sql()))
	
	jstr := "["
	
	for i, row := range rows {
		sch := Schedule{}
		sch.Fname = string(row.Str(res.Map("fname")))
		sch.Lname = string(row.Str(res.Map("lname")))
		sch.Date = string(row.Str(res.Map("pc_eventDate")))
		sch.Start = string(row.Str(res.Map("pc_startTime")))
		sch.End = string(row.Str(res.Map("pc_endTime")))
		sch.Status = string(row.Str(res.Map("pc_apptstatus")))
    	sch.Comment = string(row.Str(res.Map("pc_hometext")))
    	
		jstr += str2js(sch)  	
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