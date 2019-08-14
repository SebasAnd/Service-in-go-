package main
import (
    "database/sql"
    "net/http"
	"log"
	"fmt"
	"time"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors" 
	"encoding/json"
	"io/ioutil"
	"strings"
	_ "github.com/lib/pq"
)

/*

ANDRES CORTES SERVICE 

*/

const (
  host     = "localhost"
  port     = 5432
  user     = "postgres"
  password = "anddan"
  dbname   = "postgres"
)


type Rows struct {
    Array []string `json:"rows,omitempty"`
}

type ElementLeft struct{
	Server_changed bool `json:"server_changed",omitempty`
	Ssl_grade_left string `json:"ssl_grade_left",omitempty`
	Ssl_grade_hour_before string `json:"ssl_grade_hour_before",omitempty`
	Is_Down bool `json:"is_down",omitempty`
}
type Server struct {
   Address string `json:"address",omitempty`
   Ssl_grade string `json:"ssl_grade",omitempty`
   Country string `json:"country",omitempty`
   Owner string `json:"owner",omitempty`
}
type Server_inBD struct {
   Address string `json:"address",omitempty`
   Ssl_grade string `json:"ssl_grade",omitempty`
   Country string `json:"country",omitempty`
   Owner string `json:"owner",omitempty`
   Time time.Time `json:"domain,omitempty"`
   Domain string `json:"domain,omitempty"`
}

type ServerJson struct{
	Servers []Server `json:"servers",omitempty`
	Server_changed bool `json:"server_changed",omitempty`
	Ssl_grade string `json:"ssl_grade",omitempty`
	Previous_ssl_grade string `json:"previous_ssl_grade",omitempty`
	Logo string `json:"logo",omitempty`
	Title string `json:"title",omitempty`
	Is_Down bool `json:"is_down",omitempty`

}

/*
initialice de database conection( in this case postgres because i had a problem with cockroachdb)

*/


func initialiseDatabase(filepath string) *sql.DB {
    db, err := sql.Open("postgres",filepath)
	
    if err != nil || db == nil {
        panic("Error connecting to database")
    }
    return db
}

/*
insert in db the domain used if this is not in db already

*/




func insert_servers_table(db *sql.DB, domain string) {
   sql :=  `
    SELECT register FROM servers WHERE register = $1;`
   
     test , err := db.Query(sql, domain)
    if err != nil {
        
    }
	
	register:= ""

	for test.Next() {
	err := test.Scan(&register)
	if err != nil {
		fmt.Printf("")
	}

	
}
	
	err = test.Err()
	if err != nil {
	fmt.Println(err)
	}

	if (strings.Compare(register,"")== 0){
	verifier :=  `
    INSERT INTO servers (register)
	VALUES ($1);`
   
		_, err2 := db.Exec(verifier, domain)
		if err2 != nil {
        panic(err2)
		}
}   
}
/*

get the domains in db and retrun them in json form

*/

func listRows(db *sql.DB)[]byte{
	sql :=  `
    SELECT register FROM servers;`
   
     test , err := db.Query(sql)
    if err != nil {
        
    }
	
	register:= ""
	 s := make([]string, 0)
	for test.Next() {
	err := test.Scan(&register)
	s = append(s,register)
	if err != nil {
		fmt.Printf("")
	}
	
	}
	
	rowsList := Rows{
    Array : s,
	}

	btResulttest, _ := json.MarshalIndent(rowsList, "", "  ")
	return btResulttest

	}

/*
save the servers information in db this information is not in db already 

*/

