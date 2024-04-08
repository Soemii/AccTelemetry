package AccTelemetry

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/Soemii/goptional"
	"log"
	"time"
)

type Number interface {
	byte | int8 | int16 | uint16 | int32 | uint32 | int64 | uint64 | uint | int | float32 | float64
}

type ReadBuffer interface {
	bool | []bool | Number
}

type AccModels interface {
	OutboundMessage | InboundMessage | SessionType | SessionPhase | Nationality | CarLocation | CarModel | EventType | TrackId | DriverCategory | CupCategory
}

const BroadcastingProtocolVersion byte = 4
const InvalidSectorTime = (2 << 30) - 1

type EntryList []uint16

func writeString(b *bytes.Buffer, data string) error {
	length := int16(len(data))
	lengthErr := binary.Write(b, binary.LittleEndian, length)
	if lengthErr != nil {
		return lengthErr
	}
	n, dataErr := b.Write([]byte(data))
	if dataErr != nil {
		return dataErr
	}
	if n != len(data) {
		return errors.New("mismatch of length written bytes")
	}
	return nil
}

func readString(b *bytes.Buffer) (s string, err error) {
	length, err := readNumber[int16](b)
	buff := make([]byte, length)
	n, err := b.Read(buff)
	if int16(n) != length {
		err = errors.New("mis")
	}
	s = string(buff)
	return
}

func writeNumber[d ReadBuffer | AccModels](b *bytes.Buffer, data d) error {
	return binary.Write(b, binary.LittleEndian, data)
}

func readNumber[d ReadBuffer | AccModels](b *bytes.Buffer) (data d, err error) {
	err = binary.Read(b, binary.LittleEndian, &data)
	return
}

func divideByTen(b byte, err error) (float32, error) {
	return float32(b) / 10, err
}

func parseTime(float322 float32, err error) (time.Time, error) {
	log.Println(float322)
	return time.Now(), errors.New("cannot convert float32 to time")
}

func parseDuration(float322 float32, err error) (time.Duration, error) {
	log.Println(float322)
	return time.Duration(0), errors.New("cannot convert float32 to duration")
}

func readLap(b *bytes.Buffer) (lap LapInfo, err error) {
	var splitCount byte
	var split int32
	var isOutLap bool
	var isInLap bool
	lap.LapTimeMs, err = readNumber[int32](b)
	if err != nil {
		return
	}
	lap.CarIndex, err = readNumber[uint16](b)
	if err != nil {
		return
	}
	lap.DriverIndex, err = readNumber[uint16](b)
	if err != nil {
		return
	}
	splitCount, err = readNumber[byte](b)
	if err != nil {
		return
	}
	for i := byte(0); i < splitCount; i++ {
		split, err = readNumber[int32](b)
		if err != nil {
			return
		}
		lap.Splits = append(lap.Splits, split)
	}
	lap.IsInvalid, err = readNumber[bool](b)
	if err != nil {
		return
	}
	lap.IsValidForBest, err = readNumber[bool](b)
	if err != nil {
		return
	}
	isOutLap, err = readNumber[bool](b)
	if err != nil {
		return
	}
	isInLap, err = readNumber[bool](b)
	if err != nil {
		return
	}
	if isOutLap {
		lap.LapType = LapTypeOutlap
	} else if isInLap {
		lap.LapType = LapTypeInlap
	} else {
		lap.LapType = LapTypeRegular
	}
	for len(lap.Splits) < 3 {
		lap.Splits = append(lap.Splits, 0)
	}
	return
}

func readConnectionResponse(b *bytes.Buffer) (connectionId int32, connectionSuccess bool, isReadOnly bool, err error) {
	connectionId, err = readNumber[int32](b)
	if err != nil {
		return
	}
	connectionSuccess, err = readNumber[bool](b)
	if err != nil {
		return
	}
	isReadOnly, err = readNumber[bool](b)
	if err != nil {
		return
	}
	errMsg, err := readString(b)
	if err != nil {
		return
	}
	if errMsg != "" {
		err = errors.New(errMsg)
	}
	return
}

func readEntryListResponse(b *bytes.Buffer) (connectionId int32, entryList EntryList, err error) {
	connectionId, err = readNumber[int32](b)
	if err != nil {
		return
	}
	entryCount, err := readNumber[uint16](b)
	if err != nil {
		return
	}
	entryList = make(EntryList, entryCount)
	for i := uint16(0); i < entryCount; i++ {
		var entry uint16
		entry, err = readNumber[uint16](b)
		if err != nil {
			return
		}
		entryList[i] = entry
	}
	return
}

