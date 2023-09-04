package types

type Provider struct {
	Name      string
	AppID     int
	Version   int
	Timestamp int64
}

type MapData struct {
	Name                 string
	MatchID              string
	GameTime             int
	ClockTime            int
	Daytime              bool
	NightstalkerNight    bool
	RadiantScore         int
	DireScore            int
	GameState            string
	Paused               bool
	WinTeam              string
	CustomGameName       string
	WardPurchaseCooldown int
}

type Wearables struct {
	Wearable0  int
	Wearable1  int
	Wearable2  int
	Wearable3  int
	Wearable4  int
	Wearable5  int
	Wearable6  int
	Wearable7  int
	Wearable8  int
	Wearable9  int
	Wearable10 int
	Wearable11 int
}

type Entity struct {
	XPos          int
	YPos          int
	Image         string
	Team          int
	Yaw           int
	UnitName      string
	VisionRange   int
	Name          *string
	EventDuration *int
	YPosP         *string // parser
	XPosP         *string // parser
	TeamP         *string // parser
}

type Minimap map[string]Entity

type Player struct {
	Team2 *struct {
		Player0 Player
		Player1 Player
		Player2 Player
		Player3 Player
		Player4 Player
	}
	Team3 *struct {
		Player5 Player
		Player6 Player
		Player7 Player
		Player8 Player
		Player9 Player
	}
	SteamID             string
	AccountID           string
	Name                string
	Activity            string
	Kills               int
	Deaths              int
	Assists             int
	LastHits            int
	Denies              int
	KillStreak          int
	CommandsIssued      int
	KillList            map[string]int
	TeamName            string
	Gold                int
	GoldReliable        int
	GoldUnreliable      int
	GoldFromHeroKills   int
	GoldFromCreepKills  int
	GoldFromIncome      int
	GoldFromShared      int
	GPM                 int
	XPM                 int
	NetWorth            int
	HeroDamage          int
	TowerDamage         int
	WardsPurchased      int
	WardsPlaced         int
	WardsDestroyed      int
	RunesActivated      int
	CampsStacked        int
	SupportGoldSpent    int
	ConsumableGoldSpent int
	ItemGoldSpent       int
	GoldLostToDeath     int
	GoldSpentOnBuybacks int
}

type Hero struct {
	Team2 *struct {
		Player0 Hero
		Player1 Hero
		Player2 Hero
		Player3 Hero
		Player4 Hero
	}
	Team3 *struct {
		Player5 Hero
		Player6 Hero
		Player7 Hero
		Player8 Hero
		Player9 Hero
	}
	ID              int
	Name            any
	XPos            *int
	YPos            *int
	Level           *int
	XP              *int
	Alive           *bool
	RespawnSeconds  *int
	BuybackCost     *int
	BuybackCooldown *int
	Health          *int
	MaxHealth       *int
	HealthPercent   *int
	Mana            *int
	MaxMana         *int
	ManaPercent     *int
	Silenced        *bool
	Stunned         *bool
	Disarmed        *bool
	Magicimmune     *bool
	Hexed           *bool
	Muted           *bool
	Break           *bool
	AghanimsScepter *bool
	AghanimsShard   *bool
	Smoked          *bool
	HasDebuff       *bool
	SelectedUnit    *bool // Only available as spectator
	Talent1         *bool
	Talent2         *bool
	Talent3         *bool
	Talent4         *bool
	Talent5         *bool
	Talent6         *bool
	Talent7         *bool
	Talent8         *bool
}

type Abilities struct {
	Ability0  *Ability
	Ability1  *Ability
	Ability2  *Ability
	Ability3  *Ability
	Ability4  *Ability
	Ability5  *Ability
	Ability6  *Ability
	Ability7  *Ability
	Ability8  *Ability
	Ability9  *Ability
	Ability10 *Ability
	Ability11 *Ability
	Ability12 *Ability
	Ability13 *Ability
	Ability14 *Ability
	Ability15 *Ability
	Ability16 *Ability
	Ability17 *Ability
	Ability18 *Ability
	Ability19 *Ability
}

type Ability struct {
	Name           string
	Level          int
	CanCast        bool
	Passive        bool
	AbilityActive  bool
	Cooldown       int
	Ultimate       bool
	Charges        int
	MaxCharges     int
	ChargeCooldown int
}

type Items struct {
	Slot0     *Item
	Slot1     *Item
	Slot2     *Item
	Slot3     *Item
	Slot4     *Item
	Slot5     *Item
	Slot6     *Item
	Slot7     *Item
	Slot8     *Item
	Stash0    *Item
	Stash1    *Item
	Stash2    *Item
	Stash3    *Item
	Stash4    *Item
	Stash5    *Item
	Teleport0 *Item
	Neutral0  *Item
}

type Item struct {
	Name      string
	Purchaser *int
	CanCast   *bool
	Cooldown  *int
	ItemLevel *int
	Passive   bool
	Charges   *int
}

