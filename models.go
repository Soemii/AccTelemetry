package AccTelemetry

import "time"

type OutboundMessage byte
type InboundMessage byte

type DriverCategory byte
type LapType byte
type CarLocation byte
type SessionPhase byte
type SessionType byte
type EventType byte
type Nationality byte
type CupCategory byte

type CarModel byte
type TrackId byte

type TrackData struct {
	Id         TrackId
	Name       string
	Meters     int32
	CameraSets map[string][]string
	HudPages   []string
}

type CarInfo struct {
	Id              uint16
	Model           CarModel
	TeamName        string
	RaceNumber      int32
	CupCategory     CupCategory
	CurrentDriverId int8
	Drivers         []DriverInfo
	Nationality     Nationality
}

func (c CarInfo) GetCurrentDriver() DriverInfo {
	if c.CurrentDriverId > int8(len(c.Drivers)) {
		return c.Drivers[c.CurrentDriverId]
	}
	return DriverInfo{}
}

type BroadCastEvent struct {
	Type   EventType
	Msg    string
	TimeMs int32
	CarId  int32
}

type LapInfo struct {
	LapTimeMs      int32
	Splits         []int32
	CarIndex       uint16
	DriverIndex    uint16
	IsInvalid      bool
	IsValidForBest bool
	LapType        LapType
}

type DriverInfo struct {
	FirstName   string
	LastName    string
	ShortName   string
	Category    DriverCategory
	Nationality Nationality
}

type RealTimeCarUpdate struct {
	CarIndex       uint16
	DriverIndex    uint16
	Gear           byte
	WorldPosX      float32
	WorldPosY      float32
	Yaw            float32
	CarLocation    CarLocation
	Kmh            uint16
	Position       uint16
	TrackPosition  uint16
	SplinePosition float32
	Delta          int32
	BestSessionLap LapInfo
	LastLap        LapInfo
	CurrentLap     LapInfo
	Laps           uint16
	CupPosition    uint16
	DriverCount    byte
}

type RealTimeUpdate struct {
	EventIndex           uint16
	SessionIndex         uint16
	Phase                SessionPhase
	SessionTime          time.Time     //TODO: CHECK IF TYPE CORRECT
	RemainingTime        time.Duration //TODO: CHECK IF TYPE CORRECT
	TimeOfDay            time.Time     //TODO: CHECK IF TYPE CORRECT
	RainLevel            float32
	Clouds               float32
	Wetness              float32
	BestSessionLap       LapInfo
	BestLapCarIndex      uint16
	BestLapDriverIndex   uint16
	FocusedCarIndex      int32
	ActiveCameraSet      string
	ActiveCamera         string
	IsReplaying          bool
	ReplaySessionTime    time.Time
	ReplayRemainingTime  time.Duration
	SessionRemainingTime time.Duration //TODO: CHECK IF TYPE CORRECT
	SessionEndTime       time.Time     //TODO: CHECK IF TYPE CORRECT
	SessionType          SessionType
	AmbientTemp          byte
	TrackTemp            byte
	CurrentHudPage       string
}

//CONSTS

const (
	OutboundMessageRegisterCommandApplication   OutboundMessage = 1
	OutboundMessageUnregisterCommandApplication OutboundMessage = 9
	OutboundMessageRequestEntryList             OutboundMessage = 10
	OutboundMessageRequestTrackData             OutboundMessage = 11
	OutboundMessageChangeHudPage                OutboundMessage = 49
	OutboundMessageChangeFocus                  OutboundMessage = 50
	OutboundMessageInstantReplayRequest         OutboundMessage = 51
	OutboundMessagePlayManualReplayHighlight    OutboundMessage = 52 //ACC: TODO, but planned
	OutboundMessageSaveManualReplayHighlight    OutboundMessage = 60 //ACC: TODO, but planned: saving manual replays gives distributed clients the possibility to see the play the same replay
)

const (
	InboundMessageRegistrationResult InboundMessage = iota + 1
	InboundMessageRealtimeUpdate
	InboundMessageRealtimeCarUpdate
	InboundMessageEntryList
	InboundMessageTrackData
	InboundMessageEntryListCar
	InboundMessageBroadcastingEvent
)

const (
	DriverCategoryBronze DriverCategory = iota
	DriverCategorySilver
	DriverCategoryGold
	DriverCategoryPlatinum
	DriverCategoryError DriverCategory = 255
)

const (
	CupCategoryPro CupCategory = iota
	CupCategoryProAm
	CupCategoryAm
	CupCategorySilver
	CupCategoryNational
)

const (
	LapTypeERROR LapType = iota
	LapTypeOutlap
	LapTypeRegular
	LapTypeInlap
)

const (
	CarLocationNONE CarLocation = iota
	CarLocationTrack
	CarLocationPitlane
	CarLocationPitEntry
	CarLocationPitExit
)

