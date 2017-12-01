package survey

import (
	"time"

	"github.com/VagabondDataNinjas/gizlinebot/storage"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
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
		n.normalisePrices(errc)
		n.normaliseIslands(errc)
	}
}

func (n *Normaliser) normalisePrices(errc chan error) {
	// log.Info("Normalising prices")
	lastNormalisedAnswerId, err := n.Storage.GetLastNormalisedPriceId()
	if err != nil {
		errc <- errors.Wrap(err, "Error getting last normalised price id")
		return
	}

	_, err = n.Storage.NormalisePrices(lastNormalisedAnswerId)
	if err != nil {
		log.Errorf("Got error while normalising prices %s", err)
		errc <- errors.Wrap(err, "Error normalising prices")
		return
	}
}

func (n *Normaliser) normaliseIslands(errc chan error) {
	// log.Info("Normalising Islands")
	lastNormalisedAnswerId, err := n.Storage.GetLastNormalisedIslandId()
	if err != nil {
		errc <- errors.Wrap(err, "Error getting last normalised island id")
		return
	}

	_, err = n.Storage.NormaliseIslands(lastNormalisedAnswerId)
	if err != nil {
		errc <- errors.Wrap(err, "Error normalising islands")
		return
	}
}
