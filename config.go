package main

import (
    "log"
    "os"
    "strconv"

    "github.com/joho/godotenv"
)

type Config struct {
    SolanaRPCURL  string
    SerumAPIKey   string
    DBHost        string
    DBPort        string
    DBUser        string
    DBPassword    string
    DBName        string
    TargetWinRate float64
    MaxDrawdown   float64
}

func LoadConfig() Config {
    err := godotenv.Load()
    if err != nil {
        log.Println("No .env file found. Using environment variables.")
    }

    targetWinRate, err := strconv.ParseFloat(os.Getenv("TARGET_WIN_RATE"), 64)
    if err != nil {
        targetWinRate = 60.0 // default
    }

    maxDrawdown, err := strconv.ParseFloat(os.Getenv("MAX_DRAWDOWN"), 64)
    if err != nil {
        maxDrawdown = 20.0 // default
    }

    return Config{
        SolanaRPCURL:  os.Getenv("SOLANA_RPC_URL"),
        SerumAPIKey:   os.Getenv("SERUM_API_KEY"),
        DBHost:        os.Getenv("DB_HOST"),
        DBPort:        os.Getenv("DB_PORT"),
        DBUser:        os.Getenv("DB_USER"),
        DBPassword:    os.Getenv("DB_PASSWORD"),
        DBName:        os.Getenv("DB_NAME"),
        TargetWinRate: targetWinRate,
        MaxDrawdown:   maxDrawdown,
    }
}