func putinfo_in_DB(domain string,db *sql.DB,array []Server)ElementLeft{

	element_return := ElementLeft{
		Server_changed : false,
		Ssl_grade_left : "",
		Ssl_grade_hour_before: "",
		Is_Down : false,

		}

	sql_form:= `SELECT * FROM servers_information WHERE domain = $1;`
	

	 test , err := db.Query(sql_form, domain)
    if err != nil {
        
    }

	register:= ""
	 s := make([]string, 0)
	for test.Next() {
	err := test.Scan(&register)
	s = append(s,register)
	if err != nil {
		fmt.Printf("")
	}
	
	}
	if len(s) == 0{
		
		for i := 0; i < len(array); i++{

			insert_in_db :=  `
    		INSERT INTO servers_information (address,ssl_grade,country,owner,time_saved,domain)
			VALUES ($1,$2,$3,$4,$5,$6);`
   
			_, err2 := db.Exec(insert_in_db,array[i].Address,array[i].Ssl_grade,array[i].Country,array[i].Owner,time.Now(),domain)
			if err2 != nil {
        	panic(err2)
			}
		}

		element_return = ElementLeft{
		Server_changed : false,
		Ssl_grade_left : minimun_value_ssl(array),
		Ssl_grade_hour_before: "The domain was not in the data base",
		Is_Down : exists_server(array),

		}
		return element_return
	}else{

		get_data_db:= `SELECT address,ssl_grade,country,owner,time_saved,domain FROM servers_information WHERE domain = $1;`
	

	 	elements , err := db.Query(get_data_db, domain)
    	if err != nil {
        
    	}
    	address := ""
		ssl_grade  := ""
		country  := ""
		owner  := ""
		time_saved := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
		domain2  := "" 
	 	arraydb := make([]Server_inBD, 0)
		for elements.Next() {
			err := elements.Scan(&address,&ssl_grade,&country,&owner,&time_saved,&domain2)

			if err != nil {
				fmt.Printf("error")
			}
			//log.Println(address, ssl_grade)
			arraydb = append(arraydb,Server_inBD{Address: address,Ssl_grade: ssl_grade,Country: country,Owner: owner,Time : time_saved, Domain:domain2})
			
		}

		
		date_principal:= time.Now()
		timeStampString := date_principal.Format("2006-01-02 15:04:05")    
    	layOut := "2006-01-02 15:04:05"    
    	timeStamp, err := time.Parse(layOut, timeStampString)

		//date_now := time.Date(date_principal.Year(), date_principal.Month, date_principal.Day(), date_principal.Hours(), date_principal.Minutes(), date_principal.Seconds(),date_principal.Nanoseconds(), time.UTC)
		date_diference := timeStamp.Sub(arraydb[0].Time)

		//fmt.Println(date_diference, arraydb[0].Time,timeStamp)

		if date_diference - 3600000000000 >= 0{
			

			element_return = ElementLeft{
			Server_changed : compare_element_with_db(arraydb,array),
			Ssl_grade_left : minimun_value_ssl(array),
			Ssl_grade_hour_before: minimun_value_ssl_db(arraydb),
			Is_Down : exists_server(array),

			}
		}else{


			element_return = ElementLeft{

			Server_changed : compare_element_with_db(arraydb,array),
			Ssl_grade_left : minimun_value_ssl(array),
			Ssl_grade_hour_before: "The domain has been in the database since a few minutes ago.",
			Is_Down : exists_server(array)}
		}
		

	}
	return element_return 

}

/*compare the servers obatin in the recent search with the information in db */

func compare_element_with_db(array[]Server_inBD,coming_array[]Server)bool{
	result := false

	for i := 0; i < len(array); i++{

		if strings.Compare(array[i].Address,coming_array[i].Address) == -1||
		strings.Compare(array[i].Ssl_grade,coming_array[i].Ssl_grade) == -1||
		strings.Compare(array[i].Country,coming_array[i].Country) == -1||
		strings.Compare(array[i].Owner,coming_array[i].Owner) == -1{
			
			result = true

		}

			}

	return result
}
/*
verify if the domain has at least one ip direction, if the domain had ip direction means that the server is not down

*/
func exists_server (array[]Server)bool{

	ret_bool := true 

	for i := 0; i < len(array); i++{

		if strings.Compare(array[i].Address, "No Data") == -1{

			ret_bool = false
		}
	}

	return ret_bool
}


/*

compare the element with the elements of a list and return a number according to element location 

*/
func index_in_array(index[]string, element string)int{
	ret := -1
	

	for i := 0; i < len(index); i++{

		if strings.Compare(index[i], element) == 0{

			ret = i
			
		}
	}
	return ret
}
/*

function to ge the minimum value of a server ssl grade of the domain searched

*/
func minimun_value_ssl(array []Server)string{
	index_values := []string {"F","E","D","C","B","A","A+"}
	ret_string := "A+" 
	

	for i := 0; i < len(array); i++{
		
		if index_in_array(index_values,ret_string) > index_in_array(index_values,array[i].Ssl_grade) && strings.Compare(array[i].Ssl_grade,"No Data") == -1 {

			ret_string = array[i].Ssl_grade
		}

	}

	return ret_string

}

/*

function to ge the minimum value of a server ssl grade of the domain in db

*/
func minimun_value_ssl_db(array []Server_inBD)string{
	index_values := []string {"F","E","D","C","B","A","A+"}
	ret_string := "A+" 
	for i := 0; i < len(array); i++{

		if index_in_array(index_values,ret_string) > index_in_array(index_values,array[i].Ssl_grade)&& strings.Compare(array[i].Ssl_grade,"No Data")== -1{

			ret_string = array[i].Ssl_grade
		}

	}

	return ret_string

}



