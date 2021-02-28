package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func main() {

	fmt.Println("Running on 8080: ")

	r := mux.NewRouter()

	r.HandleFunc("/api", ServeHTTP)

	http.ListenAndServe(":8080", r)

}

type InitialRequest struct {
	//{businesses: Array(20), total: 291, region: {…}}
	Total int `json:"total"`
}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")

	st := r.URL.Query().Get("searchText")

	//Check it once to get the count for the offset
	offset, err := fetchParkingData(w, r, st)
	if err != nil {
		fmt.Println("error fetching yelp data")
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	makeFinalCall(w, r, offset, st)

}

func makeFinalCall(w http.ResponseWriter, r *http.Request, offset int, st string) {

	client := http.Client{}

	req, err := http.NewRequest("GET", "https://api.yelp.com/v3/businesses/search", nil)
	if err != nil {
		fmt.Println("error")
		fmt.Println(err)
		return
	}

	req.Header.Add("Authorization", "Bearer LUOOZvpP-MLYi9phUl5tJCf-CBDMjRxMap0_QPgkDXgaiccqKVL3_anu1acLW0vaNbdjlVwC4zjy20Flvum3HTFIDo6ct_fiq5l81R-D5WG5V2eHXgj11kOcnAs7YHYx")

	stO := strconv.Itoa(offset)

	q := req.URL.Query()
	q.Add("term", "parking")
	q.Add("location", st)
	q.Add("sort_by", "rating")
	q.Add("offset", stO)
	req.URL.RawQuery = q.Encode()

	res, err := client.Do(req)
	if err != nil {
		fmt.Println("error doing request")
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bs, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("error parsing")
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	io.WriteString(w, string(bs))

}

func calculateOffset(count int) int {
	if count < 50 {
		return 0
	}

	offset := count - 50

	if offset > 1000 {
		return 950
	}

	return offset
}

func fetchParkingData(w http.ResponseWriter, r *http.Request, st string) (int, error) {

	client := http.Client{}

	req, err := http.NewRequest("GET", "https://api.yelp.com/v3/businesses/search", nil)
	if err != nil {
		fmt.Println("error")
		fmt.Println(err)
		return 0, errors.New("err")
	}

	req.Header.Add("Authorization", "Bearer LUOOZvpP-MLYi9phUl5tJCf-CBDMjRxMap0_QPgkDXgaiccqKVL3_anu1acLW0vaNbdjlVwC4zjy20Flvum3HTFIDo6ct_fiq5l81R-D5WG5V2eHXgj11kOcnAs7YHYx")

	q := req.URL.Query()
	q.Add("term", "parking")
	q.Add("location", st)
	q.Add("sort_by", "rating")
	req.URL.RawQuery = q.Encode()

	res, err := client.Do(req)
	if err != nil {
		fmt.Println("error doing request")
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return 0, errors.New("err")
	}

	bs, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("error parsing")
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return 0, errors.New("err")
	}

	ir := InitialRequest{}

	err = json.Unmarshal(bs, &ir)
	if err != nil {
		fmt.Println("error marshaling json")
		fmt.Println(err)
		return 0, errors.New("err")
	}

	//caluclate offset using total
	offset := calculateOffset(ir.Total)

	return offset, nil

}
