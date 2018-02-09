package main

import (
	"bufio"
	"fmt"
	"github.com/boltdb/bolt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

const pbfUrl = "http://download.geofabrik.de/europe/great-britain/wales-latest.osm.pbf"

func main() {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "127.0.0.1"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "osm"
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "postgres"
	}
	dbPass := os.Getenv("DB_PASS")

	db, err := setupDb()
	if err != nil {
		log.Fatalf("Can't setup db: %v\n", err)
	}
	defer db.Close()

	ticker := time.NewTicker(time.Minute)
	for range ticker.C {
		fmt.Println("Checking for updates")
		md5, err := getMD5()
		if err != nil {
			log.Printf("Can't get MD5 from server: %v\n", err)
			continue
		}

		var oldMd5 string
		err = db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("status"))
			oldMd5 = string(b.Get([]byte("md5")))
			return nil
		})
		if err != nil {
			log.Printf("Can't old MD5 from db: %v\n", err)
			continue
		}

		if oldMd5 == md5 {
			log.Println("No change")
		} else {
			log.Println("Chahge occured, updating")

			err := download()
			if err != nil {
				log.Printf("Can't download pbf: %v", err)
			}

			err = update(dbHost, dbName, dbUser, dbPass)
			if err != nil {
				log.Printf("Can't update database: %v", err)
			}

			err = db.Update(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte("status"))
				err = b.Put([]byte("md5"), []byte(md5))
				return err
			})
			if err != nil {
				log.Printf("Can't update database: %v\n", err)
			}
		}
	}
}

func setupDb() (*bolt.DB, error) {
	db, err := bolt.Open("data.db", 0600, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("status"))
		return err
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func getMD5() (string, error) {
	md5Resp, err := http.Get(pbfUrl + ".md5")
	if err != nil {
		return "", err
	}
	defer md5Resp.Body.Close()
	md5Bytes, err := ioutil.ReadAll(md5Resp.Body)
	md5 := strings.Split(string(md5Bytes), " ")
	return md5[0], nil
}

type PassThru struct {
	io.Reader
	total   int64
	Len     int64
	lastOut float64
}

func (pt *PassThru) Read(p []byte) (int, error) {
	n, err := pt.Reader.Read(p)
	pt.total += int64(n)

	if err == nil {
		percent := (float64(pt.total) / float64(pt.Len)) * 100
		if percent-pt.lastOut > 1 {
			fmt.Printf("\rRead %d%%", int64(percent))
			pt.lastOut = percent
		}
	}

	return n, err
}

func download() error {
	log.Println("Started download")
	out, err := os.Create("latest.osm.pbf")
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(pbfUrl)
	if err != nil {
		return err
	}
	respWrapped := &PassThru{Reader: resp.Body, Len: resp.ContentLength}
	defer resp.Body.Close()

	_, err = io.Copy(out, respWrapped)
	log.Println("Finished download")
	return err
}

func update(dbHost string, dbName string, dbUser string, dbPass string) error {
	cmd := exec.Command("osm2pgsql", "-d", dbName, "-H", dbHost, "-U", dbUser, "--create", "--slim", "-G",
		"--hstore", "-C", "16000", "latest.osm.pbf")

	env := os.Environ()
	env = append(env, fmt.Sprintf("PGPASSWORD=%s", dbPass))
	cmd.Env = env

	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			log.Println(scanner.Text())
		}
	}()
	cmdErrReader, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	errScanner := bufio.NewScanner(cmdErrReader)
	go func() {
		for errScanner.Scan() {
			log.Println(errScanner.Text())
		}
	}()

	log.Println("Starting update")
	err = cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	if err != nil {
		return err
	}
	log.Println("Update done")
	return nil
}