func createjson(domain string,db *sql.DB)[]byte{


/*
get the server elements from the url given 

*/
	serv := ("https://api.ssllabs.com/api/v3/analyze?host="+domain)
	rest, err2 := http.Get(serv)

/*
get the server's answer or the error resulted

*/

    if err2 != nil {
        panic(err2)
    }

/*

handle the error if exists one 
*/
   
	var servdat map[string]interface{}

/*

interface created to save the return of the server 

*/

	err := json.NewDecoder(rest.Body).Decode(&servdat)
	if err2 != nil{
	log.Fatal(err2)
	}
	/*

handle the error if exists one 
if the error not exits save the elemets obtain from the servers in the servdat interface 
*/
   
    test := servdat["endpoints"].([]interface{})

 /*

find in the servdat interface the element endpoints al returnit as array interface


*/
	
	n := make([]Server,len(test))
 /*

create an array the type Server initially null with test size


*/	
	for a := 0; a < len(n); a++ {
	var serversinfo = servdat["endpoints"].([]interface{})[a]
	
	if(serversinfo.(map[string]interface{})["ipAddress"] == nil){

	
			n[a] = Server{Address : "No Data",
   			Ssl_grade : "No Data",
   			Country : "Not Data",
   			Owner : "Not Data"}
		
	
	}
	if(serversinfo.(map[string]interface{})["grade"] == nil){

		
		whoisurl := ("https://www.whoisxmlapi.com/whoisserver/WhoisService?apiKey=at_Xj70i3fUA41Mxnvpg6eYQvwZKhr0s&domainName="+serversinfo.(map[string]interface{})["ipAddress"].(string)+"&outputFormat=json")
	

		res, err := http.Get(whoisurl)
    	if err != nil {
        panic(err)
    	}
   
		var dat map[string]interface{}

		err = json.NewDecoder(res.Body).Decode(&dat)
		if err != nil{
		log.Fatal(err)
		}



		if(dat["WhoisRecord"].(map[string]interface{})["registryData"].(map[string]interface{})["registrant"] == nil){

			n[a] = Server{Address : serversinfo.(map[string]interface{})["ipAddress"].(string),
  			Ssl_grade : "No Data",
   			Country : "Not Data",
   			Owner : "Not Data"}
		}else{

			n[a] = Server{Address : serversinfo.(map[string]interface{})["ipAddress"].(string),
   			Ssl_grade : "Not Data",
   			Country : dat["WhoisRecord"].(map[string]interface{})["registryData"].(map[string]interface{})["registrant"].(map[string]interface{})["country"].(string),
   			Owner : dat["WhoisRecord"].(map[string]interface{})["registryData"].(map[string]interface{})["registrant"].(map[string]interface{})["organization"].(string)}
		}
	
	}

	if(serversinfo.(map[string]interface{})["grade"] != nil && serversinfo.(map[string]interface{})["ipAddress"] != nil){

		
		whoisurl := ("https://www.whoisxmlapi.com/whoisserver/WhoisService?apiKey=at_Xj70i3fUA41Mxnvpg6eYQvwZKhr0s&domainName="+serversinfo.(map[string]interface{})["ipAddress"].(string)+"&outputFormat=json")
	

		res, err := http.Get(whoisurl)
    	if err != nil {
        panic(err)
    	}
   
		var dat map[string]interface{}

		err = json.NewDecoder(res.Body).Decode(&dat)
		if err != nil{
		log.Fatal(err)
		}



		if(dat["WhoisRecord"].(map[string]interface{})["registryData"].(map[string]interface{})["registrant"] == nil){

			n[a] = Server{Address : serversinfo.(map[string]interface{})["ipAddress"].(string),
  			Ssl_grade : serversinfo.(map[string]interface{})["grade"].(string),
   			Country : "Not Data",
   			Owner : "Not Data"}
		}else{

			n[a] = Server{Address : serversinfo.(map[string]interface{})["ipAddress"].(string),
   			Ssl_grade : serversinfo.(map[string]interface{})["grade"].(string),
   			Country : dat["WhoisRecord"].(map[string]interface{})["registryData"].(map[string]interface{})["registrant"].(map[string]interface{})["country"].(string),
   			Owner : dat["WhoisRecord"].(map[string]interface{})["registryData"].(map[string]interface{})["registrant"].(map[string]interface{})["organization"].(string)}
		}
	
	}
	

	}
	/*

loop to fill the n array with the servers information and whois api information 



	*/

   
   url := "http://"+domain
	resp, err := http.Get(url)

	if err != nil {
		panic(err)
	}

/*
handle the error if exist one 
if the error dont exist get the information of a page using his url


	*/


	defer resp.Body.Close()
	

	
	

	dataInBytes, err := ioutil.ReadAll(resp.Body)
    pageContent := string(dataInBytes)
	evaluatetitle:= false

/*
use the html information as a string 
	*/

    titleStartIndex := strings.Index(pageContent, "<title>")
    if titleStartIndex == -1 {
        fmt.Println("No title element found")
		evaluatetitle = true
        
    }
/*
find the <title> element in the html string, 

 */
    
    titleStartIndex += 7

/*
to exclude the element title found already, the search of the closing tag </title> start with a  bigger index 


*/
	
    // Find the index of the closing tag
    titleEndIndex := strings.Index(pageContent, "</title>")
    if titleEndIndex == -1 {
        fmt.Println("No closing tag for title found.")
        evaluatetitle = true
    }

    /*
    find the title closing tag 


*/
	
  
	pageTitle := " "
	if (evaluatetitle){
		pageTitle = " " 
	}else{
		pageTitle = string([]byte(pageContent[titleStartIndex:titleEndIndex]))
	}
    
/*

if the element was found page title have the title element or  the page title is a string empty

*/



	iconurl := ("https://besticon-demo.herokuapp.com/allicons.json?url="+domain)
	rest2, err2 := http.Get(iconurl)
    if err2 != nil {
        panic(err2)
    }

  /*
 get the icon of a domain using another api

  */

	var servicon map[string]interface{}

	 err3 := json.NewDecoder(rest2.Body).Decode(&servicon)
	if err3 != nil{
	log.Fatal(err3)
	}
 /*
as the other api manager put the obtain information in an interface 

handle the error if exists

  */



	var iconinfo = servicon["icons"].([]interface{})[0].(map[string]interface{})["url"]
	
 /*
 get the the icons url from the api 

  */
	other_information := putinfo_in_DB(domain,db,n)

 /*
 get the information left to create el general json 

  */

	serverjson := ServerJson{Servers : nil,
	Server_changed: other_information.Server_changed,
	Ssl_grade : other_information.Ssl_grade_left,
	Previous_ssl_grade: other_information.Ssl_grade_hour_before,
	Logo : "" ,
	Title : "",
	Is_Down: other_information.Is_Down}

	if iconinfo == nil{
	serverjson = ServerJson{Servers : n,
	Server_changed: other_information.Server_changed,
	Ssl_grade : other_information.Ssl_grade_left,
	Previous_ssl_grade: other_information.Ssl_grade_hour_before,
	Logo : "" ,
	Title : pageTitle,
	Is_Down: other_information.Is_Down}
	} else {
	serverjson = ServerJson{Servers : n,
	Server_changed: other_information.Server_changed,
	Ssl_grade : other_information.Ssl_grade_left,
	Previous_ssl_grade: other_information.Ssl_grade_hour_before,
	Logo : iconinfo.(string),
	Title : pageTitle,
	Is_Down: other_information.Is_Down}

	}
/*
 create the general struct with all necesary information 

  */
	
   
   
	btResulttest, _ := json.MarshalIndent(serverjson, "", "  ")
/*
create the json accoding to the struct 

  */

	insert_servers_table(db, domain)
/*
evaluate and insert the domain information to the database

*/
	 return btResulttest
/*
 return the final json to show it when the method is used
  */
}


