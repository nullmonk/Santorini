package santorini

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"go.etcd.io/bbolt"
	"golang.org/x/crypto/bcrypt"
)

const MB = 1000000

// A Script object that players upload to contain bots
type Script struct {
	Name     string `json:"name"`     // Name of the script
	Contents string `json:"contents"` // The actual script contents
	Md5      string `json:"md5"`      // Hash of the script
	Password string `json:"password"` // Password hash needed to edit the script
}

type Database struct {
	db *bbolt.DB
}

func NewDatabase(path string) (*Database, error) {
	db, err := bbolt.Open(path, 0644, &bbolt.Options{
		Timeout: time.Second * 5,
	})
	if err != nil {
		return nil, err
	}
	return &Database{
		db: db,
	}, nil
}

func (d *Database) SaveScript(password string, script *Script) error {
	// Validation
	if len(script.Contents) > 50*MB {
		return fmt.Errorf("script contents exceed limit (50MB)")
	}
	if len(script.Name) > 256 {
		return fmt.Errorf("script name exceeds limit (256)")
	}
	hash := md5.Sum([]byte(script.Contents))
	script.Md5 = hex.EncodeToString(hash[:])

	if s, _ := d.LoadScript(script.Name); s != nil {
		if script.Md5 == s.Md5 {
			return nil // Nothing needs to be done, its the same
		}
		err := bcrypt.CompareHashAndPassword([]byte(script.Password), []byte(password))
		if err != nil {
			return fmt.Errorf("invalid password")
		}
		script.Password = s.Password
	} else {
		pass, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
		script.Password = string(pass)
	}

	// Encode to json
	raw, err := json.Marshal(script)
	if err != nil {
		return fmt.Errorf("error converting to json: %s", err)
	}
	return d.db.Update(func(tx *bbolt.Tx) error {
		b, _ := tx.CreateBucketIfNotExists([]byte("scripts"))
		return b.Put([]byte(script.Name), raw)
	})
}

func (d *Database) LoadScript(name string) (*Script, error) {
	s := &Script{}
	if err := d.db.View(func(tx *bbolt.Tx) error {
		b, _ := tx.CreateBucketIfNotExists([]byte("scripts"))
		res := b.Get([]byte(name))
		if res == nil {
			s = nil
			return nil
		}
		return json.Unmarshal(res, &s)
	}); err != nil {
		return nil, err
	}
	return s, nil
}

func (d *Database) Scripts() []string {
	scripts := make([]string, 0, 10)
	d.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("scripts"))
		if b == nil {
			return nil
		}
		b.ForEach(func(k, v []byte) error {
			scripts = append(scripts, string(k))
			return nil
		})
		return nil
	})
	return scripts
}
