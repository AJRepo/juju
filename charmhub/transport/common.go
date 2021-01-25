// Copyright 2020 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package transport

// The following contains all the common DTOs for a gathering information from
// a given store.

type Channel struct {
	Name       string   `json:"name"`
	Platform   Platform `json:"platform"`
	ReleasedAt string   `json:"released-at"`
	Risk       string   `json:"risk"`
	Track      string   `json:"track"`
}

type Platform struct {
	Architecture string `json:"architecture"`
	OS           string `json:"os"`
	Series       string `json:"series"`
}

// Download represents the download structure from CharmHub.
// Elements not used by juju but not used are: "hash-sha3-384"
// and "hash-sha-512"
type Download struct {
	HashSHA256 string `json:"hash-sha-256"`
	HashSHA384 string `json:"hash-sha-384"`
	Size       int    `json:"size"`
	URL        string `json:"url"`
}

type Entity struct {
	Categories  []Category        `json:"categories"`
	Charms      []Charm           `json:"contains-charms"`
	Description string            `json:"description"`
	License     string            `json:"license"`
	Publisher   map[string]string `json:"publisher"`
	Summary     string            `json:"summary"`
	UsedBy      []string          `json:"used-by"`
	StoreURL    string            `json:"store-url"`
}

type Category struct {
	Featured bool   `json:"featured"`
	Name     string `json:"name"`
}

type Charm struct {
	Name      string `json:"name"`
	PackageID string `json:"package-id"`
	StoreURL  string `json:"store-url"`
}
