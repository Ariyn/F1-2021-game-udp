package main

const FormulaLoggerLogo = `                                                                                                                                                                                           
    __________  ____  __  _____  ____    ___       ___
   / ____/ __ \/ __ \/  |/  / / / / /   /   |     <  /
  / /_  / / / / /_/ / /|_/ / / / / /   / /| |     / / 
 / __/ / /_/ / _, _/ /  / / /_/ / /___/ ___ |    / /  
/_/ __ \______/______________________/_/  |_|   /_/   
   / /  / __ \/ ____/ ____/ ____/ __ \                
  / /  / / / / / __/ / __/ __/ / /_/ /                
 / /__/ /_/ / /_/ / /_/ / /___/ _, _/                 
/_____\____/\____/\____/_____/_/ |_|
`

type packetData struct {
	Buf       []byte
	Size      int
	Timestamp int64
}

type Driver struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	TeamName string `json:"teamName"`
	CarIndex int    `json:"CarIndex"`
	IsAi     bool   `json:"isAi"`
}
