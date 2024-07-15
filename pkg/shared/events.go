package shared

import (
	"time"

	"github.com/google/uuid"
)

type ObservableType int

const (
	ObservableTypeHost ObservableType = iota
	ObservableTypeAgent
	ObservableTypeProcess
	ObservableTypeFile
	ObservableTypeEvent
)

type Predicate int

const (
	PredicateStarted Predicate = iota
	PredicateStopped
	PredicateExecuted
	PredicateObserved
)

type Event struct {
	Id          string         `json:"id"`
	Time        time.Time      `json:"time"`
	SubjectId   string         `json:"subject_id"`
	SubjectType ObservableType `json:"subject_type"`
	Predicate   Predicate      `json:"predicate"`
	ObjectId    string         `json:"object_id"`
	ObjectType  ObservableType `json:"object_type"`
}

func NewEvent(
	subjectId string,
	subjectType ObservableType,
	predicate Predicate,
	objectId string,
	objectType ObservableType) Event {

	return Event{
		Id:          uuid.NewString(),
		Time:        time.Now(),
		SubjectId:   subjectId,
		SubjectType: subjectType,
		Predicate:   predicate,
		ObjectId:    objectId,
		ObjectType:  objectType,
	}
}

type DataSourceEvent struct {
	Event
	DataSourceId string `json:"data_source_id"`
}

func NewDataSourceEvent(dataSourceId string, event Event) DataSourceEvent {
	return DataSourceEvent{
		Event:        event,
		DataSourceId: dataSourceId,
	}
}
