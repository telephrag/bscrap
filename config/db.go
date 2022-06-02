package config

import "time"

const DBUri = "mongodb://localhost:27017"
const DBName = "bscrapdb"
const ResultsCol = "bscrap_res"
const SourceDataCollection = "bscrap_source"
const RecordExpirationTime = time.Hour * 24
