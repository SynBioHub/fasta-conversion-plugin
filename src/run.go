package main

import (
        "bytes"
        "encoding/json"
        "errors"
        "fmt"
        "io"
        "io/ioutil"
        "net/http"
        "strings"
        "archive/zip"
        "strconv"
)

func Run(w http.ResponseWriter, r *http.Request) {
    // Unpack the request from SynBioHub
    request := SubmitRequest{}
    err := json.NewDecoder(r.Body).Decode(&request)
    if err != nil {
        fmt.Fprintln(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Open a temporary file to use as the ZIP
    tempZip, err := ioutil.TempFile("", "submit")
    if err != nil {
        fmt.Fprintln(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer tempZip.Close()

    // Create the Writer which will write zipped bytes to tempZip
    zipWriter := zip.NewWriter(tempZip)
    defer zipWriter.Close()

    // Begin setting up the manifest of the response
    response := ResponseManifest{}

    for _, file := range request.Manifest {
        // Get the FASTA file we are going to convert 
        resp, err := http.Get(file.URL)
        if err != nil {
            fmt.Fprintln(w, err.Error(), http.StatusInternalServerError)
            return
        }

        convertedFilename := file.Filename + ".converted"

        // Read the file from a buffer of bytes into a string
        buf := new(bytes.Buffer)
        buf.ReadFrom(resp.Body)
        fasta := buf.String()

        // Convert the FASTA file to SBOL using the validator
        valid, result, err := ConvertFastaToGenbank(fasta)
        if err != nil {
            fmt.Fprintln(w, err.Error(), http.StatusInternalServerError)
            return
        } else if !valid {
            fmt.Fprintln(w, result, http.StatusUnprocessableEntity)
            return
        }

        // Prepare to write the converted SBOL to the ZIP
        fileWriter, err := zipWriter.Create(convertedFilename)
        if err != nil {
            fmt.Fprintln(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // Write the SBOL to the ZIP file
        _, err = io.WriteString(fileWriter, result)
        if err != nil {
            fmt.Fprintln(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // Update our returned manifest with the file info
        resultInfo := ProcessedFileInfo{convertedFilename, []string{file.Filename}}
        response.Results = append(response.Results, resultInfo)
    }

    // Prepare to write the manifest to the zip file
    fileWriter, err := zipWriter.Create("manifest.json")
    if err != nil {
        fmt.Fprintln(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Unpack the manifest object into a JSON string
    responseJson, err := JSONMarshal(response)
    if err != nil {
        fmt.Fprintln(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Write the JSON to the zip file
    _, err = fileWriter.Write(responseJson)
    if err != nil {
        fmt.Fprintln(w, err.Error(), http.StatusInternalServerError)
        return
    }

    zipWriter.Close()
    
    // Read the completed zip file back into memory
    zipData, err := ioutil.ReadFile(tempZip.Name())
    if err != nil {
        fmt.Fprintln(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Stat the zip file so we know its size
    fstat, err := tempZip.Stat()
    if err != nil {
        fmt.Fprintln(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Set the necessary response headers
    // This one gives the filename of the response file as result.zip
    w.Header().Set("Content-Disposition", "attachment; filename=result.zip")
    // Set the filetype of the response as a zip
    w.Header().Set("Content-Type", "application/zip")
    // Set the content length of the response
    w.Header().Set("Content-Length", strconv.FormatInt(fstat.Size(), 10))
    // Write the zip data to the response
    _, err = w.Write(zipData)
    if err != nil {
        fmt.Fprintln(w, err.Error(), http.StatusInternalServerError)
        return
    }

    return
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
