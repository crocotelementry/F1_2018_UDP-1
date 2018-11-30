package main

import (
	"flag"
	"log"
	"net/http"
  "net"
  "fmt"
  "encoding/json"
  "bytes"
  "encoding/binary"



	"github.com/gorilla/websocket"
  "github.com/gorilla/mux"
  "github.com/knilssen/f1_go/structs"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

func main() {
  var dir string;
  flag.StringVar(&dir, "dir", "./static", "the directory to serve files from. Defaults to the current dir")
  flag.Parse()

  router := mux.NewRouter()
  router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))

	router.HandleFunc("/", rootHandler).Methods("GET")
	router.HandleFunc("/ws", wsHandler)

	log.Fatal(http.ListenAndServe(":8080", router))


  // addr, _ := net.ResolveUDPAddr("udp", ":5003")
  // sock, _ := net.ListenUDP("udp", addr)


  // // Set up our structs including the header struct so we are able to determine which
  // // udp packet is incoming and we can deal with it accordingly
  // var header structs.PacketHeader
  // var motion_packet structs.PacketMotionData
  // var session_packet structs.PacketSessionData
  // var lap_packet structs.PacketLapData
  // var event_packet structs.PacketEventData
  // var participant_packet structs.PacketParticipantsData
  // var car_setup_packet structs.PacketCarSetupData
  // var telemetry_packet structs.PacketCarTelemetryData
  // var car_status_packet structs.PacketCarStatusData
}



func rootHandler(w http.ResponseWriter, r *http.Request) {
  // fmt.Println(r.Header)
  // if r.Header.Get("Origin") != "http://"+r.Host {
	// 	http.Error(w, "Origin not allowed", 403)
	// 	return
	// }
  http.ServeFile(w, r, "./static/telemetry_dashboard.html")
  // conn, err := upgrader.Upgrade(w, r, nil)
	// if err != nil {
	// 	http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	// }
}
//
// func rootHandler(w http.ResponseWriter, r *http.Request) {
// 	// fmt.Fprintf(w, "home")
//   r.ServeFile(w, r, "telemetry_dashboard.html")
// }






func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}
	go echo(conn)
}



