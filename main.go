package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
)

func main() {
	var update bool
	flag.BoolVar(&update, "u", true, "Update cards before generating")

	var branch string
	flag.StringVar(&branch, "b", "develop", "Use specific branch (defaults to 'develop')")

	flag.Parse()

	log.Printf("Working with '%s' branch.", branch)

	if update {
		for countryCode, cardsURL := range cardsURLs {
			log.Printf("Downloading %s...", countryCode)
			url := fmt.Sprintf(cardsURL, branch)
			err := downloadFile(url, sourceFile(countryCode))
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	log.Println("Generating CSV...")

	var content []byte
	for countryCode, _ := range cardsURLs {
		buf, err := os.ReadFile(sourceFile(countryCode))
		if err != nil {
			log.Fatal(err)
		}

		var cards Cards
		err = json.Unmarshal(buf, &cards)
		if err != nil {
			log.Fatal(err)
		}

		sort.Sort(cards)

		for _, c := range cards {
			name := c.Name
			if strings.Contains(name, ",") {
				name = "\"" + name + "\""
			}

			s := fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s\n",
				uniqueID(c), c.ID, name, c.Pitch, foiling[c.Foiling], edition[c.Edition], c.Rarity, c.Cost, c.Power, c.Defense, c.ImageURL, strings.Join(c.Types, ","))
			content = append(content, []byte(s)...)
		}
	}

	err := os.WriteFile("cards.csv", content, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func sourceFile(countryCode string) string {
	return fmt.Sprintf("cards/%s.json", countryCode)
}

func downloadFile(url string, filepath string) error {
	output, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download: %s", response.Status)
	}

	_, err = io.Copy(output, response.Body)
	if err != nil {
		return err
	}

	return nil
}

// There is only fabled in different languages, but hence all keywords are translated it will mess with filters.
// Abandoning idea for implementing different languages for now.
// Revisit later, when dataset will have all the cards.
var cardsURLs = map[string]string{
	"en": "https://raw.githubusercontent.com/the-fab-cube/flesh-and-blood-cards/%s/json/english/card-flattened.json",
	//"fr": "https://raw.githubusercontent.com/the-fab-cube/flesh-and-blood-cards/develop/json/french/card-flattened.json",
	//"de": "https://raw.githubusercontent.com/the-fab-cube/flesh-and-blood-cards/develop/json/german/card-flattened.json",
	//"it": "https://raw.githubusercontent.com/the-fab-cube/flesh-and-blood-cards/develop/json/italian/card-flattened.json",
	//"es": "https://raw.githubusercontent.com/the-fab-cube/flesh-and-blood-cards/develop/json/spanish/card-flattened.json",
}

// Shorthand	Name
// A	Alpha
// F	First
// U	Unlimited
// N	No specified edition (used for promos, non-set releases, etc.)

var edition = map[string]string{
	"A": "Alpha",
	"F": "1st",
	"U": "Unl",
	"N": "",
}

// Shorthand	Name
// S	Standard
// R	Rainbow Foil
// C	Cold Foil
// G	Gold Cold Foil

var foiling = map[string]string{
	"S": "NF",
	"R": "RF",
	"C": "CF",
	"G": "GF",
}

func uniqueID(card Card) string {
	var uid = fmt.Sprintf("%s.%s", card.ID, foiling[card.Foiling])
	s := edition[card.Edition]
	if s == "" {
		return uid
	}
	return uid + "." + s
}

type Cards []Card

func (x Cards) Len() int           { return len(x) }
func (x Cards) Less(i, j int) bool { return x[i].Name < x[j].Name }
func (x Cards) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

type Card struct {
	UniqueID                 string        `json:"unique_id"`
	Name                     string        `json:"name"`
	Pitch                    string        `json:"pitch"`
	Cost                     string        `json:"cost"`
	Power                    string        `json:"power"`
	Defense                  string        `json:"defense"`
	Health                   string        `json:"health"`
	Intelligence             string        `json:"intelligence"`
	Types                    []string      `json:"types"`
	CardKeywords             []interface{} `json:"card_keywords"`
	AbilitiesAndEffects      []interface{} `json:"abilities_and_effects"`
	AbilityAndEffectKeywords []interface{} `json:"ability_and_effect_keywords"`
	GrantedKeywords          []interface{} `json:"granted_keywords"`
	RemovedKeywords          []interface{} `json:"removed_keywords"`
	InteractsWithKeywords    []interface{} `json:"interacts_with_keywords"`
	FunctionalText           string        `json:"functional_text"`
	FunctionalTextPlain      string        `json:"functional_text_plain"`
	TypeText                 string        `json:"type_text"`
	PlayedHorizontally       bool          `json:"played_horizontally"`
	BlitzLegal               bool          `json:"blitz_legal"`
	CcLegal                  bool          `json:"cc_legal"`
	CommonerLegal            bool          `json:"commoner_legal"`
	BlitzLivingLegend        bool          `json:"blitz_living_legend"`
	CcLivingLegend           bool          `json:"cc_living_legend"`
	BlitzBanned              bool          `json:"blitz_banned"`
	CcBanned                 bool          `json:"cc_banned"`
	CommonerBanned           bool          `json:"commoner_banned"`
	UpfBanned                bool          `json:"upf_banned"`
	BlitzSuspended           bool          `json:"blitz_suspended"`
	CcSuspended              bool          `json:"cc_suspended"`
	CommonerSuspended        bool          `json:"commoner_suspended"`
	PrintingUniqueID         string        `json:"printing_unique_id"`
	SetPrintingUniqueID      string        `json:"set_printing_unique_id"`
	ID                       string        `json:"id"`
	SetID                    string        `json:"set_id"`
	Edition                  string        `json:"edition"`
	Foiling                  string        `json:"foiling"`
	Rarity                   string        `json:"rarity"`
	Artist                   string        `json:"artist"`
	ArtVariation             interface{}   `json:"art_variation"`
	FlavorText               string        `json:"flavor_text"`
	FlavorTextPlain          string        `json:"flavor_text_plain"`
	ImageURL                 string        `json:"image_url"`
	TcgplayerProductID       string        `json:"tcgplayer_product_id"`
	TcgplayerURL             string        `json:"tcgplayer_url"`
}