func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s sslmode=disable",
    host, port, user, password, dbname)
/*
 get the database information 

  */
	db, err := sql.Open("postgres", psqlInfo)

/*
open the data base conection 


*/
	if err != nil {
	panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
	panic(err)
	}
/*
 handle the error if exists.
if the error does not exists means that the connection was successfully

  */
	

	r := chi.NewRouter()
/*
 create the chi router to show the elements using the navegador 

  */

	cors := cors.New(cors.Options{
    AllowedOrigins:   []string{"*"},
   
    AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
    ExposedHeaders:   []string{"Link"},
    AllowCredentials: true,
    MaxAge:           300, 
  })
/*
cors was created to vuejs, because if this information is not handle vuejs doesn't let you use the server information 

  */

  r.Use(cors.Handler)
  fmt.Println("welcome, write  http://localhost:3000/search/(web site domain that you want search) on your navigator")
		fmt.Println("or write  http://localhost:3000/search/listsearch on your navigator")
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome, write/search/(web site domain)"))

		
		/*
first page show without a specified route

  */
	})
	r.Route("/search", func(r chi.Router) {
	 r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("write a web domain to continue"))

		/*
/search page show without a specified domain to search 

  */
		
	})
      r.Get("/{param}",func(w http.ResponseWriter, r *http.Request) {
		param := chi.URLParam(r, "param") 
/*
/search/{param} page show the result of the search of the given domain 

  */		
		fmt.Println("app running domain:"+ param)
		w.Write(createjson(param,db))
		
	})
	r.Get("/listsearch",func(w http.ResponseWriter, r *http.Request) {
				
		w.Write(listRows(db))

		/*
/listsearch page show the list of searches

  */
		
	})
	  
	  
    })
	 
	  
	http.ListenAndServe(":3000", r)

	/*
define what port will be the server listening in this case will be http://localhost:3000/search/listsearch

  */
}