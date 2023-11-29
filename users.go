package ProfileBackend

import (
	"context"
	"encoding/json"
	pkb "github.com/FarhanRizkiM/pasetobackend"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"os"
)

// reg User
func Register(Mongoenv, dbname string, r *http.Request) string {
	resp := new(pkb.Credential)
	userdata := new(pkb.User)
	resp.Status = false
	conn := pkb.MongoCreateConnection(Mongoenv, dbname)
	err := json.NewDecoder(r.Body).Decode(&userdata)
	if err != nil {
		resp.Message = "error parsing application/json: " + err.Error()
	} else {
		resp.Status = true
		hash, err := pkb.HashPass(userdata.Password)
		if err != nil {
			resp.Message = "Gagal Hash Password" + err.Error()
		}
		pkb.InsertUserdata(conn, userdata.Username, userdata.Role, hash)
		resp.Message = "Berhasil Input data"
	}
	response := pkb.ReturnStringStruct(resp)
	return response

}

// log User
func Login(Privatekey, MongoEnv, dbname, Colname string, r *http.Request) string {
	var resp pkb.Credential
	mconn := pkb.MongoCreateConnection(MongoEnv, dbname)
	var datauser pkb.User
	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		resp.Message = "error parsing application/json: " + err.Error()
	} else {
		if pkb.PasswordValidator(mconn, Colname, datauser) {
			datarole := pkb.GetOneUser(mconn, "user", pkb.User{Username: datauser.Username})
			tokenstring, err := pkb.EncodeWithRole(datarole.Role, datauser.Username, os.Getenv(Privatekey))
			if err != nil {
				resp.Message = "Gagal Encode Token : " + err.Error()
			} else {
				resp.Status = true
				resp.Message = "Selamat Datang"
				resp.Token = tokenstring
			}
		} else {
			resp.Message = "Password Salah"
		}
	}
	return pkb.ReturnStringStruct(resp)
}

// Get Data User
func GetDataUserForAdmin(PublicKey, MongoEnv, dbname, colname string, r *http.Request) string {
	req := new(pkb.ResponseDataUser)
	conn := pkb.MongoCreateConnection(MongoEnv, dbname)
	tokenlogin := r.Header.Get("Login")
	if tokenlogin == "" {
		req.Status = false
		req.Message = "Header Login Not Found"
	} else {
		cekadmin := IsAdmin(tokenlogin, PublicKey)
		if cekadmin != true {
			req.Status = false
			req.Message = "IHHH Kamu bukan admin"
		}
		checktoken, err := pkb.DecodeGetUser(os.Getenv(PublicKey), tokenlogin)
		if err != nil {
			req.Status = false
			req.Message = "tidak ada data username : " + tokenlogin
		}
		compared := pkb.CompareUsername(conn, colname, checktoken)
		if compared != true {
			req.Status = false
			req.Message = "Data User tidak ada"
		} else {
			datauser := pkb.GetAllUser(conn, colname)
			req.Status = true
			req.Message = "data User berhasil diambil"
			req.Data = datauser
		}
	}
	return pkb.ReturnStringStruct(req)
}

// Reset Password
func ResetPassword(MongoEnv, publickey, dbname, colname string, r *http.Request) string {
	resp := new(Cred)
	req := new(pkb.User)
	conn := pkb.MongoCreateConnection(MongoEnv, dbname)
	tokenlogin := r.Header.Get("Login")
	if tokenlogin == "" {
		resp.Status = fiber.StatusBadRequest
		resp.Message = "Token login tidak ada"
	} else {
		checkadmin := IsAdmin(tokenlogin, os.Getenv(publickey))
		if !checkadmin {
			resp.Status = fiber.StatusInternalServerError
			resp.Message = "kamu bukan admin"
		} else {
			UpdatePassword(conn, pkb.User{
				Username: req.Username,
				Password: req.Password,
			})
			resp.Status = fiber.StatusOK
			resp.Message = "Berhasil reset password"
		}
	}
	return pkb.ReturnStringStruct(resp)
}

