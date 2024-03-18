package bunkers

import (
	"math/rand"
	"strconv"
	"strings"

	"rerere.org/fedi-games/domain/models"
	"rerere.org/fedi-games/games"
)

type BunkersGame struct{}

type BunkersGameState struct {
	InitTerrain Terrain
	PosA        int
	PosB        int
	PlayerA     string
	PlayerB     string
	WhosTurn    int // 1 = PlayerA, 2 = PlayerB
	Shots       []Shot
	Init        bool
}

type BunkersGameStep struct {
	Player int // 1 == PlayerA, 2 == PlayerB
	Angle  float64
	Vel    float64
}

type BunkersGameResult struct {
	Winner int // 0 == None, 1 == PlayerA, 2 == PlayerB
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
		StartY: s.Terrain().Height[startX] + 15,
		Vel:    step.Vel,
		Angle:  angle,
	}

	// check if shot hits a bunker
	_, _, hit := shot.getImpact(s.Terrain())

	// add shot to state
	s.Shots = append(s.Shots, shot)

	return BunkersGameResult{
		Winner: hit,
	}
}

func (t *BunkersGame) OnMsg(session *models.GameSession, msg games.GameMsg) (interface{}, games.GameReply, error) {
	state := session.Data.(*BunkersGameState)

	// initialize the game
	if state.Init == false {
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
	found := 0
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
		return state, games.GameReply{
			To:  []string{msg.From},
			Msg: "You must include 'power' and 'angle' in your message followed by a number. Example: angle 10 power 20",
		}, nil
	}

	step := BunkersGameStep{
		Player: state.WhosTurn,
		Vel:    vel,
		Angle:  angle,
	}

	result := state.Step(step)
	state.WhosTurn = (state.WhosTurn % 2) + 1

	if result.Winner != 0 {
		return state, games.GameReply{
			To:  []string{state.PlayerA, state.PlayerB},
			Msg: "TODO: Someone won!",
		}, nil
	} else {
		return state, games.GameReply{
			To:  []string{state.PlayerA, state.PlayerB},
			Msg: "TODO: Next Turn!",
		}, nil
	}
}