func readEntryListCarResponse(b *bytes.Buffer) (car CarInfo, err error) {
	car.Id, err = readNumber[uint16](b)
	if err != nil {
		return
	}
	car.Model, err = readNumber[CarModel](b)
	if err != nil {
		return
	}
	car.TeamName, err = readString(b)
	if err != nil {
		return
	}
	car.RaceNumber, err = readNumber[int32](b)
	if err != nil {
		return
	}
	car.CupCategory, err = readNumber[CupCategory](b)
	if err != nil {
		return
	}
	car.CurrentDriverId, err = readNumber[int8](b)
	if err != nil {
		return
	}
	car.Nationality, err = readNumber[Nationality](b)
	if err != nil {
		return
	}
	var driversCount uint8
	driversCount, err = readNumber[uint8](b)
	if err != nil {
		return
	}
	car.Drivers = make([]DriverInfo, driversCount)
	for i := uint8(0); i < driversCount; i++ {
		car.Drivers[i].FirstName, err = readString(b)
		if err != nil {
			return
		}
		car.Drivers[i].LastName, err = readString(b)
		if err != nil {
			return
		}
		car.Drivers[i].ShortName, err = readString(b)
		if err != nil {
			return
		}
		car.Drivers[i].Category, err = readNumber[DriverCategory](b)
		if err != nil {
			return
		}
		car.Drivers[i].Nationality, err = readNumber[Nationality](b)
		if err != nil {
			return
		}
	}
	return
}

func readRealtimeUpdateResponse(b *bytes.Buffer) (update RealTimeUpdate, err error) {
	update.EventIndex, err = readNumber[uint16](b)
	update.SessionIndex, err = readNumber[uint16](b)
	update.SessionType, err = readNumber[SessionType](b)
	update.Phase, err = readNumber[SessionPhase](b)
	update.SessionTime, err = parseTime(readNumber[float32](b))
	update.SessionEndTime, err = parseTime(readNumber[float32](b))
	update.FocusedCarIndex, err = readNumber[int32](b)
	update.ActiveCameraSet, err = readString(b)
	update.ActiveCamera, err = readString(b)
	update.CurrentHudPage, err = readString(b)
	update.IsReplaying, err = readNumber[bool](b)
	if update.IsReplaying {
		update.ReplaySessionTime, err = parseTime(readNumber[float32](b))
		update.ReplayRemainingTime, err = parseDuration(readNumber[float32](b))
	}
	update.TimeOfDay, err = parseTime(readNumber[float32](b))
	update.AmbientTemp, err = readNumber[byte](b)
	update.TrackTemp, err = readNumber[byte](b)
	update.Clouds, err = divideByTen(readNumber[byte](b))
	update.RainLevel, err = divideByTen(readNumber[byte](b))
	update.Wetness, err = divideByTen(readNumber[byte](b))
	update.BestSessionLap, err = readLap(b)
	return
}

func readRealtimeCarUpdateResponse(b *bytes.Buffer) (update RealTimeCarUpdate, err error) {
	update.CarIndex, err = readNumber[uint16](b)
	if err != nil {
		return
	}
	update.DriverIndex, err = readNumber[uint16](b)
	if err != nil {
		return
	}
	update.DriverCount, err = readNumber[byte](b)
	if err != nil {
		return
	}
	gear, err := readNumber[byte](b)
	if err != nil {
		return
	}
	update.Gear = gear - 2
	update.WorldPosX, err = readNumber[float32](b)
	if err != nil {
		return
	}
	update.WorldPosY, err = readNumber[float32](b)
	if err != nil {
		return
	}
	update.Yaw, err = readNumber[float32](b)
	if err != nil {
		return
	}
	update.CarLocation, err = readNumber[CarLocation](b)
	if err != nil {
		return
	}
	update.Kmh, err = readNumber[uint16](b)
	if err != nil {
		return
	}
	update.Position, err = readNumber[uint16](b)
	if err != nil {
		return
	}
	update.CupPosition, err = readNumber[uint16](b)
	if err != nil {
		return
	}
	update.TrackPosition, err = readNumber[uint16](b)
	if err != nil {
		return
	}
	update.SplinePosition, err = readNumber[float32](b)
	if err != nil {
		return
	}
	update.Laps, err = readNumber[uint16](b)
	if err != nil {
		return
	}
	update.Delta, err = readNumber[int32](b)
	if err != nil {
		return
	}
	update.BestSessionLap, err = readLap(b)
	if err != nil {
		return
	}
	update.LastLap, err = readLap(b)
	if err != nil {
		return
	}
	update.CurrentLap, err = readLap(b)
	if err != nil {
		return
	}
	return
}

