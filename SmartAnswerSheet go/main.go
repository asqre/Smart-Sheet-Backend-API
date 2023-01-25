package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"sync"
)

//backticks ;- ``
//NOTE: To convert the response of all API request in JSON format
// use "w.Header().Add("Content-Type", "application/json")"   means add these to headers which shows content type in JSON format

// Declare enum
const (
	AddSubmission          = "/addSubmission"
	GetSubmissionById      = "/getSubmissionById"
	GetAllSubmissionsByIds = "/getAllSubmissionsByIds"
	GetAllSubmissions      = "/getAllSubmissions"
	UpdateSubmissionById   = "/updateSubmissionById"
	DeleteSubmissionById   = "/deleteSubmissionById"
	GetUrlRequestCounter   = "/getUrlRequestCounter"
	//FileHandler            = "/fileHandler"
)

// FileDetails NOTE: Refactor use: it actually renames whole (same type of word).just right-click on that word and click refactor and rename.
type FileDetails struct {
	//we write their attributes in json format to use POST in postman
	FileName      string `json:"file_name"`
	LocalFileName string `json:"local_file_name"`
	FileSize      string `json:"file_size"`
	FilePath      string `json:"file_path"`
}

type Url struct {
	Url string `json:"url"`
}

// StudentInfo creating a map that will store student detail having key id
var StudentInfo = make(map[string]StudentDetails)

// create a global variable url channel
var urlChanel = make(chan string, 7)

type studentId struct {
	Id string
}

// creating a struct for decoding multi string in getAllSubmissionById
type studentIdArray struct {
	NoOfId []string `json:"noOfId"`
}

type StudentDetails struct {
	Id                  string      `json:"id"`
	StudentName         string      `json:"student_name"`
	StudentEmailAddress string      `json:"student_email_address"`
	File                FileDetails `json:"file"`
	TimeStamp           int64       `json:"time_stamp"`
	// we consider timestamp datatype as int64 coz,we use time.now.unix() to find the time, and it can be calculated easily
	// Current Unix Time is 1672722393 seconds since 1 January 1970. that's why we take int64 instead of int8, int16...
}

// creating a map for counting URL request
//var cnt = make(map[string]int)

// CriticalSection now , declare mutex(mutual exclusion)  to avoid any two or more server conflict. like addSubmission and getSubmission are not call at a time.
// that's why we use mutex
// to use mutex, we should declare struct coz, structure is the medium to marshal/unmarshal
type CriticalSection struct {
	mux         sync.Mutex
	cnt         map[string]int
	StudentData map[string]StudentDetails
	//we use this cnt, and studentsData in this struct, these are used in every API,
}

// now creating a critical section object globally  because we have to use in every API
var criticalSection = CriticalSection{
	//now , initialize the object
	//no need to initialize the mutex , it will be automatically initialized
	cnt:         make(map[string]int), //now initialize the variable of struct by :
	StudentData: make(map[string]StudentDetails),
}

