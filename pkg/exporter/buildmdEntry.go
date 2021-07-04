package exporter

import (
	"fmt"
)

// func buildMDEntry(e parse.Entry) string {
// 	link, source := mainLink(e.Site, e.Source)
// 	demo := linkSyntax(e.Demo, "Demo")
// 	clients := clientLinks(e.Clients)
// 	linkString := links([3]string{demo, source, clients})

// 	return fmt.Sprintf("- [%s](%s)] %s- %s%s `%s` `%s`", e.Name, link, pDep(e.Pdep), e.Descrip, linkString, e.License, e.Lang)
// }

func clientLinks(clientLinks []string) (clients string) {
	for i, c := range clientLinks {
		if i > 0 {
			clients += fmt.Sprintf(", %s", linkSyntax(c, "Client"))
		}
		clients += linkSyntax(c, "Client")
	}
	return
}

func mainLink(site string, source string) (string, string) {
	if site == "" {
		return source, ""
	}
	return site, linkSyntax(source, "Source Code")
}

func pDep(p bool) string {
	if p {
		return "`⚠` "
	}
	return ""
}

func links(links [3]string) (res string) {
	i := 0
	for i, l := range links {
		if l != "" {
			if i == 0 {
				res = fmt.Sprintf(" (%s", l)
			} else {
				res += fmt.Sprintf(", %s", l)
			}
		}
	}
	if i > 0 {
		res += ")"
	}
	return
}

func linkSyntax(link string, lType string) string {
	if link != "" {
		return fmt.Sprintf("[%s](%s)", lType, link)
	}
	return ""
}
