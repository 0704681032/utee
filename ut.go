package utee

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/smtp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

var (
	PlainMd5  = Md5Str("")
	PlainSha1 = Sha1Str("")
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func Md5Str(salt string) func(string) string {
	return func(s string) string {
		h := md5.New()
		io.WriteString(h, s)
		io.WriteString(h, salt)
		return hex.EncodeToString(h.Sum(nil))
	}
}

func Sha1Str(salt string) func(string) string {
	return func(s string) string {
		h := sha1.New()
		io.WriteString(h, s)
		io.WriteString(h, salt)
		return hex.EncodeToString(h.Sum(nil))
	}
}

func Sha256Str(salt string) func(string) string {
	return func(s string) string {
		h := sha256.New()
		io.WriteString(h, s)
		io.WriteString(h, salt)
		return hex.EncodeToString(h.Sum(nil))
	}
}

func HmacSha256(s string, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	io.WriteString(h, s)
	return hex.EncodeToString(h.Sum(nil))
}

func Chk(err error) {
	if err != nil {
		panic(err)
	}
}

// Truncate truncate string
func Truncate(s string, n int) string {
	if n <= 0 {
		return ""
	}

	length := utf8.RuneCountInString(s)
	if length <= n {
		return s
	}

	l := []rune{}
	for _, r := range s {
		l = append(l, r)
	}

	l = l[:(length - n)]
	return string(l)
}

// Tick unix timestamp in millisecond
func Tick(t ...time.Time) int64 {
	if len(t) == 0 {
		return time.Now().UnixNano() / 1e6
	}

	return t[0].UnixNano() / 1e6
}

// TickToTime convert unix timestamp in millisecond to time at current location
func TickToTime(tick int64) time.Time {
	return time.Unix(tick/1e3, (tick%1e3)*1e6)
}

func Md5(b []byte) []byte {
	h := md5.New()
	h.Write(b)
	return h.Sum(nil)
}

func DeleteMap(m map[string]interface{}, ks ...string) {
	for _, v := range ks {
		delete(m, v)
	}
}

func IsPemExpire(b []byte) (bool, error) {
	block, _ := pem.Decode(b)
	if block == nil {
		return false, errors.New("failed to parse certificate PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return false, err
	}
	return cert.NotAfter.Before(time.Now()), nil
}

func Shuffle(src []string) []string {
	dest := make([]string, len(src))
	perm := rand.Perm(len(src))
	for i, v := range perm {
		dest[v] = src[i]
	}
	return dest
}

func SendMail(user, password, host, to, subject, body, mailtype string) error {
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var content_type string
	if mailtype == "html" {
		content_type = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
	}

	msg := []byte("To: " + to + "\r\nFrom: " + user + "<" + user + ">\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	send_to := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, send_to, msg)
	return err
}

func ParseAddr(s string) (string, int, error) {
	a := strings.Split(s, ":")
	if len(a) != 2 {
		return "", 0, fmt.Errorf("bad url %s", s)
	}
	port, err := strconv.Atoi(a[1])
	return a[0], port, err
}

func Unique(data []interface{}) []interface{} {
	m := map[interface{}]interface{}{}

	for _, d := range data {
		m[d] = "0"
	}

	l := []interface{}{}
	for key := range m {
		l = append(l, key)
	}
	return l
}

func UniqueInt(data []int) []int {
	m := map[int]string{}

	for _, d := range data {
		m[d] = "0"
	}

	l := []int{}
	for key := range m {
		l = append(l, key)
	}
	return l
}

func UniqueStr(data []string) []string {
	m := map[string]string{}

	for _, d := range data {
		m[d] = "0"
	}

	l := []string{}
	for key := range m {
		l = append(l, key)
	}
	return l
}

// split a into several parts, no more than n
func SplitSlice(a []string, n int) [][]string {
	if len(a) < n || n == 1 {
		return [][]string{a}
	}

	result := make([][]string, n)
	for i, s := range a {
		idx := i % n
		result[idx] = append(result[idx], s)
	}
	return result
}

func IntToInf(src []int) []interface{} {
	result := []interface{}{}
	for _, v := range src {
		result = append(result, v)
	}
	return result
}

func StrToInf(src []string) []interface{} {
	result := []interface{}{}
	for _, v := range src {
		result = append(result, v)
	}
	return result
}

func PrintJson(any ...interface{}) {
	for _, obj := range any {
		b, err := json.Marshal(obj)
		fmt.Println(err, string(b))
	}
}
