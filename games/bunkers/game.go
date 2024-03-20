package bunkers

import (
	"math"
	"math/rand"
	"strconv"
	"strings"

	"rerere.org/fedi-games/domain/models"
	"rerere.org/fedi-games/games"
	"rerere.org/fedi-games/internal/acpub"
)

func NewBunkerGame() games.Game {
	return &BunkersGame{}
}

type BunkersGame struct{}

// Example implements games.Game.
func (t *BunkersGame) Example() string {
	return "@alice@example.com power 13 angle 45"
}

// Name implements games.Game.
func (BunkersGame) Name() string {
	return "bunkers"
}

// NewState implements games.Game.
func (t *BunkersGame) NewState() interface{} {
	return &BunkersGameState{}
}

// Summary implements games.Game.
func (t *BunkersGame) Summary() string {
	return `It's the year 2169, the world is ravaged by runaway global warming.<br>
	The only lives left are clones of billionaires in their apocalypse bunkers.<br>
	Their prime directives leave them no other choice than to fight until their are the only one left.
	`
}

type BunkersGameState struct {
	InitTerrain Terrain
	PosA        int
	PosB        int
	PlayerA     string
	PlayerB     string
	WhosTurn    int // 1 = PlayerA, 2 = PlayerB
	Shots       []Shot
	Init        bool
	Wind        int
}

type BunkersGameStep struct {
	Player int // 1 == PlayerA, 2 == PlayerB
	Angle  float64
	Vel    float64
}

type BunkersGameResult struct {
	Winner int // 0 == None, 1 == PlayerA, 2 == PlayerB
}

func newWind() int {
	return rand.Intn(11) - 5
}

func windMsg(wind int) string {
	w := ""
	if wind == 0 {
		w = "↔️"
	}
	if wind < 0 {
		for i := 0; i < -wind; i++ {
			w += "⬅️"
		}
	}
	if wind > 0 {
		for i := 0; i < wind; i++ {
			w += "➡️"
		}
	}
	return "Wind: " + w
}

func NewBunkersGameState(a string, b string) *BunkersGameState {
	return &BunkersGameState{
		InitTerrain: NewTerrain(),
		PosA:        50 + rand.Intn(150),
		PosB:        600 + 50 + rand.Intn(150),
		PlayerA:     a,
		PlayerB:     b,
		WhosTurn:    1,
		Shots:       make([]Shot, 0),
		Init:        true,
		Wind:        newWind(),
	}
}

func (s *BunkersGameState) Terrain() Terrain {
	t := s.InitTerrain.Copy()
	if len(s.Shots) > 1 {
		for _, shot := range s.Shots[:len(s.Shots)-1] {
			t = shot.DestroyTerrain(t)
		}
	}
	return t
}

func (s *BunkersGameState) Step(step BunkersGameStep) BunkersGameResult {
	// Construct Shot
	startX := s.PosA
	angle := step.Angle
	if step.Player == 2 {
		startX = s.PosB
		angle = -angle
	}

	shot := Shot{
		StartX: startX,
		StartY: s.Terrain().At(startX) + 15,
		Vel:    step.Vel,
		Angle:  angle,
		Wind:   s.Wind,
	}

	// check if shot hits a bunker
	trail, valid := shot.getImpact(s.Terrain())
	winner := 0
	if valid {
		p := trail[len(trail)-1]
		hitX := p.X
		hitY := p.Y

		aX := s.PosA
		aY := s.Terrain().At(aX)
		bX := s.PosB
		bY := s.Terrain().At(bX)

		daX := float64(hitX - aX)
		daY := float64(hitY - aY)
		dbX := float64(hitX - bX)
		dbY := float64(hitY - bY)

		if math.Sqrt(daX*daX+daY*daY) < float64(EXPLOSION_RADIUS) {
			winner = 2
		}
		if math.Sqrt(dbX*dbX+dbY*dbY) < float64(EXPLOSION_RADIUS) {
			winner = 1
		}
	}

	// add shot to state
	s.Shots = append(s.Shots, shot)
	s.Wind = newWind()

	return BunkersGameResult{
		Winner: winner,
	}
}

