package ProfileBackend

import pkb "github.com/FarhanRizkiM/pasetobackend"

type ResponseBack struct {
	Status  int      `json:"status"`
	Message string   `json:"message"`
	Data    []string `json:"data"`
}

type ResponseParkiran struct {
	Status  int      `json:"status"`
	Message string   `json:"message"`
	Data    Parkiran `json:"data"`
}

type ResponseParkiranBanyak struct {
	Status  int        `json:"status"`
	Message string     `json:"message"`
	Data    []Parkiran `json:"data"`
}

type Parkiran struct {
	ParkiranId     string   `json:"parkiranid" bson:"parkiranid,omitempty"`
	Nama           string   `json:"nama" bson:"nama,omitempty"`
	NPM            string   `json:"npm" bson:"npm,omitempty"`
	Jurusan        string   `json:"jurusan" bson:"jurusan,omitempty"`
	NamaKendaraan  string   `json:"namakendaraan" bson:"namakendaraan,omitempty"`
	NomorKendaraan string   `json:"nomorkendaraan" bson:"nomorkendaraan"`
	JenisKendaraan string   `json:"jeniskendaraan" bson:"jeniskendaraan"`
	Akun           pkb.User `json:"akun" bson:"akun,omitempty"`
}

type Updated struct {
	NamaKendaraan  string `json:"namakendaraan" bson:"namakendaraan"`
	NomorKendaraan string `json:"nomorkendaraan" bson:"nomorkendaraan"`
	JenisKendaraan string `json:"jeniskendaraan" bson:"jeniskendaraan"`
}

type Cred struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type ReqUsers struct {
	Username string `json:"username"`
}

type RequestParkiran struct {
	ParkiranId string `json:"parkiranid"`
}
