package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// https://github.com/PacktPublishing/Hands-On-Dependency-Injection-in-Go
// ch03/01_optimizing_for_humans

func NotSoSimple(ID int64, name string, age int, registered bool) string {
	out := &bytes.Buffer{}
	out.WriteString(strconv.FormatInt(ID, 10))
	out.WriteString("-")
	out.WriteString(strings.Replace(name, " ", "-", -1))
	out.WriteString("-")
	out.WriteString(strconv.Itoa(age))
	out.WriteString("-")
	out.WriteString(strconv.FormatBool(registered))
	return out.String()
}

func Simpler(ID int64, name string, age int, registered bool) string {
	nameWithNoSpace := strings.Replace(name, " ", "-", -1)
	return fmt.Sprintf("%d-%s-%d-%t", ID, nameWithNoSpace, age, registered)
}

type myGatter interface {
	Get(url string) (*http.Response, error)
}

func TooAbstract(getter myGatter, url string) ([]byte, error) {
	resp, err := getter.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func CommonConcept(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

// 05_boolean_param
type Pet struct {
	Name string
	Dog  bool
	Born time.Time
}

func NewPet(name string, isDog bool) Pet {
	return Pet{
		Name: name,
		Dog:  isDog,
		Born: time.Now(),
	}
}

func CreatePetsV1() {
	NewPet("Fido", true)
}

// 06_hidden_boolean.go
const (
	isDog = true
	isCat = false
)

func NewDog(name string) Pet {
	return NewPet(name, isDog)
}

func NewCat(name string) Pet {
	return NewPet(name, isCat)
}

func CreatePetsV2() {
	NewDog("Fido")
}

// 07_wide_formatter.go
type WideFormatter interface {
	ToCSV(pets []Pet) ([]byte, error)
	ToGOB(pets []Pet) ([]byte, error)
	ToJSON(pets []Pet) ([]byte, error)
}

// 08_thin_formatters.go
type ThinFormatter interface {
	Format(pets []Pet) ([]byte, error)
}

type CSVFormatter struct{}

func (f CSVFormatter) Format(pets []Pet) ([]byte, error) {
	// convert slice of pets to CSV
	return nil, nil
}

type GOBFormatter struct{}

func (f GOBFormatter) Format(pets []Pet) ([]byte, error) {
	return nil, nil
}

type JSONFormatter struct{}

func (f JSONFormatter) Format(pets []Pet) ([]byte, error) {
	// convert slice of pets to JSON
	return nil, nil
}

// 09_extra_config.go
func PetFetcher(search string, limit int, offset int, sortBy string, sortAscending bool) []Pet {
	return []Pet{}
}

func PetFetcherTypicalUsage() {
	_ = PetFetcher("Fido", 0, 0, "", true)
}

func main() {

}