// Delete User
func DeleteUserforAdmin(Mongoenv, publickey, dbname, colname string, r *http.Request) string {
	resp := new(Cred)
	req := new(ReqUsers)
	conn := pkb.MongoCreateConnection(Mongoenv, dbname)
	tokenlogin := r.Header.Get("Login")
	if tokenlogin == "" {
		resp.Status = fiber.StatusBadRequest
		resp.Message = "Token login tidak ada"
	} else {
		checkadmin := IsAdmin(tokenlogin, os.Getenv(publickey))
		if !checkadmin {
			resp.Status = fiber.StatusInternalServerError
			resp.Message = "kamu bukan admin"
		} else {
			_, err := DeleteUser(conn, colname, req.Username)
			if err != nil {
				resp.Status = fiber.StatusBadRequest
				resp.Message = "gagal hapus data"
			}
			resp.Status = fiber.StatusOK
			resp.Message = "data berhasil dihapus"
		}
	}
	return pkb.ReturnStringStruct(resp)
}

// Insert data
func InsertParkiran(MongoEnv, dbname, colname, publickey string, r *http.Request) string {
	resp := new(pkb.Credential)
	req := new(Parkiran)
	conn := pkb.MongoCreateConnection(MongoEnv, dbname)
	tokenlogin := r.Header.Get("Login")
	if tokenlogin == "" {
		resp.Status = false
		resp.Message = "Header Login Not Found"
	} else {
		checkadmin := IsAdmin(tokenlogin, os.Getenv(publickey))
		if !checkadmin {
			checkHR := IsPK(tokenlogin, os.Getenv(publickey))
			if !checkHR {
				resp.Status = false
				resp.Message = "Anda tidak bisa Insert data karena bukan HR atau admin"
			}
		} else {
			err := json.NewDecoder(r.Body).Decode(&req)
			if err != nil {
				resp.Message = "error parsing application/json: " + err.Error()
			} else {
				pass, err := pkb.HashPass(req.Akun.Password)
				if err != nil {
					resp.Status = false
					resp.Message = "Gagal Hash Code"
				}
				InsertDataParkiran(conn, colname, Parkiran{
					ParkiranId: req.ParkiranId,
					Nama:       req.Nama,
					NPM:      req.NPM,
					Jurusan:      req.Jurusan,
					NamaKendaraan:      req.NamaKendaraan,
					NomorKendaraan:      req.NomorKendaraan,
					JenisKendaraan:      req.JenisKendaraan,
					Akun: pkb.User{
						Username: req.Akun.Username,
						Password: pass,
						Role:     req.Akun.Role,
					},
				})
				pkb.InsertUserdata(conn, req.Akun.Username, req.Akun.Role, pass)
				resp.Status = true
				resp.Message = "Berhasil Insert data"
			}
		}
	}
	return pkb.ReturnStringStruct(resp)
}

// Update data
func UpdateDataParkiran(MongoEnv, dbname, publickey string, r *http.Request) string {
	req := new(pkb.Credential)
	resp := new(Parkiran)
	tokenlogin := r.Header.Get("Login")
	if tokenlogin == "" {
		req.Status = false
		req.Message = "Header Login Not Found"
	} else {
		err := json.NewDecoder(r.Body).Decode(&resp)
		if err != nil {
			req.Message = "error parsing application/json: " + err.Error()
		} else {
			checkadmin := IsAdmin(tokenlogin, os.Getenv(publickey))
			if !checkadmin {
				checkHR := IsPK(tokenlogin, os.Getenv(publickey))
				if !checkHR {
					req.Status = false
					req.Message = "Anda tidak bisa Insert data karena bukan HR atau admin"
				}
			} else {
				conn := pkb.MongoCreateConnection(MongoEnv, dbname)
				UpdateParkiran(conn, context.Background(), Parkiran{
					ParkiranId: resp.ParkiranId,
					Nama:       resp.Nama,
					NPM:      resp.NPM,
					Jurusan:      resp.Jurusan,
					NamaKendaraan:      resp.NamaKendaraan,
					NomorKendaraan:      resp.NomorKendaraan,
					JenisKendaraan:      resp.JenisKendaraan,
					Akun: pkb.User{
						Username: resp.Akun.Username,
						Password: resp.Akun.Password,
						Role:     resp.Akun.Role,
					},
				})
				req.Status = true
				req.Message = "Berhasil Update data"
			}
		}
	}
	return pkb.ReturnStringStruct(req)
}

