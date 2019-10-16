package main

import (
        "bytes"
        "encoding/json"
        "errors"
        "fmt"
        "net/http"
        "strings"
)

func main() {
        http.HandleFunc("/status", Status)
        http.HandleFunc("/run", Run)
        http.ListenAndServe(":3000", nil)
}

func Status(w http.ResponseWriter, r *http.Request) {
        // We are ready to go as long as we can connect to the validator
        _, err := http.Get("https://validator.sbolstandard.org/")
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        } else {
                fmt.Fprintln(w, "Ready to go!")
        }
}

func Run(w http.ResponseWriter, r *http.Request) {
    fasta := ">Z. ach genome\nGATTACATGATTACAGATTACA"

    valid, result, err := ConvertFastaToGenbank(fasta)
    if err != nil {
        fmt.Fprintln(w, err.Error(), http.StatusInternalServerError)
    } else if !valid {
        fmt.Fprintln(w, result, http.StatusUnprocessableEntity)
    } else {
        fmt.Fprintln(w, result)
    }
}

func ConvertFastaToGenbank(fasta string) (bool, string, error) {
        options := ValidatorOptions{"https://fasta.zach.network", "1", false}
        validatorRequest := ValidatorRequest{options, true, fasta}

        requestJson, err := JSONMarshal(validatorRequest)
        if err != nil {
                return false, "", err
        }

        requestBuf := bytes.NewReader(requestJson)

        resp, err := http.Post("https://validator.sbolstandard.org/validate/", "application/json", requestBuf)
        if err != nil {
                return false, "", err
        }
        defer resp.Body.Close()

        if resp.StatusCode != 200 {
            return false, "", errors.New("Validator responded with error code")
        }

        validatorResponse := ValidatorResponse{}
        err = json.NewDecoder(resp.Body).Decode(&validatorResponse)
        if err != nil {
            return false, "", err
        }

        if !validatorResponse.Valid {
            return false, strings.Join(validatorResponse.Errors, "\n"), nil
        }

        return true, validatorResponse.Result, nil
}

type ValidatorOptions struct {
        UriPrefix    string `json:"uri_prefix"`
        Version      string `json:"version"`
        TestEquality bool   `json:"test_equality"`
}

type ValidatorRequest struct {
        Options    ValidatorOptions `json:"options"`
        ReturnFile bool             `json:"return_file"`
        MainFile   string           `json:"main_file"`
}

type ValidatorResponse struct {
        Valid         bool     `json:"valid"`
        CheckEquality bool     `json:"check_equality"`
        Equality      bool     `json:"equality"`
        Errors        []string `json:"errors"`
        OutputFile    string   `json:"output_file"`
        Result        string   `json:"result"`
}

func JSONMarshal(t interface{}) ([]byte, error) {
        buffer := &bytes.Buffer{}
        encoder := json.NewEncoder(buffer)
        encoder.SetEscapeHTML(false)
        err := encoder.Encode(t)
        return buffer.Bytes(), err
}

