// Copyright 2020 The Cloud Native Events Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hwevent

import (
	"bytes"
	"fmt"
	"io"

	jsoniter "github.com/json-iterator/go"
)

// WriteJSON writes the in event in the provided writer.
// Note: this function assumes the input event is valid.
func WriteJSON(in *Event, writer io.Writer) error {
	stream := jsoniter.ConfigFastest.BorrowStream(writer)
	defer jsoniter.ConfigFastest.ReturnStream(stream)
	stream.WriteObjectStart()
	if in.DataContentType != nil {
		switch in.GetDataContentType() {
		case ApplicationJSON:
			stream.WriteObjectField("id")
			stream.WriteString(in.ID)
			stream.WriteMore()

			stream.WriteObjectField("type")
			stream.WriteString(in.GetType())

			if in.GetDataContentType() != "" {
				stream.WriteMore()
				stream.WriteObjectField("dataContentType")
				stream.WriteString(in.GetDataContentType())
			}

			if in.Time != nil {
				stream.WriteMore()
				stream.WriteObjectField("time")
				stream.WriteString(in.Time.String())
			}

			if in.GetDataSchema() != "" {
				stream.WriteMore()
				stream.WriteObjectField("dataSchema")
				stream.WriteString(in.GetDataSchema())
			}
		default:
			return fmt.Errorf("missing event content type")
		}
	}

	// Let's do a check on the error
	if stream.Error != nil {
		return fmt.Errorf("error while writing the event attributes: %w", stream.Error)
	}

	// Let's write the body
	data := in.GetData()
	if data != nil {
		stream.WriteMore()
		stream.WriteObjectField("data")
		if err := writeJSONData(data, writer, stream); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("data is not set")
	}
	stream.WriteObjectEnd()
	// Let's do a check on the error
	if stream.Error != nil {
		return fmt.Errorf("error while writing the event Data: %w", stream.Error)
	}

	// Let's do a check on the error
	if stream.Error != nil {
		return fmt.Errorf("error while writing the event extensions: %w", stream.Error)
	}
	return stream.Flush()
}

// WriteDataJSON writes the in data in the provided writer.
// Note: this function assumes the input event is valid.
func WriteDataJSON(in *Data, writer io.Writer) error {
	stream := jsoniter.ConfigFastest.BorrowStream(writer)
	defer jsoniter.ConfigFastest.ReturnStream(stream)
	if err := writeJSONData(in, writer, stream); err != nil {
		return err
	}
	return stream.Flush()
}
func writeJSONData(in *Data, writer io.Writer, stream *jsoniter.Stream) error {
	stream.WriteObjectStart()

	// Let's write the body
	if in != nil {
		data := in
		stream.WriteObjectField("version")
		stream.WriteString(data.GetVersion())
		stream.WriteMore()
		stream.WriteObjectField("data")
		if err := writeJSONRedfishEvent(data.Data, writer, stream); err != nil {
			return fmt.Errorf("error writing data: %w", err)
		}
		stream.WriteObjectEnd()
	} else {
		return fmt.Errorf("data version is not set")
	}

	// Let's do a check on the error
	if stream.Error != nil {
		return fmt.Errorf("error while writing the event Data: %w", stream.Error)
	}

	// Let's do a check on the error
	if stream.Error != nil {
		return fmt.Errorf("error while writing the event extensions: %w", stream.Error)
	}
	return nil
}

func writeJSONRedfishEvent(in *RedfishEvent, writer io.Writer, stream *jsoniter.Stream) error {
	stream.WriteObjectStart()

	// Let's write the body
	if in != nil {
		var err error
		data := in
		if data.OdataContext != "" {
			stream.WriteObjectField("@odata.context")
			stream.WriteString(data.OdataContext)
			stream.WriteMore()
		}
		if data.Actions != nil {
			stream.WriteObjectField("Actions")
			_, err = stream.Write(data.Actions)
			if err != nil {
				return fmt.Errorf("error writing Actions: %w", err)
			}
			stream.WriteMore()
		}
		if data.Context != "" {
			stream.WriteObjectField("Context")
			stream.WriteString(data.Context)
			stream.WriteMore()
		}
		if data.Description != "" {
			stream.WriteObjectField("Description")
			stream.WriteString(data.Description)
			stream.WriteMore()
		}
		if data.Oem != nil {
			stream.WriteObjectField("Oem")
			_, err = stream.Write(data.Oem)
			if err != nil {
				return fmt.Errorf("error writing Oem: %w", err)
			}
		}
		if data.OdataType != "" {
			stream.WriteObjectField("@odata.type")
			stream.WriteString(data.OdataType)
			stream.WriteMore()
		} else {
			return fmt.Errorf("@odata.type is not set")
		}
		if data.Events != nil {
			stream.WriteObjectField("Events")
			stream.WriteArrayStart()
			for i, v := range data.Events {
				if err := writeJSONEventRecord(&v, writer, stream); err != nil {
					return fmt.Errorf("error writing Event[%d]: %w", i, err)
				}
			}
			stream.WriteArrayEnd()
			stream.WriteMore()
		} else {
			return fmt.Errorf("field Events is not set")
		}
		if data.ID != "" {
			stream.WriteObjectField("Id")
			stream.WriteString(data.ID)
			stream.WriteMore()
		} else {
			return fmt.Errorf("field Id is not set")
		}
		if data.Name != "" {
			stream.WriteObjectField("Name")
			stream.WriteString(data.Name)
		} else {
			return fmt.Errorf("field Name is not set")
		}

		stream.WriteObjectEnd()
	} else {
		return fmt.Errorf("@odata.type is not set")
	}

	// Let's do a check on the error
	if stream.Error != nil {
		return fmt.Errorf("error while writing the event Data: %w", stream.Error)
	}

	// Let's do a check on the error
	if stream.Error != nil {
		return fmt.Errorf("error while writing the event extensions: %w", stream.Error)
	}
	return nil
}

