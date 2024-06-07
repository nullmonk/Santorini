package santorini

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"go.etcd.io/bbolt"
	"golang.org/x/crypto/bcrypt"
)

const MB = 1000000

var Storage *Database

func init() {
	pth := os.Getenv("SANTORENA_DB")
	if pth == "" {
		pth = "santorena.db"
	}
	var err error
	Storage, err = NewDatabase(pth)
	if err != nil {
		panic(err)
	}
}

// A Script object that players upload to contain bots
type Script struct {
	Name     string `json:"name"`     // Name of the script
	Contents string `json:"contents"` // The actual script contents
	Md5      string `json:"md5"`      // Hash of the script
	Password string `json:"password"` // Password hash needed to edit the script
}

type ScriptSanatized struct {
	Name     string `json:"name"`     // Name of the script
	Contents string `json:"contents"` // The actual script contents
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
	// Validation of sizes
	if len(script.Contents) > 50*MB {
		return fmt.Errorf("script contents exceed limit (50MB)")
	}
	if len(script.Name) > 256 {
		return fmt.Errorf("script name exceeds limit (256)")
	}

	// Calculate the hash, if it has not changed, no need to do anything
	hash := md5.Sum([]byte(script.Contents))
	script.Md5 = hex.EncodeToString(hash[:])
	s, _ := d.LoadScript(script.Name)
	if s != nil {
		if script.Md5 == s.Md5 {
			return nil // Nothing needs to be done, its the same
		}
		// Check the password
		err := bcrypt.CompareHashAndPassword([]byte(s.Password), []byte(password))
		if err != nil {
			return fmt.Errorf("invalid password")
		}
		script.Password = s.Password
	} else {
		// New script
		pass, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
		script.Password = string(pass)
		if len(script.Contents) < 15 {
			return fmt.Errorf("refusing to save empty script")
		}
	}

	// Encode to json
	raw, err := json.Marshal(script)
	if err != nil {
		return fmt.Errorf("error converting to json: %s", err)
	}
	return d.db.Update(func(tx *bbolt.Tx) error {
		b, _ := tx.CreateBucketIfNotExists([]byte("scripts"))
		// If the new script contents are empty, delete it
		if len(script.Contents) == 0 {
			return b.Delete([]byte(script.Name))
		}
		return b.Put([]byte(script.Name), raw)
	})
}

func (d *Database) LoadScript(name string) (*Script, error) {
	s := &Script{}
	if err := d.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("scripts"))
		if b == nil {
			s = nil
			return nil
		}
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

func (d *Database) Scripts() map[string]*ScriptSanatized {
	scripts := make(map[string]*ScriptSanatized, 10)
	d.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("scripts"))
		if b == nil {
			return nil
		}
		b.ForEach(func(k, v []byte) error {
			s := &ScriptSanatized{}
			json.Unmarshal(v, s)
			scripts[string(k)] = s
			return nil
		})
		return nil
	})
	return scripts
}

// LoadCache loads arbitrary JSON data into a script
func (d *Database) LoadCache(script string) []byte {
	var res []byte
	d.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("script_memory"))
		if b == nil {
			return nil
		}
		res = b.Get([]byte(script))
		return nil
	})
	return res
}

// SaveCache saves arbitrary JSON data from a script
func (d *Database) SaveCache(script string, memory []byte) error {
	if len(memory) > 50*MB {
		return fmt.Errorf("script memory exceeds limit (50MB)")
	}

	return d.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("script_memory"))
		if b == nil {
			return nil
		}
		return b.Put([]byte(script), memory)
	})
}

// SaveGame saves a game as JSON to the DB. If no ID is specified, generate one. If Game is nil, it is deleted
func (d *Database) SaveGame(id string, g *Game) (string, error) {
	if id == "" {
		id = uuid.New().String()
	}

	game, err := json.Marshal(g)
	if err != nil {
		return "", err
	}

	return id, d.db.Update(func(tx *bbolt.Tx) error {
		b, _ := tx.CreateBucketIfNotExists([]byte("games"))
		if b == nil {
			return nil
		}
		return b.Put([]byte(id), game)
	})
}

// Delete a game from the database
func (d *Database) DeleteGame(id string) error {
	return d.db.Update(func(tx *bbolt.Tx) error {
		b, _ := tx.CreateBucketIfNotExists([]byte("games"))
		if b == nil {
			return nil
		}
		return b.Delete([]byte(id))
	})
}

func (d *Database) LoadGame(id string) (*Game, error) {
	g := &Game{}
	err := d.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("games"))
		if b == nil {
			return nil
		}
		bytes := b.Get([]byte(id))
		return json.Unmarshal(bytes, &g)
	})
	if err != nil {
		// Delete this game save
		d.DeleteGame(id)
	}
	return g, err
}

// TODO Save stats on the different bots
