package config

import "time"

const DBUri = "mongodb://localhost:27017"
const DBName = "bscrapdb"
const BScrapResCol = "bscrap_res"
const BScrapSourceCol = "bscrap_source"
const RecordLifeTime = time.Hour * 24
