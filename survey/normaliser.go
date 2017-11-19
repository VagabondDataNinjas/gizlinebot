package survey

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/VagabondDataNinjas/gizlinebot/storage"
)

type Normaliser struct {
	Storage *storage.Sql
}

func NewNormaliser(storage *storage.Sql) *Normaliser {
	return &Normaliser{
		Storage: storage,
	}
}

func (n *Normaliser) Start(errc chan error) {
	for c := time.Tick(30 * time.Second); ; <-c {
		lastNormalisedAnswerId, err := n.Storage.GetLastNormalisedPriceAnswerId()
		if err != nil {
			log.Errorf("Error getting last normalised answerd id: %s", err)
			continue
		}

		_, err = n.Storage.NormalisePrices(lastNormalisedAnswerId)
		if err != nil {
			log.Errorf("Error normalising prices: %s", err)
			continue
		}
	}
}
