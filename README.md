# SynBioHub Submit plugin draft
This is a simple FASTA -> SBOL submit plugin

## Submit plugin architecture
 1. SynBioHub sends a submit manifest to the `/evaluate` endpoint of the plugin
    - For each file in the submission, the manifest will have a `filename`, a `url` where the file can be accessed, and a `edam` type, which is a URL representing SynBioHub's best guess at the file type.
    - The plugin should respond with a list of needs for each file, with each need containing a `filename`, and a `requirement`, where the `requirement` is 0 if the plugin will ignore the file, 1 if the plugin will reference the file but not convert it to SBOL, and 2 if the plugin will convert the file to SBOL.
    - *Demo*: This plugin simply reads the first byte of the file, and if it is '>', it assumes that it is a FASTA file and says it will convert it to SBOL.
  2. SynBioHub will decide which plugins will handle which files, and send an updated manifest to the `/run` endpoint of each plugin. The updated manifest will only contain the files which the plugin said it needed.
  3. The plugin should convert its files into SBOL. It may create one or more SBOL files. It should also create a manifest file, which describes which input files were used to create which output files.
  4. The manifest and all resulting plugin files should be zipped together by the plugin, and sent as a response to the `/run` request.

### Example evaluate manifest
```json
{
    "manifest": [{
        "url": "https://mineshaft.zach.network/4ccc8279-001a-440d-acc4-547ec51610ef.fasta",
        "filename": "zach.fasta",
        "edam": "http://edamontology.org/format_2330"
    }, {
        "url": "https://mineshaft.zach.network/ee8fc566-e416-4903-aaa7-9fce84fcd539.fasta",
        "filename": "jet.fasta",
        "edam": "http://edamontology.org/format_2330"
    }, {
        "url": "https://mineshaft.zach.network/49bc7acc-419a-47a0-9086-b4011e82f55d.sbol",
        "filename": "zach.sbol",
        "edam": "http://edamontology.org/format_2330"
    }]
}
```

### Example evaluate response
```json
{
    "manifest": [{
        "filename": "zach.fasta",
        "requirement": 2
    }, {
        "filename": "jet.fasta",
        "requirement": 2
    }, {
        "filename": "zach.sbol",
        "requirement": 0
    }]
}
```

### Example run request
```json
{
    "manifest": [{
        "url": "https://mineshaft.zach.network/4ccc8279-001a-440d-acc4-547ec51610ef.fasta",
        "filename": "zach.fasta",
        "edam": "http://edamontology.org/format_2330"
    }, {
        "url": "https://mineshaft.zach.network/ee8fc566-e416-4903-aaa7-9fce84fcd539.fasta",
        "filename": "jet.fasta",
        "edam": "http://edamontology.org/format_2330"
    }]
}
```

### Example response manifest
```json
{
    "results": [{
        "filename": "zach.fasta.converted",
        "sources": ["zach.fasta"]
    }, {
        "filename": "jet.fasta.converted",
        "sources": ["jet.fasta"]
    }]
}
