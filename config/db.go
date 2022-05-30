package config

import "time"

const DBUri = "mongodb://localhost:27017"
const DBName = "bscrapdb"
const ResultsCol = "bscrap_res"
const RawDataCol = "bscrap_source"
const RecordExpirationTime = time.Second * 60