func addSubmission(w http.ResponseWriter, r *http.Request) {

	//Resolving "Error: XMLHttpRequest error." in frontend terminal
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	log.Println(r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	//

	//convert the response of api in json format
	w.Header().Add("Content-Type", "application/json")
	//count the api called using channel, which can handle multiple request
	urlChanel <- AddSubmission
	//lock, update and unlock the API
	//criticalSection.mux.Lock()
	//criticalSection.cnt[AddSubmission]++
	//criticalSection.mux.Unlock()
	//check 1st request is post or not
	//cnt[AddSubmission]++
	if r.Method != "POST" {
		http.Error(w, " Only POST request are Allowed", http.StatusMethodNotAllowed) // statusMethodNotAllowed coz method are not post.that's why
		return
	}
	// To read a header from a http request in golang then, we use "r.Header.Get" if I receive a request of type http.Request

	// Now take data in Json format which is encoded (i.e. go Marshal)
	var studentDetail StudentDetails
	//now decode it to studentDetail (i.e, go UnMarshal) from request body which is in byte array(int data type)
	err := json.Unmarshal([]byte(r.FormValue("data")), &studentDetail) //decode is only possible for struct. so , studentDetail struct are used
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// as file is open, and it must be closed at the end .so, use defer it stored in stack and executed at the end
	defer func(file multipart.File) { //there may be an error while close the file so, we need to handle.so, simply return
		err = file.Close()
		if err != nil {
			return
		}
	}(file)
	//create a new file to store uploaded file
	out, err := os.Create(header.Filename) //to store large file we use "os"
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) //if developer is unable to create a new file then the error is from developer side.so internal servererror
		return
	}
	defer func(out *os.File) { //it will take pointer file as an input
		err = out.Close()
		if err != nil {
			return
		}
	}(out)
	//now, copy the uploaded file to the new file
	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = fmt.Fprintln(w, "File Submitted Successfully")
	if err != nil {
		return
	}

	// now we lock the one API to access share data and hold the remaining API. so, we use mutex
	criticalSection.mux.Lock()
	// now , update the studentDetail by locking another api
	criticalSection.StudentData[studentDetail.Id] = studentDetail
	// now , unlock the api
	criticalSection.mux.Unlock()
	//   now create a map of StudentInfo
	//[studentDetail.Id] = studentDetail // as it is updated 4 line above
	_, err = fmt.Fprintf(w, "Submitted Successfully of %v \n", studentDetail.Id)
	if err != nil {
		return
	}

}

func getSubmissionById(w http.ResponseWriter, r *http.Request) {
	//Resolving "Error: XMLHttpRequest error." in frontend terminal
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	log.Println(r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	//
	w.Header().Add("Content-type", "application/json")
	urlChanel <- GetSubmissionById
	//same here as above API, lock, update and unlock the API
	//criticalSection.mux.Lock()
	//criticalSection.cnt[GetSubmissionById]++
	//criticalSection.mux.Unlock()
	if r.Method != "POST" {
		http.Error(w, "Only POST request are Allowed", http.StatusMethodNotAllowed) //here, when wrong request are used then statusmethodnotallowed are used
		return
	}
	// now create a studentId struct to use decode(or unMarshal) for Using particular id
	var studentDetail studentId
	err := json.NewDecoder(r.Body).Decode(&studentDetail)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) //in this error statusBadRequest we do .coz,given request are not fond
		//here, when the error is coming from client side then statusBadRequest is used.
		return
	}
	var getDetailsByID StudentDetails
	//similarly , lock, update,unlock the api
	criticalSection.mux.Lock()
	getDetailsByID = criticalSection.StudentData[studentDetail.Id]
	criticalSection.mux.Unlock()
	//getDetailsByID = StudentInfo[studentDetail.Id]
	_, err = fmt.Fprintf(w, "submission detail of %v are %v\n", studentDetail.Id, getDetailsByID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) //here, when the error is coming due to developer/or from mys ide
		// then StatusInternalServerError is used
		return
	}

}

