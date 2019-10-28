package main

import (
        "bytes"
        "encoding/json"
)

type ProcessedFileInfo struct {
    Filename string `json:"filename"`
    Sources  []string `json:"sources"`
}

type FileInfo struct {
    URL string `json:"url"`
    Filename string `json:"filename"`
    EDAMType string `json:"edam"`
}

type SubmitRequest struct {
    Manifest []FileInfo `json:"manifest"`
}

type ResponseManifest struct {
    Results []ProcessedFileInfo `json:"results"`
}

type Need int
const (
    WillNotUse Need = 0
    WillRead   Need = 1
    WillHandle Need = 2
)

type EvaluateFileResponse struct {
    Filename    string `json:"filename"`
    Requirement Need   `json:"requirement"`
}

type EvaluateResponse struct {
    Files []EvaluateFileResponse `json:"files"`
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