func (t *BunkersGame) OnMsg(session *models.GameSession, msg games.GameMsg) (interface{}, games.GameReply, error) {
	state := session.Data.(*BunkersGameState)

	// initialize the game
	if !state.Init {
		if len(msg.To) != 1 {
			return state, games.GameReply{
				To:  []string{msg.From},
				Msg: "You must mention exactly other player",
			}, nil
		}

		state = NewBunkersGameState(
			msg.From,
			msg.To[0],
		)
	}

	// check if it's players turn
	if (state.WhosTurn == 1 && msg.From != state.PlayerA) || (state.WhosTurn == 2 && msg.From != state.PlayerB) {
		return state, games.GameReply{
			To:  []string{msg.From},
			Msg: "It's not your turn",
		}, nil
	}

	// parse  message
	parts := strings.Split(msg.Msg, " ")
	vel := 0.0
	angle := 0.0
	velFound := false
	angleFound := false
	searching := 0
	for _, p := range parts {
		if strings.ToLower(p) == "power" || strings.ToLower(p) == "velocity" {
			searching = 1
			continue
		}
		if strings.ToLower(p) == "angle" {
			searching = 2
			continue
		}

		if searching != 0 {
			value, err := strconv.ParseFloat(p, 64)
			if err != nil {
				continue
			}
			if searching == 1 {
				vel = value
				velFound = true
				searching = 0
				continue
			}
			if searching == 2 {
				angle = value
				angleFound = true
				searching = 0
				continue
			}
		}
		if angleFound && velFound {
			break
		}
	}

	// not all info given
	if !(angleFound && velFound) {
		// render state to show player the map
		img, err := Render(*state)
		if err != nil {
			return state, games.GameReply{}, err
		}

		return state, games.GameReply{
			To:  []string{msg.From},
			Msg: "You must include 'power' and 'angle' in your message followed by a number. Example: angle 45 power 60.<br>" + windMsg(state.Wind),
			Attachments: []games.GameAttachment{
				{
					Url:       img,
					MediaType: "image/png",
				},
			},
		}, nil
	}

	vel = math.Max(math.Min(vel, 1000), 0)

	step := BunkersGameStep{
		Player: state.WhosTurn,
		Vel:    vel,
		Angle:  angle,
	}

	result := state.Step(step)
	state.WhosTurn = (state.WhosTurn % 2) + 1
	state.Wind = newWind()

	img, err := Render(*state)
	if err != nil {
		return state, games.GameReply{}, err
	}

	actorA, _ := acpub.GetActor(state.PlayerA)
	actorB, _ := acpub.GetActor(state.PlayerB)

	if result.Winner != 0 {
		m := "Winner: 🎉🎉🎉 "
		if result.Winner == 1 {
			m += acpub.ActorToLink(actorA)
		} else {
			m += acpub.ActorToLink(actorB)
		}
		m += " 🎉🎉🎉"

		return state, games.GameReply{
			To:  []string{state.PlayerA, state.PlayerB},
			Msg: m,
			Attachments: []games.GameAttachment{
				{
					Url:       img,
					MediaType: "image/png",
				},
			},
		}, nil
	} else {
		m := windMsg(state.Wind) + "<br>"
		m += "🟥 " + acpub.ActorToLink(actorA) + "<br>"
		m += "🟦 " + acpub.ActorToLink(actorB) + "<br>"
		m += "Its your turn: "
		if state.WhosTurn == 1 {
			m += acpub.ActorToLink(actorA)
		} else {
			m += acpub.ActorToLink(actorB)
		}

		return state, games.GameReply{
			To:  []string{state.PlayerA, state.PlayerB},
			Msg: m,
			Attachments: []games.GameAttachment{
				{
					Url:       img,
					MediaType: "image/png",
				},
			},
		}, nil
	}
}
