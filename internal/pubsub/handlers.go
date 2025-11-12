package pubsub

import (
	"fmt"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

func HandlerPause(gs *gamelogic.GameState) func(routing.PlayingState) AckType {
	return func(r routing.PlayingState) AckType {
		defer fmt.Print("> ")
		gs.HandlePause(r)
		return Ack
	}
}

func HandlerMove(gs *gamelogic.GameState, ch *amqp.Channel) func(gamelogic.ArmyMove) AckType {
	return func(move gamelogic.ArmyMove) AckType {
		defer fmt.Print("> ")
		moveOutcome := gs.HandleMove(move)
		switch moveOutcome {
		case gamelogic.MoveOutcomeSamePlayer:
			return NackDiscard
		case gamelogic.MoveOutComeSafe:
			return Ack
		case gamelogic.MoveOutcomeMakeWar:
			err := PublishJSON(ch, routing.ExchangePerilTopic, routing.WarRecognitionsPrefix+"."+gs.Player.Username, gamelogic.RecognitionOfWar{
				Attacker: move.Player,
				Defender: gs.Player,
			})
			if err != nil {
				return NackRequeue
			}
			return Ack
		}
		fmt.Println("error: unknown move outcome")
		return NackDiscard
	}
}

func HandlerWar(gs *gamelogic.GameState, ch *amqp.Channel) func(gamelogic.RecognitionOfWar) AckType {
	return func(rw gamelogic.RecognitionOfWar) AckType {
		defer fmt.Print("> ")
		outcome, winner, loser := gs.HandleWar(rw)
		message := ""

		switch outcome {
		case gamelogic.WarOutcomeNotInvolved:
			return NackRequeue
		case gamelogic.WarOutcomeNoUnits:
			return NackDiscard
		case gamelogic.WarOutcomeOpponentWon:
			message = fmt.Sprintf("%s won a war against %s", winner, loser)
		case gamelogic.WarOutcomeYouWon:
			message = fmt.Sprintf("%s won a war against %s", winner, loser)
		case gamelogic.WarOutcomeDraw:
			message = fmt.Sprintf("A war between %s and %s resulted in a draw", winner, loser)
		default:
			fmt.Println("an error occured")
			return NackDiscard
		}

		err := PublishGameLog(routing.GameLog{
			CurrentTime: time.Now(),
			Message:     message,
			Username:    gs.GetUsername(),
		}, ch)


		if err != nil {
			return NackRequeue
		}

		return Ack
	}
}

