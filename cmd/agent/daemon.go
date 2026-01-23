package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Ctrl-Alt-GG/projectile/cmd/agent/config"
	"github.com/Ctrl-Alt-GG/projectile/cmd/agent/grpc"
	"github.com/Ctrl-Alt-GG/projectile/cmd/agent/scrapers"
	"github.com/Ctrl-Alt-GG/projectile/pkg/agentmsg"
	"github.com/Ctrl-Alt-GG/projectile/pkg/framework"
	"github.com/Ctrl-Alt-GG/projectile/pkg/model"
	"go.uber.org/zap"
	ggrpc "google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func baseGameDataFromConfig(cfg config.GameData) model.GameServerStaticData {
	return model.GameServerStaticData{
		Game:      cfg.Game,
		Name:      cfg.Name,
		Addresses: cfg.AddressesOverride,
	}
}

func initGameServerData(loggger *zap.Logger, cfg config.GameData, scraper scrapers.Scraper) model.GameServerData {
	data := model.GameServerData{ // we will re-use this
		GameServerStaticData: baseGameDataFromConfig(cfg),
	}
	if len(data.GameServerStaticData.Addresses) == 0 {
		loggger.Debug("Figuring out game server addresses...")
		// TODO!!
	}
	data.GameServerStaticData.Capabilities = scraper.Capabilities()
	return data
}

func daemon() {
	logger := framework.SetupLogger(false)
	defer logger.Sync()
	logger.Info("Starting Projectile agent...")

	cfg, err := config.LoadConfig(logger, "")
	if err != nil {
		logger.Error("Failed to load config", zap.Error(err))
		return
	}

	wg := sync.WaitGroup{}

	// Scraper part

	scraper, err := scrapers.FromConfig(cfg.Scraper)
	if err != nil {
		logger.Error("Failed to instantiate scraper from config!", zap.Error(err))
		return
	}
	updateCh := make(chan model.GameServerDynamicData, 1)
	stopCh := make(chan any)

	go func() {
		wg.Add(1)
		defer wg.Done()
		defer close(updateCh)

		l := logger.With(zap.String("src", "scraper"))
		t := time.NewTicker(cfg.Scraper.Interval)

		defer t.Stop()

		parentCtx, parentCancel := context.WithCancel(context.Background())
		defer parentCancel()

		for {
			select {
			case <-t.C:
				l.Debug("Running scraper...")
				ctx, cancel := context.WithTimeout(parentCtx, cfg.Scraper.Timeout)
				d, err := scraper.Scrape(ctx, l)
				cancel()
				if err != nil {
					l.Error("Error while running the scraper", zap.Error(err))
				} else {
					updateCh <- d
				}
			case <-stopCh:
				return
			}
		}
	}()

	// Message assembler part

	msgCh := make(chan *agentmsg.GameServer, 1)

	go func() {
		wg.Add(1)
		defer wg.Done()
		defer close(msgCh)

		l := logger.With(zap.String("src", "msg_passing"))
		data := initGameServerData(l, cfg.GameData, scraper)

		for update := range updateCh {
			data.GameServerDynamicData = update
			select {
			case msgCh <- data.ToProtobuf():
				// nothing
			default:
				l.Warn("Dropped message!")
			}
		}
	}()

	// the GRPC client

	cm := grpc.NewClientManagerFromConfig(cfg.Server)

	go func() {
		wg.Add(1)
		defer wg.Done()
		outerL := logger.With(zap.String("src", "grpcclient"))

		streamCtx, cancel := context.WithCancel(context.Background())
		defer cancel()

		var stream ggrpc.ClientStreamingClient[agentmsg.GameServer, emptypb.Empty] = nil

		recvErrCh := make(chan error, 1)

		for msg := range msgCh {
			sendOk := false
		RetryLoop:
			for try := range cfg.Server.MaxRetry { // give it 6 tries
				l := outerL.With(zap.Int("try", try))

				// first check if we received an error while we were sleeping...
				var recvErr error
				select {
				case recvErr = <-recvErrCh:
				default:
				}

				if recvErr != nil {
					// that must mean that the stream has broken
					err = stream.CloseSend()
					if err != nil {
						l.Error("Error when closing send channel", zap.Error(err))
					}
					stream = nil
				}

				if stream == nil {
					c, err := cm.GetGRPCClient(l)
					if err != nil {
						l.Error("Failed to acquire GRPC client!", zap.Error(err))
						cm.Close(l)
						time.Sleep(time.Second)
						continue RetryLoop
					}

					// we probably have a good c now...
					stream, err = c.Updates(streamCtx)
					if err != nil {
						l.Error("Failed to open Updates channel", zap.Error(err))
						cm.Close(l)
						time.Sleep(time.Second)
						continue RetryLoop
					}

					go func() {
						localErr := stream.RecvMsg(nil) // Hopefully RecvMsg always exits, when we close the stream
						if localErr != nil {
							recvErrCh <- localErr
						}
						l.Debug("RecvMsg returned")
					}()

				}

				err = stream.Send(msg)
				if err != nil {
					l.Error("Error while sending message on the GRPC stream channel!", zap.Error(err))
					err = stream.CloseSend()
					if err != nil {
						l.Error("Error when closing send channel", zap.Error(err))
					}
					stream = nil
					time.Sleep(time.Second)
					continue RetryLoop
				}

				t := time.NewTimer(time.Second) // give this much time to the server, to respond
				select {
				case recvErr = <-recvErrCh:
					l.Error("The server sent back an error, breaking the stream!", zap.Error(recvErr))
					err = stream.CloseSend()
					if err != nil {
						l.Error("Error when closing send channel", zap.Error(err))
					}
					stream = nil
					time.Sleep(time.Second)
					continue RetryLoop
				case <-t.C:

				}

				l.Debug("Message successfully sent!")
				sendOk = true
				break
			}
			if !sendOk {
				outerL.Warn("Failed to send message, and ran out of retry attempts! Message dropped!")
			}
		}

	}()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGTERM, syscall.SIGINT)

	<-signalCh

	logger.Info("Stopping agent!")
	close(stopCh)
	wg.Wait()

	logger.Debug("Sending withdraw message...")
	c, err := cm.GetGRPCClient(logger)
	if err != nil {
		logger.Error("Failed to get client to send the withdraw message!", zap.Error(err))
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	_, err = c.Withdraw(ctx, nil)
	if err != nil {
		logger.Error("Failed to send the withdraw message", zap.Error(err))
	}

	cm.Close(logger)
	logger.Info("Good bye!")
}
