package AccTelemetry

import (
	"bytes"
	"errors"
	"github.com/Soemii/goptional"
	"log"
	"net"
	"time"
)

const ReadBufferSize = 32 * 1024

func NewAccUDPClient(address string, displayName string, connectionPassword string, msRealtimeUpdateInterval int32, commandPassword string, timeoutMs int32) *AccUDPClient {
	return &AccUDPClient{
		address:                  address,
		displayName:              displayName,
		connectionPassword:       connectionPassword,
		msRealtimeUpdateInterval: msRealtimeUpdateInterval,
		commandPassword:          connectionPassword,
		timeoutMs:                timeoutMs,
	}
}

type AccUDPClient struct {
	conn *net.UDPConn

	timeOutDuration time.Duration

	connectionId int32

	listening bool

	ErrChannel                    chan error
	BroadCastEventChannel         chan BroadCastEvent
	TrackDataEventChannel         chan TrackData
	EntryListCarEventChannel      chan CarInfo
	EntryListEventChannel         chan EntryList
	RealtimeUpdateEventChannel    chan RealTimeUpdate
	RealtimeCarUpdateEventChannel chan RealTimeCarUpdate

	address                  string
	displayName              string
	connectionPassword       string
	msRealtimeUpdateInterval int32
	commandPassword          string
	timeoutMs                int32
}

func (c *AccUDPClient) Connect() (err error) {
	var addr *net.UDPAddr
	log.Printf("Try to connect to %s", c.address)
	addr, err = net.ResolveUDPAddr("udp", c.address)
	if err != nil {
		log.Printf("[Error] Cannot reslove Address: %v", err)
		return err
	}
	c.conn, err = net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Printf("[Error] Cannot connect: %v", err)
		return err
	}
	var buffer bytes.Buffer
	err = writeRegistration(&buffer, c.displayName, c.connectionPassword, c.msRealtimeUpdateInterval, c.commandPassword)
	if err != nil {
		return
	}
	err = c.sendBuffer(buffer)
	go c.listen()
	return
}

func (c *AccUDPClient) listen() {
	var buff [ReadBufferSize]byte
	for c.listening {
		err := c.conn.SetDeadline(time.Now().Add(c.timeOutDuration))
		if err != nil {
			c.ErrChannel <- err
			continue
		}
		n, err := c.conn.Read(buff[:])
		if err != nil {
			c.ErrChannel <- err
			continue
		}
		if n == ReadBufferSize {
			log.Println("buffer not big enough")
			continue
		}
		buffer := bytes.NewBuffer(buff[:n])
		inboundMessage, err := readNumber[InboundMessage](buffer)
		if err != nil {
			c.ErrChannel <- err
			continue
		}
		switch inboundMessage {
		case InboundMessageRegistrationResult:
			var success bool
			var onlyRead bool
			c.connectionId, success, onlyRead, err = readConnectionResponse(buffer)
			if err != nil {
				c.ErrChannel <- err
				continue
			}
			log.Printf("UDPConnection acknowledge ConnectionID: %v | Success: %v | onlyRead: %v", c.connectionId, success, onlyRead)
		case InboundMessageBroadcastingEvent:
			var event BroadCastEvent
			event, err = readBroadcastingEvent(buffer)
			if err != nil {
				c.ErrChannel <- err
				continue
			}
			c.BroadCastEventChannel <- event
		case InboundMessageTrackData:
			var event TrackData
			_, event, err = readTrackDataResponse(buffer)
			if err != nil {
				c.ErrChannel <- err
				continue
			}
			c.TrackDataEventChannel <- event
		case InboundMessageEntryListCar:
			var event CarInfo
			event, err = readEntryListCarResponse(buffer)
			if err != nil {
				c.ErrChannel <- err
				continue
			}
			c.EntryListCarEventChannel <- event
		case InboundMessageEntryList:
			var event EntryList
			_, event, err = readEntryListResponse(buffer)
			if err != nil {
				c.ErrChannel <- err
				continue
			}
			c.EntryListEventChannel <- event
		case InboundMessageRealtimeUpdate:
			var event RealTimeUpdate
			event, err = readRealtimeUpdateResponse(buffer)
			if err != nil {
				c.ErrChannel <- err
				continue
			}
			c.RealtimeUpdateEventChannel <- event
		case InboundMessageRealtimeCarUpdate:
			var event RealTimeCarUpdate
			event, err = readRealtimeCarUpdateResponse(buffer)
			if err != nil {
				c.ErrChannel <- err
				continue
			}
			c.RealtimeCarUpdateEventChannel <- event
		default:
			c.ErrChannel <- errors.New("unknown inboundMessagetype")
		}
	}
}

func (c *AccUDPClient) Disconnect() (err error) {
	var buffer bytes.Buffer
	err = writeDisconnect(&buffer, c.connectionId)
	if err != nil {
		return
	}
	err = c.sendBuffer(buffer)
	if err != nil {
		return
	}
	err = c.conn.Close()
	c.listening = false
	return
}

func (c *AccUDPClient) sendBuffer(buffer bytes.Buffer) (err error) {
	var n int
	n, err = c.conn.Write(buffer.Bytes())
	if n != buffer.Len() {
		return errors.New("mismatch of length written bytes")
	}
	return
}

func (c *AccUDPClient) RequestTrackData() (err error) {
	var buffer bytes.Buffer
	err = writeTrackDataRequest(&buffer, c.connectionId)
	if err != nil {
		return
	}
	err = c.sendBuffer(buffer)
	return
}

func (c *AccUDPClient) RequestEntryList() (err error) {
	var buffer bytes.Buffer
	err = writeEntryListRequest(&buffer, c.connectionId)
	if err != nil {
		return
	}
	err = c.sendBuffer(buffer)
	return
}

func (c *AccUDPClient) RequestInstantReplay(startSessionTime float32, durationMS float32, initialFocusedCarIndex int32, initialCameraSet string, initialCamera string) (err error) {
	var buffer bytes.Buffer
	err = writeRequestInstantReplay(&buffer, c.connectionId, startSessionTime, durationMS, initialFocusedCarIndex, initialCameraSet, initialCamera)
	if err != nil {
		return
	}
	err = c.sendBuffer(buffer)
	return
}

func (c *AccUDPClient) RequestHudPage(hudPage string) (err error) {
	var buffer bytes.Buffer
	err = writeRequestHudPage(&buffer, c.connectionId, hudPage)
	if err != nil {
		return
	}
	err = c.sendBuffer(buffer)
	return
}

func (c *AccUDPClient) RequestFocusedCar(carIndex goptional.Optional[uint16], cameraSet goptional.Optional[string], camera goptional.Optional[string]) (err error) {
	var buffer bytes.Buffer
	err = writeFocusedCar(&buffer, c.connectionId, carIndex, cameraSet, camera)
	if err != nil {
		return
	}
	err = c.sendBuffer(buffer)
	return
}