type Buildings struct {
	DotaBadguysTower1Top Building
	DotaBadguysTower2Top Building
	DotaBadguysTower3Top Building
	DotaBadguysTower1Mid Building
	DotaBadguysTower2Mid Building
	DotaBadguysTower3Mid Building
	DotaBadguysTower1Bot Building
	DotaBadguysTower2Bot Building
	DotaBadguysTower3Bot Building
	DotaBadguysTower4Top Building
	DotaBadguysTower4Bot Building
	BadRaxMeleeTop       Building
	BadRaxRangeTop       Building
	BadRaxMeleeMid       Building
	BadRaxRangeMid       Building
	BadRaxMeleeBot       Building
	BadRaxRangeBot       Building
	DotaBadguysFort      Building
}

type Building struct {
	Health    int
	MaxHealth int
}

type Draft struct {
	Activeteam              int
	Pick                    bool
	ActiveteamTimeRemaining int
	RadiantBonusTime        int
	DireBonusTime           int
	Team2                   TeamDraft
	Team3                   TeamDraft
}

type TeamDraft struct {
	Pick0ID    int
	Pick0Class string
	Pick1ID    int
	Pick1Class string
	Pick2ID    int
	Pick2Class string
	Pick3ID    int
	Pick3Class string
	Pick4ID    int
	Pick4Class string
	Ban0ID     int
	Ban0Class  string
	Ban1ID     int
	Ban1Class  string
	Ban2ID     int
	Ban2Class  string
	Ban3ID     int
	Ban3Class  string
	Ban4ID     int
	Ban4Class  string
	Ban5ID     int
	Ban5Class  string
	Ban6ID     int
	Ban6Class  string
}

type DotaEventTypes string

const (
	RoshanKilled  DotaEventTypes = "roshan_killed"
	AegisPickedUp DotaEventTypes = "aegis_picked_up"
	AegisDenied   DotaEventTypes = "aegis_denied"
	Tip           DotaEventTypes = "tip"
	BountyPickup  DotaEventTypes = "bounty_rune_pickup"
	CourierKilled DotaEventTypes = "courier_killed"
)

var ValidEventTypes = map[DotaEventTypes]bool{
	RoshanKilled:  true,
	AegisPickedUp: true,
	AegisDenied:   true,
	Tip:           true,
	BountyPickup:  true,
}

type DotaEvent struct {
	GameTime         int
	EventType        DotaEventTypes
	SenderPlayerID   int
	ReceiverPlayerID int
	TipAmount        int
	CourierTeam      string
	KillerPlayerID   int
	OwningPlayerID   int
	PlayerID         int
	Team             string
	BountyValue      int
	TeamGold         int
	KilledByTeam     string
	Snatched         bool
}

type Auth struct {
	Token string `json:"token"`
}

type Packet struct {
	Provider  Provider
	Map       *MapData
	Player    *Player
	Minimap   Minimap
	Hero      *Hero
	Abilities *Abilities
	Items     *Items
	Buildings *struct {
		Radiant Buildings
		Dire    Buildings
	}
	Draft      *Draft
	Events     []DotaEvent
	Previously *Packet
	Added      *Packet
	Auth       *Auth
}

type GCMatchData struct {
	Result int
	Match  *Match
	Vote   int
}

type Match struct {
	Duration            int
	Starttime           int
	Players             []Player
	MatchID             ID
	TowerStatus         []int
	BarracksStatus      []int
	Cluster             int
	FirstBloodTime      int
	ReplaySalt          int
	ServerIP            *string
	ServerPort          *string
	LobbyType           int
	HumanPlayers        int
	AverageSkill        *int
	GameBalance         *int
	RadiantTeamID       *int
	DireTeamID          *int
	LeagueID            int
	RadiantTeamName     *string
	DireTeamName        *string
	RadiantTeamLogo     *string
	DireTeamLogo        *string
	RadiantTeamLogoURL  *string
	DireTeamLogoURL     *string
	RadiantTeamComplete *string
	DireTeamComplete    *string
	PositiveVotes       int
	NegativeVotes       int
	GameMode            int
	PicksBans           []interface{}
	MatchSeqNum         *int
	ReplayState         int
	RadiantGuildID      *string
	DireGuildID         *string
	RadiantTeamTag      *string
	DireTeamTag         *string
	SeriesID            int
	SeriesType          int
	BroadcasterChannels []interface{}
	Engine              int
	CustomGameData      *string
	MatchFlags          int
	PrivateMetadataKey  *string
	RadiantTeamScore    int
	DireTeamScore       int
	MatchOutcome        int
	TournamentID        *string
	TournamentRound     *string
	PreGameDuration     int
	MVPAccountID        []interface{}
	Coaches             []interface{}
	Level               string
	Timestamp           string
}

type ID struct {
	Low      int
	High     int
	Unsigned bool
}

type HeroDamage struct {
	PreReduction  int
	PostReduction int
	DamageType    int
}

type PermanentBuff struct {
	PermanentBuff int
	StackCount    int
	GrantTime     int
}
