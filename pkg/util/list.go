package util

// List is the total struct
type List struct {
	Entries  []Entry `json:"Entries"`
	LangList []Langs `json:"Langs" yaml:"-"`
	//CatList  []Cats		`json:"Cats", yaml:""`
	TagList []Tags `json:"Tags" yaml:"-"`
}

// Entry is the structure of each entry
type Entry struct {
	ID      int      `json:"ID" yaml:"ID"`
	Name    string   `json:"N" yaml:"Name"`
	Descrip string   `json:"D" yaml:"Description,flow"`
	Source  string   `json:"Sr" yaml:"Source Code"`
	Demo    string   `json:"Dem,omitempty" yaml:"Demo,omitempty"`
	Clients []string `json:"CL,omitempty" yaml:"Clients,omitempty"`
	Site    string   `json:"Si,omitempty" yaml:"Website,omitempty"`
	License []string `json:"Li" yaml:"License"`
	Lang    []string `json:"La" yaml:"Languages"`
	//Cat     string   `json:"C,omitempty"`
	Tags    []string `json:"T" yaml:"Tags"`
	NonFree bool     `json:"NF,omitempty" yaml:"NonFree,omitempty"`
	Pdep    bool     `json:"P,omitempty" yaml:"ProprietaryDependency,omitempty"`
	Stars   int      `json:"stars,omitempty" yaml:"-"`
	Created string   `json:"create,omitempty" yaml:"-"`
	Updated string   `json:"update,omitempty" yaml:"-"`
	Errors  []string `json:"-" yaml:"Errors,omitempty"`
	Warns   []string `json:"-" yaml:"Warnings,omitempty"`
}