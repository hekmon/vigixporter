package hubeau

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-querystring/query"
)

// ObservationsRequest allows to set request parameter for https://hubeau.eaufrance.fr/page/api-hydrometrie#/hydrometrie/observations
type ObservationsRequest struct {
	BBox       []float64       `url:"bbox,omitempty"`           // Rectangle d'emprise de l'objet demandé, emprise au format : min longitude, min latitude, max longitude, max latitude avec les coordonnées en WGS84 (EPSG:4326), le point doit être utilisé comme séparateur décimal, exemple : 1.6194,47.7965,2.1910,47.9988
	EntityCode []string        `url:"code_entite,omitempty"`    // Le site hydrométrique est identifié par un code unique au niveau national construit selon la règle suivante : [Code de la zone hydrographique sur 4 caractères sur laquelle est situé le site hydrologique] + [Numéro incrémental sans signification particulière sur 4 caractères]. Par exemple, J4310010 pour un site localisé sur la zone hydrologique J431.Code de la station hydrométrique. Chaque station est identifiée par un code unique pour un site donné. L'identifiant complet de la station hydrométrique est le code du site + le code de la station sur 2 caractères. Possibilité d'utiliser un pattern (ex: K*
	Cursor     string          `url:"cursor,omitempty"`         // Curseur de pagination
	StartDate  time.Time       `url:"date_debut_obs,omitempty"` // Date de début observation hydro (exprimée en TU). Ne peut pas être antérieure de plus d'1 mois par rapport à la date actuelle, les formats de date (ISO 8601) supportés : yyyy-MM-dd, yyyy-MM-dd'T'HH:mm:ss, yyyy-MM-dd'T'HH:mm:ssXXX, exemples : 2018-12-01, 2018-12-11T00:00:01, 2018-12-11T00:00:01Z
	EndDate    time.Time       `url:"date_fin_obs,omitempty"`   // Date de fin observation hydro (exprimée en TU), les formats de date (ISO 8601) supportés : yyyy-MM-dd, yyyy-MM-dd'T'HH:mm:ss, yyyy-MM-dd'T'HH:mm:ssXXX, exemples : 2018-12-01, 2018-12-11T00:00:01, 2018-12-11T00:00:01Z
	Distance   float64         `url:"distance,omitempty"`       // Rayon de recherche en kilomètre, le point doit être utilisé comme séparateur décimal, exemple : 30
	Type       ObservationType `url:"grandeur_hydro,omitempty"` // Which measurements to return
	Latitude   float64         `url:"latitude,omitempty"`       // Latitude du point en WGS84 pour la recherche par rayon, le point doit être utilisé comme séparateur décimal, exemple : 47.829
	Longitude  float64         `url:"longitude,omitempty"`      // Longitude du point en WGS84 pour la recherche par rayon, le point doit être utilisé comme séparateur décimal, exemple : 1.937
	Size       int             `url:"size,omitempty"`           // Number of value to return (default: 20, max: 20000)
	Sort       Sort            `url:"sort,omitempty"`           // Ordre de tri (asc ou desc) sur la date d'observation (si la valeur n'est pas renseignée, la valeur par défaut est desc)
	Timestep   int             `url:"timestep,omitempty"`       // Pas de temps fixe exprimé en minutes avec des limites (de 10 minutes à 60 minutes), ne fonctionne que pour la recherche sur un seul code entite, le service n'a pas de pagination
}

// ObservationType defines the type of observation
type ObservationType string

const (
	// ObservationTypeHeight represents the water height
	ObservationTypeHeight ObservationType = "H"
	// ObservationTypeSpeed represents the water speed
	ObservationTypeSpeed ObservationType = "Q"
	// ObservationTypeHeightAndSpeed represents both the water height and speed (used in query only)
	ObservationTypeHeightAndSpeed ObservationType = "H,Q"
)

// Sort represents the order of sorting
type Sort string

const (
	// SortAscending represents an ascending sort
	SortAscending Sort = "sort"
	// SortDescending represents a descending sort
	SortDescending Sort = "desc"
)

// GetObservations maps https://hubeau.eaufrance.fr/page/api-hydrometrie#/hydrometrie/observations
func (c *Controller) GetObservations(ctx context.Context, parameters ObservationsRequest) (response ObservationsResponse, err error) {
	urlValues, err := query.Values(parameters)
	if err != nil {
		err = fmt.Errorf("can't convert query parameters to url values: %w", err)
		return
	}
	if err = c.request(ctx, "GET", "hydrometrie/observations_tr", urlValues, &response); err != nil {
		err = fmt.Errorf("getting checks failed: %w", err)
		return
	}
	return
}

// ObservationsResponse represents the answer payload for GetObservations
type ObservationsResponse struct {
	Count      int           `json:"count"`       // Le nombre total de résultat
	First      string        `json:"first"`       // URL de la 1er page des résultats
	Prev       *string       `json:"prev"`        // Toujours null
	Next       string        `json:"next"`        // URL de la page suivante
	APIVersion string        `json:"api_version"` // Version de l'API (https://semver.org/)
	Data       []Observation `json:"data"`        // Les résultats de la requête sous forme de liste
}

// Observation represents a single observation
type Observation struct {
	SiteCode             string          `json:"code_site"`    // Le site hydrométrique est identifié par un code unique au niveau national construit selon la règle suivante : [Code de la zone hydrographique sur 4 caractères sur laquelle est situé le site hydrologique] + [Numéro incrémental sans signification particulière sur 4 caractères]. Par exemple, J4310010 pour un site localisé sur la zone hydrologique J431.
	StationCode          string          `json:"code_station"` // Code de la station hydrométrique. Chaque station est identifiée par un code unique pour un site donné. L'identifiant complet de la station hydrométrique est le code du site + le code de la station sur 2 caractères.
	Type                 ObservationType `json:"grandeur_hydro"`
	SerieStart           time.Time       `json:"date_debut_serie"`
	SerieEnd             time.Time       `json:"date_fin_serie"`
	SerieStatus          SerieStatus     `json:"statut_serie"` // Statut de la série : niveau de validité de la donnée. Valeurs possibles : 0, 4, 8, 12, 16 (sans validation, brut, corrigé, pré-validé, validé). Voir http://services.sandre.eaufrance.fr/References/1.3.0/jeuDonnees.php?recherche=510&function=getFicheNsa&v=3.1
	SystemCodeAltiSerie  int             `json:"code_systeme_alti_serie"`
	ObsDate              time.Time       `json:"date_obs"`
	ObsResultat          float64         `json:"resultat_obs"`
	ObsMethodCode        int             `json:"code_methode_obs"`
	ObsMethodLib         string          `json:"libelle_methode_obs"`
	ObsQualificationCode int             `json:"code_qualification_obs"`
	ObsQualificationLib  string          `json:"libelle_qualification_obs"`
	ObsHydroCont         bool            `json:"continuite_obs_hydro"`
	Longitude            float64         `json:"longitude"`
	Latitude             float64         `json:"latitude"`
}

// SerieStatus represents the status of a custom observation/serie.
type SerieStatus int

const (
	// SerieStatusNoValidation represent a serie which will not have a validation
	SerieStatusNoValidation SerieStatus = 0
	// SerieStatusRaw represents a serie which has not been altered
	SerieStatusRaw SerieStatus = 4
	// SerieStatusCorrected represents a serie which has received a correction
	SerieStatusCorrected SerieStatus = 8
	// SerieStatusPreValidated represents a serie which is pre validated
	SerieStatusPreValidated SerieStatus = 12
	// SerieStatusValidated represents a serie which is validated
	SerieStatusValidated SerieStatus = 16
)