// Get One
func GetOneEmployee(PublicKey, MongoEnv, dbname, colname string, r *http.Request) string {
	req := new(ResponseParkiran)
	resp := new(RequestParkiran)
	conn := pkb.MongoCreateConnection(MongoEnv, dbname)
	tokenlogin := r.Header.Get("Login")
	if tokenlogin == "" {
		req.Status = fiber.StatusBadRequest
		req.Message = "Header Login Not Found"
	} else {
		err := json.NewDecoder(r.Body).Decode(&resp)
		if err != nil {
			req.Message = "error parsing application/json: " + err.Error()
		} else {
			checkadmin := IsAdmin(tokenlogin, os.Getenv(PublicKey))
			if !checkadmin {
				checkHR := IsPK(tokenlogin, os.Getenv(PublicKey))
				if !checkHR {
					req.Status = fiber.StatusBadRequest
					req.Message = "Anda tidak bisa Get data karena bukan HR atau admin"
				}
			} else {
				datauser := GetOneParkiranData(conn, colname, resp.ParkiranId)
				req.Status = fiber.StatusOK
				req.Message = "data User berhasil diambil"
				req.Data = datauser
			}
		}
	}
	return pkb.ReturnStringStruct(req)
}

// GetAll
func GetAllParkiran(PublicKey, Mongoenv, dbname, colname string, r *http.Request) string {
	req := new(ResponseParkiranBanyak)
	conn := pkb.MongoCreateConnection(Mongoenv, dbname)
	tokenlogin := r.Header.Get("Login")
	if tokenlogin == "" {
		req.Status = fiber.StatusBadRequest
		req.Message = "Header Login Not Found"
	} else {
		checkadmin := IsAdmin(tokenlogin, os.Getenv(PublicKey))
		if !checkadmin {
			checkHR := IsPK(tokenlogin, os.Getenv(PublicKey))
			if !checkHR {
				req.Status = fiber.StatusBadRequest
				req.Message = "Anda tidak bisa Insert data karena bukan HR atau admin"
			}
		} else {
			datauser := GetAllParkiranData(conn, colname)
			req.Status = fiber.StatusOK
			req.Message = "data User berhasil diambil"
			req.Data = datauser
		}
	}
	return pkb.ReturnStringStruct(req)
}

// Delete Data
func DeleteParkiran(Mongoenv, publickey, dbname, colname string, r *http.Request) string {
	resp := new(Cred)
	req := new(RequestParkiran)
	conn := pkb.MongoCreateConnection(Mongoenv, dbname)
	tokenlogin := r.Header.Get("Login")
	if tokenlogin == "" {
		resp.Status = fiber.StatusBadRequest
		resp.Message = "Token login tidak ada"
	} else {
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			resp.Message = "error parsing application/json: " + err.Error()
		} else {
			checkadmin := IsAdmin(tokenlogin, os.Getenv(publickey))
			if !checkadmin {
				resp.Status = fiber.StatusInternalServerError
				resp.Message = "kamu bukan admin"
			} else {
				_, err := DeleteParkiranData(conn, colname, req.ParkiranId)
				if err != nil {
					resp.Status = fiber.StatusBadRequest
					resp.Message = "gagal hapus data"
				}
				resp.Status = fiber.StatusOK
				resp.Message = "data berhasil dihapus"
			}
		}
	}
	return pkb.ReturnStringStruct(resp)
}