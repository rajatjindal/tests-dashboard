package junit

import "encoding/xml"

type Testsuites struct {
	XMLName   xml.Name  `xml:"testsuites"`
	Text      string    `xml:",chardata"`
	Testsuite TestSuite `xml:"testsuite"`
}

type TestSuite struct {
	Text      string     `xml:",chardata"`
	Name      string     `xml:"name,attr"`
	Package   string     `xml:"package,attr"`
	ID        string     `xml:"id,attr"`
	Errors    string     `xml:"errors,attr"`
	Failures  string     `xml:"failures,attr"`
	Tests     string     `xml:"tests,attr"`
	Skipped   string     `xml:"skipped,attr"`
	Testcases []TestCase `xml:"testcase"`
	SystemOut string     `xml:"system-out"`
	SystemErr string     `xml:"system-err"`
}

type TestCase struct {
	Text      string   `xml:",chardata"`
	Classname string   `xml:"classname,attr"`
	Name      string   `xml:"name,attr"`
	Time      string   `xml:"time,attr"`
	Failure   *Failure `xml:"failure"`
	SystemOut string   `xml:"system-out"`
}

type Failure struct {
	Text string `xml:",chardata"`
	Type string `xml:"type,attr"`
}