func getAllSubmissionsByIds(w http.ResponseWriter, r *http.Request) {
	//Resolving "Error: XMLHttpRequest error." in frontend terminal
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	log.Println(r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	//
	w.Header().Add("Content-type", "application/json")
	urlChanel <- GetAllSubmissionsByIds
	//criticalSection.mux.Lock()
	//criticalSection.cnt[GetAllSubmissionsByIds]++
	//criticalSection.mux.Unlock()
	//cnt[GetAllSubmissionsByIds]++
	//as we have to return student detail of set of id request given by client so decode/ unmarshal it
	if r.Method != "POST" {
		http.Error(w, "Only POST request are Allowed", http.StatusMethodNotAllowed)
		return
	}
	// client give the array of id, so we have to create a struct of array instead of using direct string. for decode
	var idDetails studentIdArray
	err := json.NewDecoder(r.Body).Decode(&idDetails)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//for this we have to traverse through the map and return
	for i, val := range idDetails.NoOfId { //we get v as a particular id and access details by map
		criticalSection.mux.Lock()
		_, err = fmt.Fprintf(w, "Details of %vst student are %v \n", i+1, criticalSection.StudentData[val])
		criticalSection.mux.Unlock()
		//_, err = fmt.Fprintf(w, "Details of %vst student are %v \n", i+1, StudentInfo[val])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func getAllSubmissions(w http.ResponseWriter, r *http.Request) {
	//Resolving "Error: XMLHttpRequest error." in frontend terminal
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	log.Println(r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	//
	w.Header().Add("Content-type", "application/json")
	urlChanel <- GetAllSubmissions
	//criticalSection.mux.Lock()
	//criticalSection.cnt[GetAllSubmissions]++
	//criticalSection.mux.Unlock()
	//cnt[GetAllSubmissions]++
	//get method is used
	if r.Method != "GET" {
		http.Error(w, "Only GET Request are Allowed", http.StatusMethodNotAllowed)
		return
	}
	err := json.NewEncoder(w).Encode(criticalSection.StudentData)
	//_, err := fmt.Fprintf(w, "All availale data of students are %v\n", criticalSection.StudentData)
	//_, err := fmt.Fprintf(w, "All available data of students are %v \n", StudentInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func updateSubmissionById(w http.ResponseWriter, r *http.Request) {
	//Resolving "Error: XMLHttpRequest error." in frontend terminal
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	log.Println(r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	//
	w.Header().Add("Content-type", "application/json")
	urlChanel <- UpdateSubmissionById
	//criticalSection.mux.Lock()
	//criticalSection.cnt[UpdateSubmissionById]++
	//criticalSection.mux.Unlock()
	if r.Method != "POST" {
		http.Error(w, "Only POST request are Allowed", http.StatusMethodNotAllowed)
		return
	}
	// we get the detail of one student, and we have to update it
	//same decode/unmarshal the json format of student detail
	var studentDetail StudentDetails
	err := json.NewDecoder(r.Body).Decode(&studentDetail)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// now, we have a detail of student which we have to update by their ID
	criticalSection.mux.Lock()
	criticalSection.StudentData[studentDetail.Id] = studentDetail
	criticalSection.mux.Unlock()
	//StudentInfo[studentDetail.Id] = studentDetail
	criticalSection.mux.Lock()
	_, err = fmt.Fprintf(w, "Updated Detail of required %v ID are %v\n", studentDetail.Id, criticalSection.StudentData[studentDetail.Id])
	criticalSection.mux.Unlock()
	//_, err = fmt.Fprintf(w, "Updated Detail of required %v ID are %v\n", studentDetail.Id, StudentInfo[studentDetail.Id])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func deleteSubmissionById(w http.ResponseWriter, r *http.Request) {
	//Resolving "Error: XMLHttpRequest error." in frontend terminal
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	log.Println(r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	//
	w.Header().Add("Content-type", "application/json")
	urlChanel <- DeleteSubmissionById
	//criticalSection.mux.Lock()
	//criticalSection.cnt[DeleteSubmissionById]++
	//criticalSection.mux.Unlock()
	//cnt[DeleteSubmissionById]++
	if r.Method != "POST" {
		http.Error(w, "Only POST request are allowed ", http.StatusMethodNotAllowed)
		return
	}
	//let we get the no of id, so we have to create ID struct(for decode/marshal use)
	var deletionID studentId
	err := json.NewDecoder(r.Body).Decode(&deletionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//now, delete the particular id ny delete() function
	delete(StudentInfo, deletionID.Id)
	_, err = fmt.Fprintln(w, "required student ID deleted")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func getUrlRequestCounter(w http.ResponseWriter, r *http.Request) {
	//Resolving "Error: XMLHttpRequest error." in frontend terminal
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	log.Println(r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	//
	w.Header().Add("Content-type", "application/json")
	if r.Method != "POST" {
		http.Error(w, "Only POST request are Accessible", http.StatusMethodNotAllowed)
		return
	}
	var url Url
	err := json.NewDecoder(r.Body).Decode(&url)
	switch url.Url {
	case AddSubmission:
		criticalSection.mux.Lock()
		_, err := fmt.Fprintln(w, "No. of Times called addSubmission ", criticalSection.cnt[AddSubmission])
		criticalSection.mux.Unlock()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case GetSubmissionById:
		criticalSection.mux.Lock()
		_, err = fmt.Fprintf(w, "No. of Times called GetSubmissionById %v", criticalSection.cnt[GetSubmissionById])
		criticalSection.mux.Unlock()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case GetAllSubmissionsByIds:
		criticalSection.mux.Lock()
		_, err = fmt.Fprintf(w, "No. of Times called GetAllSubmissionsByIds %v", criticalSection.cnt[GetAllSubmissionsByIds])
		criticalSection.mux.Unlock()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case GetAllSubmissions:
		criticalSection.mux.Lock()
		_, err = fmt.Fprintf(w, "No. of Times called GetAllSubmissions %v", criticalSection.cnt[GetAllSubmissions])
		criticalSection.mux.Unlock()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case UpdateSubmissionById:
		criticalSection.mux.Lock()
		_, err = fmt.Fprintf(w, "No. of Times called UpdateSubmissionById %v", criticalSection.cnt[UpdateSubmissionById])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case DeleteSubmissionById:
		criticalSection.mux.Lock()
		_, err = fmt.Fprintf(w, "No. of Times called DeleteSubmissionById %v", criticalSection.cnt[DeleteSubmissionById])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// In Go, the sync.WaitGroup type is used to wait for a collection of goroutines to finish
// creating new wait group
var wg sync.WaitGroup

// UrlCounter creating function to run independently in go routine
func UrlCounter(ch chan string) {
	//close the channel at last
	//we have to close multiple things (close(ch)),wg.Done at a time.so, we make a func and done it.
	defer func() {
		close(ch) //it closes the channel
		wg.Done() //indicate that goroutine is done
		//this tells about go routine is done to main thread
	}()
	//make infinite loop to take receive  multiple request
	for {
		//take value from channel
		var api = <-ch
		//lock,update,delete the particular api
		criticalSection.mux.Lock()
		criticalSection.cnt[api]++
		criticalSection.mux.Unlock()
	}
}

func main() {
	//creating channel
	//c := make(chan int)
	//go func() {
	//	c <- 42
	//}()
	//val := <-c
	//fmt.Println(val)
	////creating a goroutine
	//go func() {
	//	fmt.Println("Hello")
	//}()
	//fmt.Println("world")

	http.HandleFunc(AddSubmission, addSubmission)
	http.HandleFunc(GetSubmissionById, getSubmissionById)
	http.HandleFunc(GetAllSubmissionsByIds, getAllSubmissionsByIds)
	http.HandleFunc(GetAllSubmissions, getAllSubmissions)
	http.HandleFunc(UpdateSubmissionById, updateSubmissionById)
	http.HandleFunc(DeleteSubmissionById, deleteSubmissionById)
	http.HandleFunc(GetUrlRequestCounter, getUrlRequestCounter)
	//http.HandleFunc(FileHandler, fileHandler)
	//add 1 go routine to the wait-group
	wg.Add(1)
	//goroutines function call
	go UrlCounter(urlChanel)

	//fmt.Println("Server is up")
	err := http.ListenAndServe(":8080", nil)
	close(urlChanel)
	if err != nil {
		log.Fatal("Listen And Serve :8080", nil)
	}
	//close the channel
	//wait for goroutine to finish
	wg.Wait()

}