const (
	SessionPhaseNONE SessionPhase = iota
	SessionPhaseStarting
	SessionPhasePreFormation
	SessionPhaseFormationLap
	SessionPhasePreSession
	SessionPhaseSession
	SessionPhaseSessionOver
	SessionPhasePostSession
	SessionPhaseResultUI
)

const (
	SessionTypePractice        SessionType = 0
	SessionTypeQualifying      SessionType = 4
	SessionTypeSuperpole       SessionType = 9
	SessionTypeRace            SessionType = 10
	SessionTypeHotlap          SessionType = 11
	SessionTypeHotstint        SessionType = 12
	SessionTypeHotlapSuperpole SessionType = 13
	SessionTypeReplay          SessionType = 14
)

const (
	EventTypeNone            EventType = iota
	EventTypeGreenFlag                 // !!! Never send out (last checked: ACC v1.3.12)
	EventTypeSessionOver               // !!! Never send out (last checked: ACC v1.3.12)
	EventTypePenaltyCommMsg            // !!! Never send out (last checked: ACC v1.3.12)
	EventTypeAccident                  // !!! Never send out (last checked: ACC v1.3.12)
	EventTypeLapCompleted              // self-explanatory
	EventTypeBestSessionLap            // self-explanatory
	EventTypeBestPersonalLap           // self-explanatory
)

const (
	NationalityAny Nationality = iota
	NationalityItaly
	NationalityGermany
	NationalityFrance
	NationalitySpain
	NationalityGreatBritain
	NationalityHungary
	NationalityBelgium
	NationalitySwitzerland
	NationalityAustria
	NationalityRussia
	NationalityThailand
	NationalityNetherlands
	NationalityPoland
	NationalityArgentina
	NationalityMonaco
	NationalityIreland
	NationalityBrazil
	NationalitySouthAfrica
	NationalityPuertoRico
	NationalitySlovakia
	NationalityOman
	NationalityGreece
	NationalitySaudiArabia
	NationalityNorway
	NationalityTurkey
	NationalitySouthKorea
	NationalityLebanon
	NationalityArmenia
	NationalityMexico
	NationalitySweden
	NationalityFinland
	NationalityDenmark
	NationalityCroatia
	NationalityCanada
	NationalityChina
	NationalityPortugal
	NationalitySingapore
	NationalityIndonesia
	NationalityUSA
	NationalityNewZealand
	NationalityAustralia
	NationalitySanMarino
	NationalityUAE
	NationalityLuxembourg
	NationalityKuwait
	NationalityHongKong
	NationalityColombia
	NationalityJapan
	NationalityAndorra
	NationalityAzerbaijan
	NationalityBulgaria
	NationalityCuba
	NationalityCzechRepublic
	NationalityEstonia
	NationalityGeorgia
	NationalityIndia
	NationalityIsrael
	NationalityJamaica
	NationalityLatvia
	NationalityLithuania
	NationalityMacau
	NationalityMalaysia
	NationalityNepal
	NationalityNewCaledonia
	NationalityNigeria
	NationalityNorthernIreland
	NationalityPapuaNewGuinea
	NationalityPhilippines
	NationalityQatar
	NationalityRomania
	NationalityScotland
	NationalitySerbia
	NationalitySlovenia
	NationalityTaiwan
	NationalityUkraine
	NationalityVenezuela
	NationalityWales
	NationalityIran
	NationalityBahrain
	NationalityZimbabwe
	NationalityChineseTaipei
	NationalityChile
	NationalityUruguay
	NationalityMadagascar
)

//TODO: CHECK

const (
	CarModelMercedes       CarModel = 1
	CarModelFerrari        CarModel = 2
	CarModelLexus          CarModel = 15
	CarModelLamborghini    CarModel = 16
	CarModelAudi           CarModel = 19
	CarModelAstonMartin    CarModel = 20
	CarModelPorsche        CarModel = 23
	CarModelAlpineGT4      CarModel = 50
	CarModelAstonMartinGT4 CarModel = 51
	CarModelBMWGT4         CarModel = 53
	CarModelChevroletGT4   CarModel = 55
	CarModelKTMGT4         CarModel = 57
	CarModelMaseratiGT4    CarModel = 58
	CarModelMcLarenGT4     CarModel = 59
	CarModelMercedesGT4    CarModel = 60
	CarModelPorscheGT4     CarModel = 61
)

const (
	TrackIdBrandsHatch TrackId = iota + 1
	TrackIdSpa
	TrackIdMonza
	TrackIdMisano
	TrackIdPaulRicard
	TrackIdSilverstone
	TrackIdHungaroring
	TrackIdNurburgring
	TrackIdBarcelona
	TrackIdZolder
	TrackIdZandvoort
	TrackIdBathurst
	TrackIdLagunaSeca
	TrackIdSuzuka
)