func readTrackDataResponse(b *bytes.Buffer) (connectionId int32, trackData TrackData, err error) {
	var cameraCount byte
	var cameraSetCount byte
	var hudCount byte
	var camSetName string
	var cameraName string
	var hudName string
	connectionId, err = readNumber[int32](b)
	if err != nil {
		return
	}
	trackData.Name, err = readString(b)
	if err != nil {
		return
	}
	trackData.Id, err = readNumber[TrackId](b)
	if err != nil {
		return
	}
	trackData.Meters, err = readNumber[int32](b)
	if err != nil {
		return
	}
	cameraSetCount, err = readNumber[byte](b)
	if err != nil {
		return
	}
	trackData.CameraSets = make(map[string][]string)
	for i := byte(0); i < cameraSetCount; i++ {
		camSetName, err = readString(b)
		if err != nil {
			return
		}
		cameraCount, err = readNumber[byte](b)
		if err != nil {
			return
		}
		for j := byte(0); j < cameraCount; j++ {
			cameraName, err = readString(b)
			if err != nil {
				return
			}
			trackData.CameraSets[camSetName] = append(trackData.CameraSets[camSetName], cameraName)
		}
	}
	hudCount, err = readNumber[byte](b)
	for i := byte(0); i < hudCount; i++ {
		hudName, err = readString(b)
		if err != nil {
			return
		}
		trackData.HudPages = append(trackData.HudPages, hudName)
	}
	return
}

func readBroadcastingEvent(b *bytes.Buffer) (event BroadCastEvent, err error) {
	event.Type, err = readNumber[EventType](b)
	if err != nil {
		return
	}
	event.Msg, err = readString(b)
	if err != nil {
		return
	}
	event.TimeMs, err = readNumber[int32](b)
	if err != nil {
		return
	}
	event.CarId, err = readNumber[int32](b)
	if err != nil {
		return
	}
	return
}

func writeRegistration(b *bytes.Buffer, displayName string, connectionPassword string, msRealtimeUpdateInterval int32, commandPassword string) error {
	if err := writeNumber(b, OutboundMessageRegisterCommandApplication); err != nil {
		return err
	}
	if err := writeNumber(b, BroadcastingProtocolVersion); err != nil {
		return err
	}
	if err := writeString(b, displayName); err != nil {
		return err
	}
	if err := writeString(b, connectionPassword); err != nil {
		return err
	}
	if err := writeNumber(b, msRealtimeUpdateInterval); err != nil {
		return err
	}
	if err := writeString(b, commandPassword); err != nil {
		return err
	}
	return nil
}

func writeDisconnect(b *bytes.Buffer, connectionId int32) error {
	if err := writeNumber(b, OutboundMessageUnregisterCommandApplication); err != nil {
		return err
	}
	if err := writeNumber(b, connectionId); err != nil {
		return err
	}
	return nil
}

func writeEntryListRequest(b *bytes.Buffer, connectionId int32) (err error) {
	if err = writeNumber(b, OutboundMessageRequestEntryList); err != nil {
		return
	}
	if err = writeNumber(b, connectionId); err != nil {
		return
	}
	return
}

func writeTrackDataRequest(b *bytes.Buffer, connectionId int32) (err error) {
	if err = writeNumber(b, OutboundMessageRequestTrackData); err != nil {
		return
	}
	if err = writeNumber(b, connectionId); err != nil {
		return
	}
	return
}

func writeFocusedCar(b *bytes.Buffer, connectionId int32, carIndex goptional.Optional[uint16], cameraSet goptional.Optional[string], camera goptional.Optional[string]) (err error) {
	if err = writeNumber(b, OutboundMessageChangeFocus); err != nil {
		return
	}
	if err = writeNumber(b, connectionId); err != nil {
		return
	}
	if err = writeNumber(b, carIndex.Present()); err != nil {
		return
	}
	err = carIndex.If(func(carId uint16) (err error) {
		if err = writeNumber(b, carId); err != nil {
			return
		}
		return
	})
	if err != nil {
		return
	}
	if err = writeNumber(b, cameraSet.Present() && camera.Present()); err != nil {
		return
	}
	if cameraSet.Present() && camera.Present() {
		if err = writeString(b, cameraSet.Get()); err != nil {
			return
		}
		if err = writeString(b, camera.Get()); err != nil {
			return
		}
	}
	return
}

func writeRequestInstantReplay(b *bytes.Buffer, connectionId int32, startSessionTime float32, durationMS float32, initialFocusedCarIndex int32, initialCameraSet string, initialCamera string) (err error) {
	if err = writeNumber(b, OutboundMessageInstantReplayRequest); err != nil {
		return
	}
	if err = writeNumber(b, connectionId); err != nil {
		return
	}
	if err = writeNumber(b, startSessionTime); err != nil {
		return
	}
	if err = writeNumber(b, durationMS); err != nil {
		return
	}
	if err = writeNumber(b, initialFocusedCarIndex); err != nil {
		return
	}
	if err = writeString(b, initialCameraSet); err != nil {
		return
	}
	if err = writeString(b, initialCamera); err != nil {
		return
	}
	return
}

func writeRequestHudPage(b *bytes.Buffer, connectionId int32, hudPage string) (err error) {
	if err = writeNumber(b, OutboundMessageInstantReplayRequest); err != nil {
		return
	}
	if err = writeNumber(b, connectionId); err != nil {
		return
	}
	if err = writeString(b, hudPage); err != nil {
		return
	}
	return
}
