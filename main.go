package main

import (
	"embed"
//	"fmt"
	. "github.com/yoonjin67/tk9.0"
	"log"
	"os"
	"os/exec"
	"strings"
)

//go:embed tcl-awthemes
var awthemes embed.FS

type PkgInfo struct {
	id     string
	text   string
	values []string
}

func SearchPkg(pkglist []string) []PkgInfo {
	foundResult := make([]PkgInfo, 0, 4096)
	if len(pkglist) > 0 {
		for _, pkg := range pkglist {
			log.Println("Package is ", pkg)
			cmd := exec.Command("dnf", "search", pkg)
			cmd.Env = os.Environ()
			cmd.Env = append(cmd.Env, "LANG=C")
			out, err := cmd.CombinedOutput()
			log.Println(string(out))
			if err != nil {
				print(out)
				panic(err)
			}
			results := strings.Split(string(out), "\n")
			for _, line := range results {
				pkginfo := strings.Split(line, ":")
				for idx, v := range pkginfo {
					pkginfo[idx] = strings.TrimSpace(v)
				}
				if len(pkginfo) < 2 {
					continue
				}
				if pkginfo[0] == "Matched Fields" {
					continue
				}
				info := strings.Split(pkginfo[0], ".")
				for idx, v := range info {
					info[idx] = strings.TrimSpace(v)
				}
				if len(info) == 2 && len(pkginfo) == 2 {
					if info[0] == "" || info[1] == "" || pkginfo[1] == "" {
						continue
					}
					values := []string{
						info[0],
						info[1],
						Quote(pkginfo[1]),
					}
					itm := PkgInfo{
						id:     info[0],
						text:   Quote(info[0]),
						values: values,
					}
					log.Println("Package Name: ", info[0])
					log.Println("Package Info: ", values[2])
					foundResult = append(foundResult, itm)
				}
			}
		}
	}
	return foundResult
}

func main() {
	DefaultTheme("awbreeze", "tcl-awthemes")
	tr := TreeView()
	tr.SelectMode("extended")
	tr.Configure("-show", "headings")

	tr.Columns([]string{"name", "arch", "info"})

	tr.Heading("#0", "-text", Quote("PKGMAN"), "-anchor", "w")
	tr.Heading("name", "-text", Quote("Name"), "-anchor", "w")
	tr.Heading("arch", "-text", Quote("Arch"), "-anchor", "w")
	tr.Heading("info", "-text", Quote("Info"), "-anchor", "w")

	tr.Column("#0", "-width 1600 -anchor w")
	tr.Column("name", "-width 400 -anchor w")
	tr.Column("arch", "-width 400 -anchor w")
	tr.Column("info", "-width 800 -anchor w")
	e := TEntry()
	res := make([]PkgInfo, 0)
	SearchBtn := Button(Txt("Search!"), Command(func() {
		search := e.Get()
		log.Println(search)
		if len(search) != 0 {
			if len(res) > 0 {
				delList := make([]string, 0, 4096)
				for _, v := range res {
					delList = append(delList, v.id)
				}
				tr.Delete(delList)
			}
			searchlist := strings.Split(search, ",")
			for i, itm := range searchlist {
				searchlist[i] = strings.TrimSpace(itm)
				log.Println("SearchList Item", itm)
			}
			res = SearchPkg(searchlist)
			for i, r := range res {
				log.Println("Inserting ", r.id)
				for _, val := range r.values {
					log.Println("Values: ", val)
				}
				tr.Insert(r.id, r.text, 0, i, r.values)
			}
		}
	}))
	Pack(
		Label(Txt("DNF-Tk Package Manager")),
		e,
		SearchBtn,
		tr,
	)
	App.Wait()
}