func writeJSONEventRecord(in *EventRecord, writer io.Writer, stream *jsoniter.Stream) error {
	stream.WriteObjectStart()

	// Let's write the body
	if in != nil {
		var err error
		data := in
		if data.Actions != nil {
			stream.WriteObjectField("Actions")
			_, err = stream.Write(data.Actions)
			if err != nil {
				return fmt.Errorf("error writing Oem: %w", err)
			}
			stream.WriteMore()
		}
		if data.Context != "" {
			stream.WriteObjectField("Context")
			stream.WriteString(data.Context)
			stream.WriteMore()
		}
		stream.WriteObjectField("EventGroupId")
		stream.WriteInt(data.EventGroupID)
		stream.WriteMore()
		if data.EventID != "" {
			stream.WriteObjectField("EventId")
			stream.WriteString(data.EventID)
			stream.WriteMore()
		}
		if data.EventTimestamp != "" {
			stream.WriteObjectField("EventTimestamp")
			stream.WriteString(data.EventTimestamp)
			stream.WriteMore()
		}
		if data.Message != "" {
			stream.WriteObjectField("Message")
			stream.WriteString(data.Message)
			stream.WriteMore()
		}
		if data.MessageArgs != nil {
			stream.WriteObjectField("MessageArgs")
			stream.WriteArrayStart()
			count := 0
			for _, v := range data.MessageArgs {
				if count > 0 {
					stream.WriteMore()
				}
				count++
				stream.WriteString(v)
			}
			stream.WriteArrayEnd()
			stream.WriteMore()
		}
		if data.Oem != nil {
			stream.WriteObjectField("Oem")
			_, err = stream.Write(data.Oem)
			if err != nil {
				return fmt.Errorf("error writing Oem: %w", err)
			}
			stream.WriteMore()
		}
		if data.OriginOfCondition != "" {
			stream.WriteObjectField("OriginOfCondition")
			stream.WriteString(data.OriginOfCondition)
			stream.WriteMore()
		}
		if data.Severity != "" {
			stream.WriteObjectField("Severity")
			stream.WriteString(data.Severity)
			stream.WriteMore()
		}
		if data.Resolution != "" {
			stream.WriteObjectField("Resolution")
			stream.WriteString(data.Resolution)
			stream.WriteMore()
		}
		if data.EventType != "" {
			stream.WriteObjectField("EventType")
			stream.WriteString(data.EventType)
			stream.WriteMore()
		} else {
			return fmt.Errorf("field EventType is not set")
		}
		if data.MessageID != "" {
			stream.WriteObjectField("MessageId")
			stream.WriteString(data.MessageID)
			stream.WriteMore()
		} else {
			return fmt.Errorf("field MessageId is not set")
		}
		if data.MemberID != "" {
			stream.WriteObjectField("MemberId")
			stream.WriteString(data.MemberID)
		} else {
			return fmt.Errorf("field MemberId is not set")
		}

		stream.WriteObjectEnd()
	} else {
		return fmt.Errorf("field EventType is not set")
	}

	// Let's do a check on the error
	if stream.Error != nil {
		return fmt.Errorf("error while writing the event Data: %w", stream.Error)
	}

	// Let's do a check on the error
	if stream.Error != nil {
		return fmt.Errorf("error while writing the event extensions: %w", stream.Error)
	}
	return nil
}

// MarshalJSON implements a custom json marshal method used when this type is
// marshaled using json.Marshal.
func (e Event) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	err := WriteJSON(&e, &buf)
	return buf.Bytes(), err
}

// MarshalJSON implements a custom json marshal method used when this type is
// marshaled using json.Marshal.
func (d Data) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	err := WriteDataJSON(&d, &buf)
	return buf.Bytes(), err
}
