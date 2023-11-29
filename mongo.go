package ProfileBackend

import (
	"context"

	pkb "github.com/FarhanRizkiM/pasetobackend"
	"github.com/aiteung/atdb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func InsertDataParkiran(MongoConn *mongo.Database, colname string, emp Parkiran) (InsertedID interface{}) {
	req := new(Parkiran)
	req.ParkiranId = emp.ParkiranId
	req.Nama = emp.Nama
	req.NPM = emp.NPM
	req.Jurusan = emp.Jurusan
	req.NamaKendaraan = emp.NamaKendaraan
	req.NomorKendaraan = emp.NomorKendaraan
	req.JenisKendaraan = emp.JenisKendaraan
	req.Akun = emp.Akun
	return pkb.InsertOneDoc(MongoConn, colname, req)
}

func GetAllParkiranData(Mongoconn *mongo.Database, colname string) []Parkiran {
	data := atdb.GetAllDoc[[]Parkiran](Mongoconn, colname)
	return data
}

func DeleteUser(Mongoconn *mongo.Database, colname, username string) (deleted interface{}, err error) {
	filter := bson.M{"username": username}
	data := atdb.DeleteOneDoc(Mongoconn, colname, filter)
	return data, err
}

func UpdateParkiran(Mongoconn *mongo.Database, ctx context.Context, emp Parkiran) (UpdateId interface{}, err error) {
	filter := bson.D{{"parkiranid", emp.ParkiranId}}
	res, err := Mongoconn.Collection("parkiran").ReplaceOne(ctx, filter, emp)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func UpdatePassword(mongoconn *mongo.Database, user pkb.User) (Updatedid interface{}) {
	filter := bson.D{{"username", user.Username}}
	pass, _ := pkb.HashPass(user.Password)
	update := bson.D{{"$Set", bson.D{
		{"password", pass},
	}}}
	res, err := mongoconn.Collection("user").UpdateOne(context.Background(), filter, update)
	if err != nil {
		return "gagal update data"
	}
	return res
}

func DeleteParkiranData(mongoconn *mongo.Database, colname, PakId string) (deletedid interface{}, err error) {
	filter := bson.M{"parkiranid": PakId}
	data := atdb.DeleteOneDoc(mongoconn, colname, filter)
	return data, err
}

func GetOneParkiranData(mongoconn *mongo.Database, colname, Pakid string) (dest Parkiran) {
	filter := bson.M{"parkiranid": Pakid}
	dest = atdb.GetOneDoc[Parkiran](mongoconn, colname, filter)
	return
}
