package ipban

import (
	"bufio"
	"net/http"
	"os"
	"strings"
)

func IPBan(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		banListLoc := os.Getenv("BAN_LIST")
		banList := []string{}

		banListFile, err := os.Open(banListLoc)
		if err != nil {
			panic(err)
		}

		scanner := bufio.NewScanner(banListFile)
		for scanner.Scan() {
			banList = append(banList, scanner.Text())
		}

		if IPInList(r.RemoteAddr, banList) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			panic("IP Banned")
		}

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func IPInList(ip string, list []string) bool {
	// Ignore ports
	ipReal := strings.Split(ip, ":")[0]

	for _, v := range list {
		if v == ipReal {
			return true
		}
	}
	return false
}