func echo(conn *websocket.Conn) {
  // Loop, checking for udp packets and dealing with them accordingly
  addr, _ := net.ResolveUDPAddr("udp", ":5003")
  sock, _ := net.ListenUDP("udp", addr)

  // Set up our structs including the header struct so we are able to determine which
  // udp packet is incoming and we can deal with it accordingly
  var header structs.PacketHeader
  var motion_packet structs.PacketMotionData
  var session_packet structs.PacketSessionData
  var lap_packet structs.PacketLapData
  var event_packet structs.PacketEventData
  var participant_packet structs.PacketParticipantsData
  var car_setup_packet structs.PacketCarSetupData
  var telemetry_packet structs.PacketCarTelemetryData
  var car_status_packet structs.PacketCarStatusData

  for {
    // Create a buffer to read the incoming udp packets
    // Read the udp packets and if we get an error while reading, print out the error
    buf := make([]byte, 1341)
    _, _, err := sock.ReadFromUDP(buf)
    if err != nil {
      fmt.Println(err)
    }


    // Set two new readers in which to cast into our structs.
    // One is for the header, which we determine what packet we have and then use the other
    // for the whole packet.
    header_bytes_reader := bytes.NewReader(buf[0:21])
    packet_bytes_reader := bytes.NewReader(buf)

    // Read the binary of the udp packet header into our struct
    if err := binary.Read(header_bytes_reader, binary.LittleEndian, &header); err != nil {
		    fmt.Println("binary.Read header failed:", err)
  	}

    // Depending on which packet we have, which we find by looking at header.M_packetId
    // We use a switch statement to then read the whole binary udp packet into its associated struct
    switch header.M_packetId {
    case 0:
        // If the packet we received is a motion_packet, read its binary into our motion_packet struct
        if err := binary.Read(packet_bytes_reader, binary.LittleEndian, &motion_packet); err != nil {
    		    fmt.Println("binary.Read motion_packet failed:", err)
      	}
        // fmt.Println(motion_packet)
        json_motion_packet, err := json.Marshal(motion_packet)
        if err != nil {
          fmt.Println(err)
        }
        if err := conn.WriteMessage(websocket.TextMessage, json_motion_packet); err != nil {
    			log.Printf("Websocket error writing motion_packet: %s", err)
    		}
    case 1:
        // If the packet we received is the session_packet, read its binary into our session_packet struct
        if err := binary.Read(packet_bytes_reader, binary.LittleEndian, &session_packet); err != nil {
    		    fmt.Println("binary.Read session_packet failed:", err)
      	}
        // fmt.Println(session_packet)
        json_session_packet, err := json.Marshal(session_packet)
        if err != nil {
          fmt.Println(err)
        }
        if err := conn.WriteMessage(websocket.TextMessage, json_session_packet); err != nil {
    			log.Printf("Websocket error writing session_packet: %s", err)
    		}
    case 2:
        // If the packet we received is the lap_packet, read its binary into our lap_packet struct
        if err := binary.Read(packet_bytes_reader, binary.LittleEndian, &lap_packet); err != nil {
    		    fmt.Println("binary.Read lap_packet failed:", err)
      	}
        // fmt.Println(lap_packet)
        json_lap_packet, err := json.Marshal(lap_packet)
        if err != nil {
          fmt.Println(err)
        }
        if err := conn.WriteMessage(websocket.TextMessage, json_lap_packet); err != nil {
    			log.Printf("Websocket error writing lap_packet: %s", err)
    		}
    case 3:
        // If the packet we received is the event_packet, read its binary into our event_packet struct
        if err := binary.Read(packet_bytes_reader, binary.LittleEndian, &event_packet); err != nil {
    		    fmt.Println("binary.Read event_packet failed:", err)
      	}
        // fmt.Println(event_packet)
        json_event_packet, err := json.Marshal(event_packet)
        if err != nil {
          fmt.Println(err)
        }
        if err := conn.WriteMessage(websocket.TextMessage, json_event_packet); err != nil {
    			log.Printf("Websocket error writing event_packet: %s", err)
    		}
    case 4:
        // If the packet we received is the participant_packet, read its binary into our participant_packet struct
        if err := binary.Read(packet_bytes_reader, binary.LittleEndian, &participant_packet); err != nil {
    		    fmt.Println("binary.Read participant_packet failed:", err)
      	}
        json_participant_packet, err := json.Marshal(participant_packet)
        if err != nil {
          fmt.Println(err)
        }

        if err = conn.WriteMessage(websocket.TextMessage, json_participant_packet); err != nil {
    			log.Printf("Websocket error writing participant_packet: %s", err)
    		}
    case 5:
        // If the packet we received is the car_setup_packet, read its binary into our car_setup_packet struct
        if err := binary.Read(packet_bytes_reader, binary.LittleEndian, &car_setup_packet); err != nil {
    		    fmt.Println("binary.Read car_setup_packet failed:", err)
      	}
        // fmt.Println(car_setup_packet)
        json_car_setup_packet, err := json.Marshal(car_setup_packet)
        if err != nil {
          fmt.Println(err)
        }
        if err := conn.WriteMessage(websocket.TextMessage, json_car_setup_packet); err != nil {
    			log.Printf("Websocket error writing car_setup_packet: %s", err)
    		}
    case 6:
        // If the packet we received is the telemetry_packet, read its binary into our telemetry_packet struct
        if err := binary.Read(packet_bytes_reader, binary.LittleEndian, &telemetry_packet); err != nil {
    		    fmt.Println("binary.Read telemetry_packet failed:", err)
      	}
        // fmt.Println(telemetry_packet)
        json_telemetry_packet, err := json.Marshal(telemetry_packet)
        if err != nil {
          fmt.Println(err)
        }
        if err := conn.WriteMessage(websocket.TextMessage, json_telemetry_packet); err != nil {
    			log.Printf("Websocket error writing telemetry_packet: %s", err)
    		}
    case 7:
        // If the packet we received is the car_status_packet, read its binary into our car_status_packet struct
        if err := binary.Read(packet_bytes_reader, binary.LittleEndian, &car_status_packet); err != nil {
    		    fmt.Println("binary.Read car_status_packet failed:", err)
      	}
        // fmt.Println(car_status_packet)
        json_car_status_packet, err := json.Marshal(car_status_packet)
        if err != nil {
          fmt.Println(err)
        }
        if err := conn.WriteMessage(websocket.TextMessage, json_car_status_packet); err != nil {
    			log.Printf("Websocket error writing car_status_packet: %s", err)
    		}
    }

  }
}
